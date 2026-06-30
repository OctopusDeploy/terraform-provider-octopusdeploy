#!/usr/bin/env python3
"""
Compare Terraform provider against Octopus Server API source (controllers + Resource/Command/Request).
Uses api-coverage-out/server_endpoints.json and server_contracts.json (from extract_server_source.py).
Run from repo root: python tools/api-coverage/compare_server_source.py
"""
import json
import re
import sys
from pathlib import Path
from collections import defaultdict

REPO_ROOT = Path(__file__).resolve().parent.parent.parent
SCHEMAS_DIR = REPO_ROOT / "octopusdeploy_framework" / "schemas"
OUT_DIR = REPO_ROOT / "api-coverage-out"


def pascal_to_snake(s):
    return re.sub(r"(?<!^)(?=[A-Z])", "_", s).lower()


def extract_tfsdk_attributes(filepath):
    text = filepath.read_text(encoding="utf-8")
    return list(dict.fromkeys(re.findall(r'tfsdk:"([^"]+)"', text)))


def extract_schema_map_attribute_keys(filepath):
    """Extract attribute names from schema Attributes map keys (e.g. \"description\": Get...). Used when schema uses programmatic API without tfsdk tags."""
    text = filepath.read_text(encoding="utf-8")
    # Match "snake_case_key": at start of line (with indent) followed by Get, util., or resourceSchema. / datasourceSchema.
    return list(dict.fromkeys(re.findall(r'^\s*"([a-z][a-z0-9_]*)"\s*:\s*(?:Get|util\.|resourceSchema\.|datasourceSchema\.)', text, re.MULTILINE)))


def build_schema_attributes_by_file():
    out = {}
    for f in SCHEMAS_DIR.glob("*.go"):
        if "test" in f.name:
            continue
        if f.name.startswith("schema_"):
            stem = f.stem.replace("schema_", "")
        else:
            stem = f.stem
        attrs = set(extract_tfsdk_attributes(f)) | set(extract_schema_map_attribute_keys(f))
        if stem not in out:
            out[stem] = set()
        out[stem].update(attrs)
    return out


# Server contract name -> TF schema file stem(s). Details subtypes map to specific TF resources.
# AccountResource is not mapped: compare concrete account types (e.g. UsernamePasswordAccountResource) to their TF resources instead.
SERVER_CONTRACT_TO_TF_SCHEMA = {
    "GenericOidcAccountResource": ["generic_oidc_account"],
    "AzureOidcAccountResource": ["azure_subscription_account"],  # Azure OIDC in Sashimi
    "AmazonWebServicesOidcAccountResource": ["amazon_web_services_account"],
    "ProjectResource": ["project"],
    "EnvironmentResource": ["environment"],
    "LifecycleResource": ["lifecycle"],
    "SpaceResource": ["space"],
    "ChannelResource": ["channel"],
    "TeamResource": ["team"],
    "TenantResource": ["tenant"],
    "UserResource": ["user"],
    "VariableResource": ["variable"],
    "ActionTemplateResource": ["step_template"],
    "CertificateResource": ["certificate"],  # octopusdeploy_certificate
    "LibraryVariableSetResource": ["library_variable_set"],
    "MachineProxyResource": ["machine_proxy"],
    "RunbookResource": ["runbook"],
    "DeploymentFreezeResource": ["deployment_freeze"],
    "ProcessResource": ["process"],
    "WorkerResource": ["listening_tentacle_worker", "ssh_connection_worker"],
    "ScopedUserRoleResource": ["scoped_user_role"],
    "TagSetResource": ["tag_set"],
    "TagResource": ["tag"],
    "ProjectGroupResource": ["project_group"],
    "ParentEnvironmentResource": ["parent_environment"],
    "ScriptModuleResource": ["script_modules"],
    "GitCredentialResource": ["git_credential"],
    "TentacleCertificateResource": ["schema_tentacle_certificate"],
    "KubernetesMonitorResource": ["kubernetes_monitor"],
    "ServiceAccountOidcIdentityResource": ["service_account_oidc_identity"],
    "SpaceDefaultLifecycleReleaseRetentionPolicyResource": ["space_default_lifecycle_release_retention_policy"],
    "SpaceDefaultLifecycleTentacleRetentionPolicyResource": ["space_default_lifecycle_tentacle_retention_policy"],
    "SpaceDefaultRunbookRetentionPolicyResource": ["space_default_runbooks_retention_policy"],
    "TenantProjectResource": ["tenant_projects"],
    "DeploymentFreezeProjectResource": ["project_deployment_freeze"],
    "ProjectAutoCreateReleaseResource": ["project_auto_create_release"],
    "ProjectVersioningStrategyResource": ["project_versioning_strategy"],
    "ProcessStepResource": ["process_step"],
    "ProcessChildStepResource": ["process_child_step"],
    "ProcessTemplatedStepResource": ["process_templated_step"],
    "ProcessTemplatedChildStepResource": ["process_templated_child_step"],
    "DeploymentFreezeTenantResource": ["deployment_freeze_tenant"],
    "UsernamePasswordAccountResource": ["username_password_account"],
    "AzureServicePrincipalAccountResource": ["azure_subscription_account"],
    "AmazonWebServicesAccountResource": ["amazon_web_services_account"],
}
# Property names that represent the same concept (e.g. API uses environment_ids, TF uses environments).
PROPERTY_EQUIVALENTS = [
    {"environment_ids", "environments"},
    {"tenant_ids", "tenants"},
]


