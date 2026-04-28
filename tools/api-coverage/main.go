// Package main extracts paths and definitions from Octopus Deploy swagger.json
// for comparison with the Terraform provider. Run: go run . from this directory.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Swagger struct {
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
}

type PathInfo struct {
	Path   string   `json:"path"`
	Methods []string `json:"methods"`
	Tags   []string `json:"tags"`
}

type DefinitionInfo struct {
	Name       string   `json:"name"`
	Properties []string `json:"properties"`
	ReadOnly   []string `json:"read_only,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-swagger.json>\n", os.Args[0])
		os.Exit(1)
	}
	swaggerPath := os.Args[1]
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read swagger: %v\n", err)
		os.Exit(1)
	}
	var swagger Swagger
	if err := json.Unmarshal(data, &swagger); err != nil {
		fmt.Fprintf(os.Stderr, "Parse swagger: %v\n", err)
		os.Exit(1)
	}

	outDir := filepath.Join(filepath.Dir(swaggerPath), "api-coverage-out")
	os.MkdirAll(outDir, 0755)

	// --- Extract all paths ---
	var paths []PathInfo
	for path, v := range swagger.Paths {
		vm, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		var methods []string
		var tags []string
		for method, op := range vm {
			if method[0] >= 'a' && method[0] <= 'z' {
				methods = append(methods, strings.ToUpper(method))
				if opm, ok := op.(map[string]interface{}); ok {
					if t, ok := opm["tags"].([]interface{}); ok {
						for _, tag := range t {
							if s, ok := tag.(string); ok {
								tags = append(tags, s)
							}
						}
					}
				}
			}
		}
		sort.Strings(methods)
		paths = append(paths, PathInfo{Path: path, Methods: methods, Tags: tags})
	}
	sort.Slice(paths, func(i, j int) bool { return paths[i].Path < paths[j].Path })
	pathBytes, _ := json.MarshalIndent(paths, "", "  ")
	os.WriteFile(filepath.Join(outDir, "paths.json"), pathBytes, 0644)

	// --- Unique path "resources" (first segment after /spaces/ or top-level) ---
	resourceSet := make(map[string]struct{})
	for _, p := range paths {
		path := p.Path
		path = strings.TrimPrefix(path, "/")
		parts := strings.Split(path, "/")
		var resource string
		for i, part := range parts {
			if part == "spaces" && i+1 < len(parts) {
				// /spaces/{spaceIdentifier}/accounts -> accounts
				resource = parts[i+2]
				break
			}
			if part == "{spaceId}" && i+1 < len(parts) {
				resource = parts[i+1]
				break
			}
			if !strings.HasPrefix(part, "{") && part != "" {
				resource = part
				break
			}
		}
		if resource != "" {
			resourceSet[resource] = struct{}{}
		}
	}
	var resources []string
	for r := range resourceSet {
		resources = append(resources, r)
	}
	sort.Strings(resources)
	resBytes, _ := json.MarshalIndent(resources, "", "  ")
	os.WriteFile(filepath.Join(outDir, "path_resources.json"), resBytes, 0644)

	// --- Extract all definitions (focus on *Resource and *Command) ---
	var defs []DefinitionInfo
	for name, v := range swagger.Definitions {
		vm, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		props, ok := vm["properties"].(map[string]interface{})
		if !ok {
			continue
		}
		var propNames []string
		var readOnly []string
		for propName, propVal := range props {
			propNames = append(propNames, propName)
			if pv, ok := propVal.(map[string]interface{}); ok {
				if ro, _ := pv["readOnly"].(bool); ro {
					readOnly = append(readOnly, propName)
				}
			}
		}
		sort.Strings(propNames)
		sort.Strings(readOnly)
		defs = append(defs, DefinitionInfo{Name: name, Properties: propNames, ReadOnly: readOnly})
	}
	sort.Slice(defs, func(i, j int) bool { return defs[i].Name < defs[j].Name })
	defBytes, _ := json.MarshalIndent(defs, "", "  ")
	os.WriteFile(filepath.Join(outDir, "definitions.json"), defBytes, 0644)

	// --- Only *Resource definitions (for property comparison) ---
	var resourceDefs []DefinitionInfo
	for _, d := range defs {
		if strings.HasSuffix(d.Name, "Resource") && d.Name != "Resource" && d.Name != "AccountDetailsResource" {
			resourceDefs = append(resourceDefs, d)
		}
	}
	resDefBytes, _ := json.MarshalIndent(resourceDefs, "", "  ")
	os.WriteFile(filepath.Join(outDir, "definitions_resources.json"), resDefBytes, 0644)

	// --- All tags from paths ---
	tagSet := make(map[string]struct{})
	for _, p := range paths {
		for _, t := range p.Tags {
			tagSet[t] = struct{}{}
		}
	}
	var tags []string
	for t := range tagSet {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	tagBytes, _ := json.MarshalIndent(tags, "", "  ")
	os.WriteFile(filepath.Join(outDir, "tags.json"), tagBytes, 0644)

	fmt.Printf("Paths: %d\n", len(paths))
	fmt.Printf("Path resources: %d\n", len(resources))
	fmt.Printf("Definitions: %d\n", len(defs))
	fmt.Printf("*Resource definitions: %d\n", len(resourceDefs))
	fmt.Printf("Tags: %d\n", len(tags))
	fmt.Printf("Output: %s\n", outDir)
}
