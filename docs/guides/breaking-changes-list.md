---
page_title: "Breaking Changes List"
subcategory: "Upgrades & Migrations"
---

# Breaking Changes
This page details the breaking changes and deprecations we're managing according to our [Breaking Changes policy](./breaking-changes-policy.md).

## Announced

| Version | Deprecated                                              | Replacement | Migration Guide | Enactment | Completion |
|---------|---------------------------------------------------------|-------------|-----------------|-----------|------------|
| [v0.37.1](https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/releases/tag/v0.37.1) | `octopusdeploy_project.versioning_strategy` (attribute) | `octopusdeploy_project_versioning_strategy` (new resource) | [Guide](./migration-guide-v0.37.1.md) | 2025-06-04 | 2025-12-04 |
| [v1.0.0](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/releases/tag/v1.0.0) | `octopusdeploy_deployment_process` | `octopusdeploy_process` | [Guide](./migration-guide-v1.0.0.md) | 2025-11-05 | 2026-05-05 |
| [v1.0.0](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/releases/tag/v1.0.0) | `octopusdeploy_runbook_process` | `octopusdeploy_process` | [Guide](./migration-guide-v1.0.0.md) | 2025-11-05 | 2026-05-05 |
| [v1.3.0](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/releases/tag/v1.3.0) | `octopusdeploy_auto_create_release` | `octopusdeploy_built_in_trigger` | [Guide](./migration-guide-v1.3.0.md) | 2025-02-08 | 2026-08-08 |

## Enacted

We aren't currently tracking any enacted Breaking Changes.


## Completed

| Version | Deprecated | Replacement | Migration Guide | Enactment | Completion |
|---------|------------|-------------|-----------------|-----------|------------|
| [v0.8.0](https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/releases/tag/v0.8.0) | Many changes* | Many changes* | [Guide](./migration-guide-v0.8.0.md) | N/A | N/A |


-> *The breaking changes in v0.8.0 were all done years before our Breaking Changes policy was officially established, so didn't follow those procedures, but the Migration Guide is noted here just in case. 
