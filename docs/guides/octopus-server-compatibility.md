---
page_title: "Octopus Server compatibility"
subcategory: "Upgrades & Migrations"
---

# Version Matrix
Provider supports major versions of Octopus Server.
Provider may work with older versions of Octopus Server, see [Breaking Changes](./breaking-changes-list.md) and [Breaking Changes Policy](./breaking-changes-policy.md) for possible incompatibility 

## Supported versions

| Server Version | Comment     |
| -------------- | ----------- |
| 2025.2         | Recommended |
| 2025.1         |             |
| 2024.4         |             |

## Partially supported resources

The table below shows provider resources and attributes which have limited compatibility with Octopus Server

| Resource                                                                                                                                                    | Server Version  | Comment                                                                                            |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------- | -------------------------------------------------------------------------------------------------- |
| [octopusdeploy_machine_policy.machine_package_cache_retention_policy](./../resources/machine_policy.md#nestedblock--machine_package_cache_retention_policy) | 2025.3 - latest | _Machine package cache rention settings is only available from 2025.3_                             |
| [octopusdeploy_deployment_freeze](./../resources/deployment_freeze.md)                                                                                      | 2025.1 - latest | _Resource were available in earlier versions, but provider is compatible only from version 2025.1_ |
| [octopusdeploy_deployment_freeze_project](./../resources/deployment_freeze_project.md)                                                                      | 2025.1 - latest |                                                                                                    |
| [octopusdeploy_deployment_freeze_tenant](./../resources/deployment_freeze_tenant.md)                                                                        | 2025.1 - latest |                                                                                                    |

