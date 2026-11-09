#!/usr/bin/env python3
"""
Full API vs Terraform provider coverage comparison.
Extracts TF schemas, maps API paths/definitions to provider resources, diffs properties.
Run from repo root: python tools/api-coverage/compare_coverage.py
"""
import json
import re
import os
from pathlib import Path
from collections import defaultdict

REPO_ROOT = Path(__file__).resolve().parent.parent.parent
SCHEMAS_DIR = REPO_ROOT / "octopusdeploy_framework" / "schemas"
OUT_DIR = REPO_ROOT / "api-coverage-out"
FRAMEWORK_PROVIDER = REPO_ROOT / "octopusdeploy_framework" / "framework_provider.go"
SDK_PROVIDER = REPO_ROOT / "octopusdeploy" / "provider.go"


def pascal_to_snake(s):
    """Convert PascalCase to snake_case."""
    return re.sub(r"(?<!^)(?=[A-Z])", "_", s).lower()


def extract_tfsdk_attributes(filepath):
    """Extract all tfsdk attribute names from a Go file."""
    text = filepath.read_text(encoding="utf-8")
    return list(dict.fromkeys(re.findall(r'tfsdk:"([^"]+)"', text)))


def extract_framework_resources():
    """Parse framework_provider.go for NewXxxResource -> octopusdeploy_xxx."""
    text = FRAMEWORK_PROVIDER.read_text(encoding="utf-8")
    # NewProjectResource -> project, NewSpaceDefaultLifecycleReleaseRetentionPolicyResource -> space_default_lifecycle_release_retention_policy
    pattern = r"New(\w+Resource)"
    matches = re.findall(pattern, text)
    result = []
    for m in matches:
        # NewProjectResource -> project, NewSpaceDefaultLifecycleReleaseRetentionPolicyResource -> space_default_lifecycle_release_retention_policy
        name = m
        if name.endswith("Resource"):
            base = name[: -len("Resource")]
        else:
            base = name
        # CamelCase to snake_case for resource type
        tf_name = pascal_to_snake(base)
        result.append(("framework", f"octopusdeploy_{tf_name}", name))
    return result


def extract_sdk_resources():
    """Parse provider.go ResourcesMap for octopusdeploy_xxx."""
    text = SDK_PROVIDER.read_text(encoding="utf-8")
    pattern = r'"octopusdeploy_([^"]+)"\s*:\s*resource\w+\('
    matches = re.findall(pattern, text)
    return [("sdk", f"octopusdeploy_{m}", None) for m in matches]


def schema_file_to_resource_name(schema_file_stem):
    """Map schema filename (without .go) to TF resource name. project -> octopusdeploy_project."""
    return "octopusdeploy_" + schema_file_stem.replace("_", "_")


def build_schema_attributes_by_file():
    """Build map: schema_file_stem -> set(attribute names)."""
    out = {}
    for f in SCHEMAS_DIR.glob("*.go"):
        if f.name.startswith("schema_") and "test" not in f.name:
            stem = f.stem.replace("schema_", "")
        else:
            stem = f.stem
        attrs = extract_tfsdk_attributes(f)
        if stem not in out:
            out[stem] = set()
        out[stem].update(attrs)
    return out