def _equiv_set(prop):
    """Return the set of names equivalent to prop (including prop itself), or {prop} if none."""
    for group in PROPERTY_EQUIVALENTS:
        if prop in group:
            return group
    return {prop}


def _in_other(prop, other_set):
    """True if prop or any equivalent is in other_set."""
    return bool(_equiv_set(prop) & other_set)


# Convention: XxxResource -> xxx
def get_tf_schema_stems_for_contract(name):
    if name in SERVER_CONTRACT_TO_TF_SCHEMA:
        return SERVER_CONTRACT_TO_TF_SCHEMA[name]
    if name.endswith("Resource"):
        base = name[:-len("Resource")]
        return [pascal_to_snake(base)]
    return []


# Schema stem -> Terraform resource type name (octopusdeploy_*). Some stems differ from TF type.
STEM_TO_TF_TYPE_OVERRIDES = {
    "amazon_web_services_account": "octopusdeploy_aws_account",
    "tenant_projects": "octopusdeploy_tenant_project",
    "project_deployment_freeze": "octopusdeploy_deployment_freeze_project",
    "space_default_runbooks_retention_policy": "octopusdeploy_space_default_runbook_retention_policy",
    "script_modules": "octopusdeploy_script_module",
}
def stem_to_tf_type(stem):
    if stem in STEM_TO_TF_TYPE_OVERRIDES:
        return STEM_TO_TF_TYPE_OVERRIDES[stem]
    # e.g. schema_tentacle_certificate -> octopusdeploy_tentacle_certificate
    s = stem.replace("schema_", "") if stem.startswith("schema_") else stem
    return f"octopusdeploy_{s}"


def schema_stem_for_lookup(stem):
    """Return schema_by_file key(s) for this stem."""
    if stem.startswith("schema_"):
        return [stem, stem.replace("schema_", "")]
    return [stem]


