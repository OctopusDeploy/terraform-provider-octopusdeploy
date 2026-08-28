#!/usr/bin/env python3
"""
Extract API endpoints and contract properties from Octopus Server C# source.
API contracts: classes ending in Resource, Command, Request (and Details subtypes for accounts).
Run: python extract_server_source.py <path-to-Octopus-Server-source>
e.g. python extract_server_source.py C:/Source/Octopus/Server/source
"""
import json
import re
import sys
from pathlib import Path
from collections import defaultdict

def main():
    if len(sys.argv) < 2:
        print("Usage: extract_server_source.py <path-to-Octopus-Server-source>", file=sys.stderr)
        sys.exit(1)
    server_root = Path(sys.argv[1])
    if not server_root.is_dir():
        print(f"Not a directory: {server_root}", file=sys.stderr)
        sys.exit(1)

    # --- 1. Extract endpoints from controllers ---
    # Base route prefixes by controller base class (from SpaceScopedWithDefaultSpaceApiController.cs etc.)
    route_prefixes = {
        "SpaceScopedWithDefaultSpaceApiController": ["api/", "api/{spaceId}/", "api/spaces/{spaceIdentifier}/"],
        "SpaceScopedApiController": ["api/{spaceId}/", "api/spaces/{spaceIdentifier}/"],
        "SystemScopedApiController": ["api/"],
        "MixedScopedApiController": ["api/", "api/{spaceId}/", "api/spaces/{spaceId}/"],
        "SpaceScopedBffController": ["bff/spaces/{spaceId}/"],
        "SystemScopedBffController": ["bff/"],
        "MixedScopedBffController": ["bff/", "bff/spaces/{spaceId}/"],
    }
    controllers_dir = server_root / "Octopus.Server" / "Web" / "Controllers"
    endpoints = []
    if controllers_dir.exists():
        for cs_file in controllers_dir.rglob("*.cs"):
            text = cs_file.read_text(encoding="utf-8", errors="replace")
            # Base class: class FooController : SpaceScopedWithDefaultSpaceApiController
            base_match = re.search(r"class\s+\w+\s*\([^)]*\)\s*:\s*(\w+)|class\s+\w+\s*:\s*(\w+)", text)
            base_class = None
            if base_match:
                base_class = (base_match.group(1) or base_match.group(2) or "").strip()
            prefixes = route_prefixes.get(base_class, ["api/"])
            # HttpGet("path"), HttpPost("path"), etc. - path can be @"path" or "path"
            for method in ["HttpGet", "HttpPost", "HttpPut", "HttpDelete", "HttpPatch"]:
                for m in re.finditer(rf'\[{method}\s*\(\s*@?"([^"]+)"\s*\)\]', text):
                    path = m.group(1).replace("\\", "")
                    for prefix in prefixes:
                        full = (prefix.rstrip("/") + "/" + path.lstrip("/")).replace("//", "/")
                        endpoints.append({"method": method.replace("Http", "").upper(), "route": full})
                # Multiple routes on one attribute: [HttpGet("a"), HttpGet("b")]
                for m in re.finditer(rf'\[{method}\s*\(\s*@?"([^"]+)"\s*\)\s*\]', text):
                    path = m.group(1).replace("\\", "")
                    for prefix in prefixes:
                        full = (prefix.rstrip("/") + "/" + path.lstrip("/")).replace("//", "/")
                        endpoints.append({"method": method.replace("Http", "").upper(), "route": full})
            # Also match [HttpGet("path")] without @
            for method in ["HttpGet", "HttpPost", "HttpPut", "HttpDelete", "HttpPatch"]:
                for m in re.finditer(rf'\[{method}\s*\(\s*"([^"]+)"\s*\)\]', text):
                    path = m.group(1)
                    for prefix in prefixes:
                        full = (prefix.rstrip("/") + "/" + path.lstrip("/")).replace("//", "/")
                        endpoints.append({"method": method.replace("Http", "").upper(), "route": full})

    # Dedupe endpoints
    seen = set()
    unique_endpoints = []
    for e in endpoints:
        key = (e["method"], e["route"])
        if key not in seen:
            seen.add(key)
            unique_endpoints.append(e)

    # --- 2. Extract contract properties from *Resource.cs, *Command.cs, *Request.cs ---
    # Search in Octopus.Core (Resources, Features, Accounts), MessageContracts, Sashimi.*
    contract_dirs = [
        server_root / "Octopus.Core" / "Resources",
        server_root / "Octopus.Core" / "Features",
        server_root / "Octopus.Core" / "Accounts",
        server_root / "Octopus.Server.MessageContracts",
        server_root / "Octopus.Server.Extensibility.Sashimi.Azure.Accounts",
        server_root / "Octopus.Server.Extensibility.Sashimi.Aws.Accounts",
        server_root / "Octopus.Server.Extensibility.Sashimi.Server.Contracts",
    ]
    # class name -> { "properties": [...], "base": base_class_name or None }
    contracts = {}
    # Match property name: allow { get or { set (so auto-props and get/set with backing field both match); allow nullable types (?)
    property_re = re.compile(r"public\s+(?:readonly\s+)?[\w<>,\s\[\]?]+\s+(\w+)\s*\{\s*(?:get|set)", re.MULTILINE)
    for d in contract_dirs:
        if not d.exists():
            continue
        for cs_file in d.rglob("*.cs"):
            if "Test" in cs_file.parts or ".Tests" in cs_file.parts:
                continue
            text = cs_file.read_text(encoding="utf-8", errors="replace")
            for m in re.finditer(r"class\s+(\w+)\s*:\s*([^,{]+)\s*[,{]|class\s+(\w+)\s*\{", text):
                class_name = m.group(1) or m.group(3)
                base_part = m.group(2).strip() if m.group(2) else None
                if not (class_name.endswith("Resource") or class_name.endswith("Command") or
                        class_name.endswith("Request") or class_name.endswith("Details")):
                    continue
                # First identifier of base (e.g. "ResourceBase<AccountResource>" -> "ResourceBase", "AccountDetailsResource" -> "AccountDetailsResource")
                base_class = None
                if base_part:
                    first = base_part.split("<")[0].strip().split()[-1]  # last token before < or space
                    if first and first != "class" and first[0].isupper():
                        base_class = first
                props = list(dict.fromkeys(property_re.findall(text)))
                if class_name not in contracts:
                    contracts[class_name] = {"properties": [], "base": base_class}
                for p in props:
                    if p not in contracts[class_name]["properties"]:
                        contracts[class_name]["properties"].append(p)
                if base_class and contracts[class_name]["base"] is None:
                    contracts[class_name]["base"] = base_class

    # Merge inherited properties: each contract gets own props + all ancestor props
    def get_all_properties(name, seen=None):
        seen = seen or set()
        if name not in contracts or name in seen:
            return []
        seen.add(name)
        entry = contracts[name]
        base = entry.get("base")
        base_props = get_all_properties(base, seen) if base and base in contracts else []
        return list(dict.fromkeys(base_props + entry["properties"]))

    merged = {}
    for name in contracts:
        if name.endswith("Resource") or name.endswith("Command") or name.endswith("Request") or name.endswith("Details"):
            merged[name] = sorted(set(get_all_properties(name)))

    # AccountDetailsResource and its subclasses are always returned inside AccountResource; include AccountResource props
    account_resource_props = merged.get("AccountResource", [])
    if account_resource_props:
        def has_account_details_ancestor(n, seen=None):
            seen = seen or set()
            if n not in contracts or n in seen:
                return False
            seen.add(n)
            if n == "AccountDetailsResource":
                return True
            b = contracts[n].get("base")
            return b and has_account_details_ancestor(b, seen)
        for name, props in list(merged.items()):
            if has_account_details_ancestor(name) or name == "AccountDetailsResource":
                merged[name] = sorted(set(props) | set(account_resource_props))

    contracts = merged

    # --- 3. Output ---
    out_dir = Path(__file__).resolve().parent.parent.parent / "api-coverage-out"
    out_dir.mkdir(exist_ok=True)
    (out_dir / "server_endpoints.json").write_text(json.dumps(sorted(unique_endpoints, key=lambda x: (x["route"], x["method"])), indent=2), encoding="utf-8")
    contracts_list = [{"name": k, "properties": v} for k, v in sorted(contracts.items())]
    (out_dir / "server_contracts.json").write_text(json.dumps(contracts_list, indent=2), encoding="utf-8")
    print("Endpoints:", len(unique_endpoints))
    print("Contracts:", len(contracts))
    print("Output:", out_dir)


if __name__ == "__main__":
    main()