# API Resource name (without 'Resource') to primary schema file stem(s) for Framework.
# Partial list; we extend by convention: XxxResource -> xxx.go
API_RESOURCE_TO_SCHEMA = {
    "Account": ["username_password_account", "certificate", "amazon_web_services_account",
                "azure_subscription_account", "generic_oidc_account"],  # multiple
    "Project": ["project"],
    "Environment": ["environment"],
    "Lifecycle": ["lifecycle"],
    "Space": ["space"],
    "Channel": ["channel"],
    "Team": ["team"],
    "Tenant": ["tenant"],
    "User": ["user"],
    "Variable": ["variable"],
    "ActionTemplate": ["step_template"],
    "Feed": ["feed", "docker_container_registry_feed", "nuget_feed", "npm_feed", "maven_feed",
             "helm_feed", "github_repository_feed", "s3_feed", "oci_registry_feed",
             "google_container_registry_feed", "azure_container_registry_feed"],
    "Certificate": ["certificate"],
    "LibraryVariableSet": ["library_variable_set"],
    "MachineProxy": ["machine_proxy"],
    "Runbook": ["runbook"],
    "DeploymentFreeze": ["deployment_freeze"],
    "Process": ["process"],
    "Worker": ["listening_tentacle_worker", "ssh_connection_worker"],
    "ScopedUserRole": ["scoped_user_role"],
    "TagSet": ["tag_set"],
    "Tag": ["tag"],
    "ProjectGroup": ["project_group"],
    "ParentEnvironment": ["parent_environment"],
    "ScriptModule": ["script_modules"],
    "GitCredential": ["git_credential"],
    "TentacleCertificate": ["schema_tentacle_certificate"],
    "KubernetesMonitor": ["kubernetes_monitor"],
    "ServiceAccountOidcIdentity": ["service_account_oidc_identity"],
    "SpaceDefaultLifecycleReleaseRetentionPolicy": ["space_default_lifecycle_release_retention_policy"],
    "SpaceDefaultLifecycleTentacleRetentionPolicy": ["space_default_lifecycle_tentacle_retention_policy"],
    "SpaceDefaultRunbookRetentionPolicy": ["space_default_runbooks_retention_policy"],
    "TenantProject": ["tenant_projects"],
    "DeploymentFreezeProject": ["project_deployment_freeze"],
    "ProjectAutoCreateRelease": ["project_auto_create_release"],
    "ProjectVersioningStrategy": ["project_versioning_strategy"],
    "ProcessStep": ["process_step"],
    "ProcessChildStep": ["process_child_step"],
    "ProcessTemplatedStep": ["process_templated_step"],
    "ProcessTemplatedChildStep": ["process_templated_child_step"],
    "DeploymentFreezeTenant": ["deployment_freeze_tenant"],
}


