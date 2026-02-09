# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build and install to local Terraform plugins folder
make

# Build only (without install)
make build

# Generate documentation
make docs

# Run unit tests
go test -i $(go list ./... | grep -v 'vendor') || exit 1

# Run acceptance tests (requires running Octopus Deploy instance)
TF_ACC=1 go test ./... -v -timeout 120m

# Run specific test(s) with BYO Octopus instance
TF_ACC_LOCAL=1 OCTOPUS_URL='http://localhost:<port>' OCTOPUS_APIKEY='<key>' \
  go test -run "^(?:TestAccOctopusDeployNuGetFeedBasic|TestAccResourceBuiltInTrigger)$" -timeout 0 ./...

# Run test with auto-created container
LICENSE=<base64-license> OCTOTESTIMAGEURL='docker.packages.octopushq.com/octopusdeploy/octopusdeploy' \
  go test -run "^(?:TestName)$" -timeout 0 ./... -createSharedContainer=true
```

## Architecture

This is a Terraform provider for Octopus Deploy that uses a **muxed provider architecture** combining:
- **Terraform Plugin Framework** (`octopusdeploy_framework/`) - New resources go here
- **Terraform Plugin SDK v2** (`octopusdeploy/`) - Legacy resources, no new additions allowed

The `main.go` combines both providers using `tf6muxserver`, allowing gradual migration from SDK to Framework.

### Key Directories

- `octopusdeploy_framework/` - Framework-based resources and data sources (preferred for new code)
  - `schemas/` - Reusable schema definitions
  - `util/` - Helper utilities
- `octopusdeploy/` - SDK-based resources (legacy, do not add new resources here)
- `internal/` - Shared utilities (errors, test helpers, deprecation)
- `docs/` - Auto-generated documentation (run `make docs`)
- `examples/` - Terraform example configurations
- `templates/` - Documentation templates for tfplugindocs

### Creating New Resources

New resources MUST use Terraform Plugin Framework in `octopusdeploy_framework/`. A GitHub action blocks SDK additions.

Resource pattern:
1. Create resource file: `resource_<name>.go`
2. Create schema in `schemas/<name>.go`
3. Add to `framework_provider.go` Resources() list
4. Create acceptance test: `resource_<name>_test.go`

Use nested attributes instead of blocks (blocks are only for SDK migration compatibility).

### Server Compatibility

Enforce version/feature requirements in resource's Configure method:
```go
f.Config.EnsureResourceCompatibilityByVersion(resourceName, "2025.1")
f.Config.EnsureResourceCompatibilityByFeature(resourceName, "FeatureToggleName")
```

## Conventions

- Use [conventional commits](https://www.conventionalcommits.org/) for commit messages
- Squash and merge PRs for clean git history
- Provider configured via `OCTOPUS_URL`, `OCTOPUS_APIKEY`/`OCTOPUS_API_KEY`, or `OCTOPUS_ACCESS_TOKEN` env vars
- Local debug provider source: `octopus.com/com/octopusdeploy`
- Registry provider source: `OctopusDeploy/octopusdeploy`

## Debugging

1. Use GoLand run configuration "Run provider" or `dlv debug . -- --debug`
2. Export the `TF_REATTACH_PROVIDERS` env var from provider output
3. Run terraform commands in the same terminal session
