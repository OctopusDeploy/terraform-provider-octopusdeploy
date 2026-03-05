#!/usr/bin/env python3
"""Extract paths and definitions from Octopus Deploy swagger.json."""
import json
import os
import sys
from pathlib import Path

def main():
    if len(sys.argv) < 2:
        print("Usage: extract_swagger.py <path-to-swagger.json>", file=sys.stderr)
        sys.exit(1)
    swagger_path = Path(sys.argv[1])
    data = swagger_path.read_text(encoding="utf-8")
    swagger = json.loads(data)
    out_dir = swagger_path.parent / "api-coverage-out"
    out_dir.mkdir(exist_ok=True)

    paths_obj = swagger.get("paths", {})
    paths_list = []
    for path, v in paths_obj.items():
        if not isinstance(v, dict):
            continue
        methods = []
        tags = []
        for method, op in v.items():
            if method and method[0].islower():
                methods.append(method.upper())
                if isinstance(op, dict) and "tags" in op:
                    tags.extend(op["tags"])
        paths_list.append({"path": path, "methods": sorted(methods), "tags": list(dict.fromkeys(tags))})
    paths_list.sort(key=lambda x: x["path"])
    (out_dir / "paths.json").write_text(json.dumps(paths_list, indent=2), encoding="utf-8")

    resource_set = set()
    for p in paths_list:
        path = p["path"].strip("/")
        parts = path.split("/")
        resource = ""
        for i, part in enumerate(parts):
            if part == "spaces" and i + 2 < len(parts):
                resource = parts[i + 2]
                break
            if part == "{spaceId}" and i + 1 < len(parts):
                resource = parts[i + 1]
                break
            if part and not part.startswith("{"):
                resource = part
                break
        if resource:
            resource_set.add(resource)
    (out_dir / "path_resources.json").write_text(
        json.dumps(sorted(resource_set), indent=2), encoding="utf-8"
    )

    defs_obj = swagger.get("definitions", {})
    defs_list = []
    resource_defs = []
    for name, v in defs_obj.items():
        if not isinstance(v, dict) or "properties" not in v:
            continue
        props = v["properties"]
        prop_names = sorted(props.keys())
        read_only = sorted(k for k, p in props.items() if isinstance(p, dict) and p.get("readOnly"))
        defs_list.append({"name": name, "properties": prop_names, "read_only": read_only})
        if name.endswith("Resource") and name not in ("Resource", "AccountDetailsResource"):
            resource_defs.append({"name": name, "properties": prop_names, "read_only": read_only})
    defs_list.sort(key=lambda x: x["name"])
    resource_defs.sort(key=lambda x: x["name"])
    (out_dir / "definitions.json").write_text(json.dumps(defs_list, indent=2), encoding="utf-8")
    (out_dir / "definitions_resources.json").write_text(
        json.dumps(resource_defs, indent=2), encoding="utf-8"
    )

    tag_set = set()
    for p in paths_list:
        tag_set.update(p["tags"])
    (out_dir / "tags.json").write_text(json.dumps(sorted(tag_set), indent=2), encoding="utf-8")

    print("Paths:", len(paths_list))
    print("Path resources:", len(resource_set))
    print("Definitions:", len(defs_list))
    print("*Resource definitions:", len(resource_defs))
    print("Tags:", len(tag_set))
    print("Output:", out_dir)

if __name__ == "__main__":
    main()