def get_all_tf_resource_types():
    """All Terraform resource type names (Framework + SDK). Resources only (no data sources)."""
    out = set()
    # SDK: from provider.go ResourcesMap only (exclude DataSourcesMap)
    sdk_provider = REPO_ROOT / "octopusdeploy" / "provider.go"
    if sdk_provider.exists():
        text = sdk_provider.read_text(encoding="utf-8")
        start = text.find("ResourcesMap: map[string]*schema.Resource{")
        end = text.find("Schema: map[string]*schema.Schema", start) if start >= 0 else -1
        if start >= 0 and end >= 0:
            block = text[start:end]
            for m in re.findall(r'"octopusdeploy_([a-z0-9_]+)"', block):
                out.add(f"octopusdeploy_{m}")
    # Framework: resource_*.go GetTypeName("...") or GetTypeName(schemas.X)
    framework_dir = REPO_ROOT / "octopusdeploy_framework"
    schema_const_to_stem = _build_schema_const_to_stem()
    for f in (framework_dir / "resource_*.go").parent.glob("resource_*.go"):
        text = f.read_text(encoding="utf-8")
        for m in re.findall(r'GetTypeName\("([^"]+)"\)', text):
            out.add(f"octopusdeploy_{m}")
        for m in re.findall(r"GetTypeName\(schemas\.(\w+)\)", text):
            if m in schema_const_to_stem:
                out.add(f"octopusdeploy_{schema_const_to_stem[m]}")
    return out


def _build_schema_const_to_stem():
    """Map schema constant name to value (stem) by parsing schemas/*.go."""
    out = {}
    for f in SCHEMAS_DIR.glob("*.go"):
        if "test" in f.name:
            continue
        text = f.read_text(encoding="utf-8")
        for m in re.findall(r'const\s+(\w+)\s*=\s*"([^"]+)"', text):
            out[m[0]] = m[1]
    return out


