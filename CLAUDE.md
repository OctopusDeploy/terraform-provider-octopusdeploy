# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the official Terraform Provider for Octopus Deploy that enables Infrastructure as Code management of Octopus Deploy resources. The project uses Go 1.22+ and implements a dual-provider architecture migrating from HashiCorp's Plugin SDK to the modern Plugin Framework.

## Common Commands

```bash
# Build and install locally
make install            # Default target - builds and installs to local Terraform plugins directory
make build              # Build binary only
make test               # Run unit tests
make testacc            # Run acceptance tests (requires TF_ACC=1)
make docs               # Generate documentation using tfplugindocs
make release            # Build cross-platform binaries

# Integration testing (choose one approach)
TF_ACC_LOCAL=1 OCTOPUS_URL='http://localhost:port' OCTOPUS_APIKEY='key' go test ./...
go test -createSharedContainer=true ./...

# Documentation generation
go generate main.go
```

## Architecture Overview

### Dual Provider Structure

The codebase implements **provider multiplexing** combining two providers:

1. **Legacy SDK Provider** (`/octopusdeploy/`)
   - Uses HashiCorp Terraform Plugin SDK v2
   - Status: Legacy, no new resources allowed
   - Contains existing resources with `resource_*.go` and `schema_*.go` files

2. **Modern Framework Provider** (`/octopusdeploy_framework/`)
   - Uses HashiCorp Terraform Plugin Framework
   - Status: Active development, all new resources go here
   - Organized with `schemas/` directory and `util/` helpers

3. **Provider Multiplexing** (`main.go`)
   - Uses `tf6muxserver` to combine both providers
   - Upgrades SDK provider to Protocol 6 for compatibility

### Key Directories

- `internal/` - Shared utilities, test helpers, error handling
- `examples/` - Usage examples and guides  
- `ci-scripts/` - Build automation and CI helpers
- `docs/` - Auto-generated documentation (do not edit manually)
- `templates/` - Documentation templates for generation

## Development Rules

### Resource Development
- **New resources MUST go in `octopusdeploy_framework/`** using the Plugin Framework
- A GitHub Action prevents new SDK additions automatically
- Framework resources should implement server version compatibility checks:
  ```go
  func (r *resourceName) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
      r.Config = ResourceConfiguration(req, resp)
      if r.Config != nil {
          diags := r.Config.EnsureResourceCompatibilityByVersion("resourceName", "2025.1")
          resp.Diagnostics.Append(diags...)
      }
  }
  ```

### Testing Requirements
- All resources need acceptance tests covering CRUD lifecycle
- Use containerized Octopus testing when possible with `testcontainers-go`
- Integration tests support both BYO Octopus instances and containerized testing

### Schema Organization
- **Framework resources**: Use structured schemas in `schemas/` directory
- **SDK resources**: Inline schema definitions (legacy pattern)
- Prefer nested attributes over blocks in Framework resources

## Environment Variables

### Required for Testing
```bash
OCTOPUS_URL=<octopus-server-url>
OCTOPUS_APIKEY=<api-key>           # or OCTOPUS_API_KEY
TF_ACC=1                           # Required for acceptance tests
```

### Optional for Advanced Testing
```bash
LICENSE=<base64-octopus-license>   # For containerized testing
OCTOTESTIMAGEURL=<docker-image>    # Octopus container image
OCTOTESTVERSION=<version>          # Container version
```

## Local Development Setup

1. Run `make install` to build and install locally
2. Version must match `VERSION` in Makefile (currently 0.7.102)
3. Use local provider in Terraform configs:
   ```hcl
   terraform {
     required_providers {
       octopusdeploy = {
         source  = "octopus.com/com/octopusdeploy"
         version = "0.7.102"
       }
     }
   }
   ```

## Documentation

- Documentation is auto-generated from templates in `/templates/`
- Never edit files in `/docs/` directly
- Run `make docs` after schema changes
- Templates use Go template syntax with provider metadata