def main():
    # Load API extraction output
    paths = json.loads((OUT_DIR / "paths.json").read_text(encoding="utf-8"))
    path_resources = json.loads((OUT_DIR / "path_resources.json").read_text(encoding="utf-8"))
    defs_resources = json.loads((OUT_DIR / "definitions_resources.json").read_text(encoding="utf-8"))
    tags = json.loads((OUT_DIR / "tags.json").read_text(encoding="utf-8"))

    # Build TF schema attributes by file
    schema_by_file = build_schema_attributes_by_file()

    # Framework + SDK resource list
    framework_rs = extract_framework_resources()
    sdk_rs = extract_sdk_resources()
    all_tf_resources = [(r[1]) for r in framework_rs + sdk_rs]

    # API Resource name -> TF attributes (union of mapped schema files)
    def get_tf_attrs_for_api_resource(api_name):
        if not api_name.endswith("Resource"):
            return set()
        base = api_name[: -len("Resource")]
        stems = API_RESOURCE_TO_SCHEMA.get(base)
        if not stems:
            # Convention: ProjectResource -> project
            stems = [pascal_to_snake(base)]
        attrs = set()
        for stem in stems:
            if stem in schema_by_file:
                attrs.update(schema_by_file[stem])
            # try with schema_ prefix
            if f"schema_{stem}" in schema_by_file:
                attrs.update(schema_by_file[f"schema_{stem}"])
        return attrs

    # Path resource -> TF coverage. Normalize key to lowercase for lookup.
    path_resource_to_tf = {}
    def norm(s):
        return s.lower().replace("-", "")

    # Manual mapping: normalized path resource -> list of TF resource names
    mapping = {
        "accounts": [r[1] for r in framework_rs + sdk_rs if "account" in r[1]],
        "actiontemplates": ["octopusdeploy_step_template", "octopusdeploy_community_step_template"],
        "artifacts": [],
        "buildinformation": [],
        "certificates": ["octopusdeploy_certificate", "octopusdeploy_tentacle_certificate"],
        "channels": ["octopusdeploy_channel"],
        "communityactiontemplates": ["octopusdeploy_community_step_template"],
        "configuration": [],
        "dashboard": [],
        "dashboardconfiguration": [],
        "deploymentfreezes": ["octopusdeploy_deployment_freeze", "octopusdeploy_deployment_freeze_project", "octopusdeploy_deployment_freeze_tenant", "octopusdeploy_project_deployment_freeze"],
        "deploymentprocesses": ["octopusdeploy_process", "octopusdeploy_process_step", "octopusdeploy_process_child_step"] + [r[1] for r in sdk_rs if "process" in r[1]],
        "deployments": [],
        "deploymentsettings": [],
        "deploymenttargettags": [],
        "environments": ["octopusdeploy_environment", "octopusdeploy_parent_environment"],
        "events": [],
        "feeds": [r[1] for r in framework_rs if "feed" in r[1]],
        "gitcredentials": ["octopusdeploy_git_credential", "octopusdeploy_platform_hub_git_credential"],
        "libraryvariablesets": ["octopusdeploy_library_variable_set"],
        "lifecycles": ["octopusdeploy_lifecycle"],
        "machines": [r[1] for r in sdk_rs if "deployment_target" in r[1] or "worker" in r[1]] + ["octopusdeploy_listening_tentacle_worker", "octopusdeploy_ssh_connection_worker", "octopusdeploy_kubernetes_monitor"],
        "machinepolicies": [r[1] for r in sdk_rs if "machine_policy" in r[1]],
        "parentenvironments": ["octopusdeploy_parent_environment"],
        "projectgroups": ["octopusdeploy_project_group"],
        "projects": ["octopusdeploy_project", "octopusdeploy_project_versioning_strategy", "octopusdeploy_project_auto_create_release"],
        "projecttriggers": ["octopusdeploy_git_trigger", "octopusdeploy_built_in_trigger"] + [r[1] for r in sdk_rs if "trigger" in r[1]],
        "proxies": ["octopusdeploy_machine_proxy"],
        "releases": [],
        "retentionpolicies": ["octopusdeploy_space_default_lifecycle_release_retention_policy", "octopusdeploy_space_default_lifecycle_tentacle_retention_policy", "octopusdeploy_space_default_runbook_retention_policy"],
        "runbooks": ["octopusdeploy_runbook"],
        "runbookprocesses": ["octopusdeploy_runbook"] + [r[1] for r in sdk_rs if "runbook" in r[1]],
        "scopeduserroles": ["octopusdeploy_scoped_user_role"],
        "spaces": ["octopusdeploy_space"],
        "tagsets": ["octopusdeploy_tag_set", "octopusdeploy_tag"],
        "teams": ["octopusdeploy_team"],
        "tenants": ["octopusdeploy_tenant"],
        "tenantvariables": ["octopusdeploy_tenant_project_variable", "octopusdeploy_tenant_common_variable"],
        "users": ["octopusdeploy_user"],
        "userroles": [r[1] for r in sdk_rs if "user_role" in r[1]] + ["octopusdeploy_scoped_user_role"],
        "variables": ["octopusdeploy_variable", "octopusdeploy_tenant_project_variable", "octopusdeploy_tenant_common_variable"],
        "workerpools": [r[1] for r in sdk_rs if "worker_pool" in r[1]],
        "workers": ["octopusdeploy_listening_tentacle_worker", "octopusdeploy_ssh_connection_worker", "octopusdeploy_kubernetes_monitor"] + [r[1] for r in sdk_rs if "worker" in r[1]],
    }
    for pr in path_resources:
        n = norm(pr)
        if n in mapping:
            path_resource_to_tf[pr] = mapping[n]
        else:
            # Try substring match
            found = []
            for _, tf_name, _ in framework_rs + sdk_rs:
                r_short = tf_name.replace("octopusdeploy_", "").replace("_", "")
                if n in r_short or r_short in n:
                    found.append(tf_name)
            path_resource_to_tf[pr] = list(dict.fromkeys(found))

    def get_coverage_for_path_segment(segment):
        """Resolve path resource segment to TF resources (normalize to match any key)."""
        n = norm(segment)
        for key, tf_list in path_resource_to_tf.items():
            if norm(key) == n:
                return tf_list
        return []

    # Build full report
    report = []
    report.append("# Full API vs Terraform Provider Coverage")
    report.append("")
    report.append("Generated from swagger and provider codebase.")
    report.append("")
    report.append("---")
    report.append("")
    report.append("## 1. API paths (all endpoints)")
    report.append("")
    report.append("| Path | Methods | Covered by TF |")
    report.append("|------|---------|----------------|")
    for p in paths:
        path = p["path"]
        methods = ", ".join(sorted(set(p["methods"])))
        parts = path.strip("/").split("/")
        res = ""
        for i, part in enumerate(parts):
            if part == "spaces" and i + 2 < len(parts):
                res = parts[i + 2]
                break
            if part == "{spaceId}" and i + 1 < len(parts):
                res = parts[i + 1]
                break
            if part and not part.startswith("{"):
                res = part
                break
        covered = get_coverage_for_path_segment(res) if res else []
        if covered:
            covered_str = ", ".join(sorted(set(covered))[:3])
            if len(covered) > 3:
                covered_str += f" (+{len(covered)-3} more)"
        else:
            covered_str = "—"
        report.append(f"| `{path}` | {methods} | {covered_str} |")

    report.append("")
    report.append("## 2. Path resources (unique API areas) and TF coverage")
    report.append("")
    report.append("| API path resource | TF resources | Status |")
    report.append("|-------------------|--------------|--------|")
    seen_norm = set()
    for pr in sorted(path_resources):
        if norm(pr) in seen_norm:
            continue
        seen_norm.add(norm(pr))
        tf_list = get_coverage_for_path_segment(pr)
        if tf_list:
            report.append(f"| {pr} | {', '.join(sorted(set(tf_list))[:5])}{' ...' if len(tf_list) > 5 else ''} | ✅ |")
        else:
            report.append(f"| {pr} | — | ❌ Missing |")

    report.append("")
    report.append("## 3. API tags (operation groups)")
    report.append("")
    report.append("| Tag |")
    report.append("|-----|")
    for t in sorted(tags):
        report.append(f"| {t} |")

    report.append("")
    report.append("## 4. API *Resource definitions – property coverage")
    report.append("")
    report.append("For each API resource type, we list properties that appear in the API but not in the mapped Terraform schema (missing or different name).")
    report.append("")

    for dr in defs_resources:
        api_name = dr["name"]
        api_props = set(dr["properties"])
        read_only = set(dr.get("read_only", []))
        tf_attrs = get_tf_attrs_for_api_resource(api_name)
        api_snake = {pascal_to_snake(p) for p in api_props}
        # Standard exclusions: Id->id, Links, LastModifiedBy, LastModifiedOn, SpaceId
        api_snake.discard("links")
        api_snake.discard("id")  # id is in TF
        missing_in_tf = api_snake - tf_attrs
        # Exclude read-only that are often not in TF
        missing_in_tf -= {pascal_to_snake(p) for p in read_only}
        if missing_in_tf and api_name not in ("AccountResourceCollection", "AccountUsageResource", "ActionTemplateSearchResource", "ActionTemplateUsageResource"):
            report.append(f"### {api_name}")
            report.append("")
            report.append("| API property (snake_case) | In TF? |")
            report.append("|--------------------------|--------|")
            for ap in sorted(api_props):
                sn = pascal_to_snake(ap)
                in_tf = "✅" if sn in tf_attrs else ("—" if ap in read_only else "❌")
                report.append(f"| {ap} (`{sn}`) | {in_tf} |")
            report.append("")

    report.append("## 5. Terraform resources (Framework + SDK)")
    report.append("")
    report.append("| Source | Resource name |")
    report.append("|--------|---------------|")
    for src, tf_name, _ in framework_rs:
        report.append(f"| Framework | `{tf_name}` |")
    for src, tf_name, _ in sdk_rs:
        report.append(f"| SDK | `{tf_name}` |")

    out_md = REPO_ROOT / "docs" / "API_COVERAGE_FULL.md"
    out_md.write_text("\n".join(report), encoding="utf-8")
    print("Report written to", out_md)
    print("Path resources:", len(path_resources), "with coverage:", sum(1 for pr in path_resources if path_resource_to_tf.get(pr)))
    print("API *Resource definitions:", len(defs_resources))
    print("Schema files:", len(schema_by_file))


if __name__ == "__main__":
    main()