def main():
    server_contracts_path = OUT_DIR / "server_contracts.json"
    server_endpoints_path = OUT_DIR / "server_endpoints.json"
    if not server_contracts_path.exists() or not server_endpoints_path.exists():
        print("Run extract_server_source.py <Server-source-path> first.", file=sys.stderr)
        sys.exit(1)
    contracts_list = json.loads(server_contracts_path.read_text(encoding="utf-8"))
    endpoints = json.loads(server_endpoints_path.read_text(encoding="utf-8"))
    schema_by_file = build_schema_attributes_by_file()
    all_tf_resources = get_all_tf_resource_types()

    def has_schema(stem):
        for key in schema_stem_for_lookup(stem):
            if key in schema_by_file:
                return True
        return False

    def contract_in_tf(contract_name):
        stems = get_tf_schema_stems_for_contract(contract_name)
        if not stems:
            return False
        for s in stems:
            if has_schema(s) and stem_to_tf_type(s) in all_tf_resources:
                return True
        return False

    tf_resources_from_api = set()
    for c in contracts_list:
        if not c["name"].endswith("Resource"):
            continue
        for stem in get_tf_schema_stems_for_contract(c["name"]):
            if has_schema(stem):
                tf_resources_from_api.add(stem_to_tf_type(stem))

    server_resource_contracts = [c for c in contracts_list if c["name"].endswith("Resource")]
    server_not_in_tf = sorted([c["name"] for c in server_resource_contracts if not contract_in_tf(c["name"])])
    tf_not_in_api = sorted(all_tf_resources - tf_resources_from_api)

    # Build contract name -> properties
    contracts = {c["name"]: set(c["properties"]) for c in contracts_list}

    def get_writable_props(resource_name):
        """Properties on the Modify command for this resource (writable); empty set if no Modify command."""
        if not resource_name.endswith("Resource"):
            return set()
        modify_name = "Modify" + resource_name.replace("Resource", "") + "Command"
        return contracts.get(modify_name, set())

    # Path resource -> TF coverage (reuse same mapping as compare_coverage for endpoint coverage)
    path_resource_to_tf = {
        "accounts": ["octopusdeploy_username_password_account", "octopusdeploy_certificate", "octopusdeploy_generic_oidc_account"],
        "actiontemplates": ["octopusdeploy_step_template", "octopusdeploy_community_step_template"],
        "certificates": ["octopusdeploy_certificate", "octopusdeploy_tentacle_certificate"],
        "channels": ["octopusdeploy_channel"],
        "environments": ["octopusdeploy_environment", "octopusdeploy_parent_environment"],
        "feeds": [],  # populated below
        "libraryvariablesets": ["octopusdeploy_library_variable_set"],
        "lifecycles": ["octopusdeploy_lifecycle"],
        "projectgroups": ["octopusdeploy_project_group"],
        "projects": ["octopusdeploy_project"],
        "proxies": ["octopusdeploy_machine_proxy"],
        "retentionpolicies": ["octopusdeploy_space_default_lifecycle_release_retention_policy", "octopusdeploy_space_default_lifecycle_tentacle_retention_policy", "octopusdeploy_space_default_runbook_retention_policy"],
        "runbooks": ["octopusdeploy_runbook"],
        "spaces": ["octopusdeploy_space"],
        "tagsets": ["octopusdeploy_tag_set", "octopusdeploy_tag"],
        "teams": ["octopusdeploy_team"],
        "tenants": ["octopusdeploy_tenant"],
        "users": ["octopusdeploy_user"],
        "variables": ["octopusdeploy_variable", "octopusdeploy_tenant_project_variable", "octopusdeploy_tenant_common_variable"],
        "workers": ["octopusdeploy_listening_tentacle_worker", "octopusdeploy_ssh_connection_worker", "octopusdeploy_kubernetes_monitor"],
        "workerpools": [],
        "deploymentfreezes": ["octopusdeploy_deployment_freeze", "octopusdeploy_deployment_freeze_project", "octopusdeploy_deployment_freeze_tenant"],
        "parentenvironments": ["octopusdeploy_parent_environment"],
        "projecttriggers": ["octopusdeploy_git_trigger", "octopusdeploy_built_in_trigger"],
        "scopeduserroles": ["octopusdeploy_scoped_user_role"],
        "deploymentprocesses": ["octopusdeploy_process"],
        "runbookprocesses": ["octopusdeploy_runbook"],
    }
    def norm(s):
        return s.lower().replace("-", "").replace("_", "")
    def get_coverage_for_route(route):
        # route like api/spaces/{id}/accounts -> accounts
        parts = route.strip("/").split("/")
        for i, p in enumerate(parts):
            if p == "spaces" and i + 2 < len(parts):
                return path_resource_to_tf.get(norm(parts[i + 2]), [])
            if p == "api" and i + 1 < len(parts):
                return path_resource_to_tf.get(norm(parts[i + 1]), [])
        return []

    skip_read_only = {"Id", "Links", "LastModifiedBy", "LastModifiedOn", "SpaceId", "Slug", "VariableSetId", "DeploymentProcessId"}
    skip_diff_snake = {pascal_to_snake(p) for p in skip_read_only}  # also exclude these from "missing from API" so they are not listed as differing
    # Type discriminators (e.g. account_type, feed_type) determine the concrete resource in Terraform; do not list as missing
    def is_type_discriminator(snake_name):
        return snake_name.endswith("_type")

    def get_tf_attrs_for_stem(stem):
        attrs = set()
        for key in schema_stem_for_lookup(stem):
            if key in schema_by_file:
                attrs.update(schema_by_file[key])
        return attrs

    report = []
    report.append("# Terraform Provider vs Octopus Server API (Source Code)")
    report.append("")
    report.append("Comparison against the **Octopus Server C# source** (controllers + Resource/Command/Request contracts).")
    report.append("")
    report.append("---")
    report.append("")
    report.append("## 1. Resources in the Server/API not in Terraform")
    report.append("")
    report.append("Server `*Resource` contracts that have **no corresponding Terraform resource** (no mapping or no schema/provider registration).")
    report.append("")
    if server_not_in_tf:
        for name in server_not_in_tf:
            report.append(f"- **{name}**")
        report.append("")
    else:
        report.append("_(None)_")
        report.append("")
    report.append("---")
    report.append("")
    report.append("## 2. Resources in Terraform not in the Server/API")
    report.append("")
    report.append("Terraform resources that are **not the target of any Server contract** in the current mapping (Framework + SDK).")
    report.append("")
    if tf_not_in_api:
        for name in tf_not_in_api:
            report.append(f"- **{name}**")
        report.append("")
    else:
        report.append("_(None)_")
        report.append("")
    report.append("---")
    report.append("")
    report.append("## 3. Property differences (Server/API vs Terraform)")
    report.append("")
    report.append("Resources that exist in both: for each, properties that **differ** are listed with whether the property is **missing from Terraform** (API has it, TF does not) or **missing from Server/API** (TF has it, API does not).")
    report.append("")
    report.append("**Note:** `*_ids` and the corresponding resource name are treated as the same (e.g. `environment_ids` ↔ `environments`, `tenant_ids` ↔ `tenants`) and are not listed as differing. Properties ending in `_type` (e.g. `account_type`, `feed_type`) are type discriminators that determine the concrete resource in Terraform and are not listed as missing. **Read-only properties** (those that exist on the Resource but not on the corresponding `ModifyXxxCommand`) are excluded from \"Missing from Terraform\" so only writable API properties are compared.")
    report.append("")

    # Per (contract, stem) pair where both exist, list differing properties
    for c in contracts_list:
        name = c["name"]
        if not name.endswith("Resource"):
            continue
        stems = get_tf_schema_stems_for_contract(name)
        if not stems:
            continue
        api_props = set(c["properties"])
        api_snake = {pascal_to_snake(p) for p in api_props} - {pascal_to_snake(p) for p in skip_read_only}
        writable_props = get_writable_props(name)
        api_snake_writable = {pascal_to_snake(p) for p in api_props if p in writable_props} - {pascal_to_snake(p) for p in skip_read_only} if writable_props else api_snake
        for stem in stems:
            if not has_schema(stem) or stem_to_tf_type(stem) not in all_tf_resources:
                continue
            tf_attrs = get_tf_attrs_for_stem(stem)
            # Only consider writable props (on Modify command) for "missing from TF"; equivalences, skip list, type discriminators
            missing_from_tf = {p for p in api_snake_writable if not _in_other(p, tf_attrs) and not is_type_discriminator(p)}
            missing_from_api = {a for a in tf_attrs if not _in_other(a, api_snake) and a not in skip_diff_snake and not is_type_discriminator(a)}
            if not missing_from_tf and not missing_from_api:
                continue
            tf_type = stem_to_tf_type(stem)
            report.append(f"### {name} ↔ **{tf_type}**")
            report.append("")
            report.append("| Property | Missing from Terraform | Missing from Server/API |")
            report.append("|----------|------------------------|-------------------------|")
            all_props = sorted(missing_from_tf | missing_from_api)
            for sn in all_props:
                from_tf = "✓" if sn in missing_from_tf else ""
                from_api = "✓" if sn in missing_from_api else ""
                report.append(f"| `{sn}` | {from_tf} | {from_api} |")
            report.append("")

    report.append("---")
    report.append("")
    report.append("## 4. Server endpoints and TF coverage (reference)")
    report.append("")
    report.append("| Method | Route | TF coverage |")
    report.append("|--------|-------|-------------|")
    for e in sorted(endpoints, key=lambda x: (x["route"], x["method"]))[:500]:
        cov = get_coverage_for_route(e["route"])
        cov_str = ", ".join(cov[:2]) + (" ..." if len(cov) > 2 else "") if cov else "—"
        report.append(f"| {e['method']} | `{e['route']}` | {cov_str} |")
    if len(endpoints) > 500:
        report.append(f"| ... | _(and {len(endpoints) - 500} more)_ | |")
    report.append("")

    out_md = REPO_ROOT / "docs" / "API_COVERAGE_SERVER_SOURCE.md"
    out_md.write_text("\n".join(report), encoding="utf-8")
    print("Report written to", out_md)
    print("Section 1 - Server/API not in Terraform:", len(server_not_in_tf))
    print("Section 2 - Terraform not in Server/API:", len(tf_not_in_api))
    print("Endpoints:", len(endpoints), "| Contracts:", len(contracts_list))


if __name__ == "__main__":
    main()
