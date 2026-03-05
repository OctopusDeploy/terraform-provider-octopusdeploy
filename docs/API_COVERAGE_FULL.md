# Full API vs Terraform Provider Coverage

Generated from swagger (2026.1) and provider codebase. This document lists **all** paths, path resources, API tags, and *Resource definition properties with TF coverage and property-level diff.

## Executive summary

| Metric | Count |
|--------|-------|
| **API paths (endpoints)** | 1,314 |
| **Path resources (unique API areas)** | 102 (normalized, deduped) |
| **Path resources with TF coverage** | 45 |
| **Path resources with no TF resource** | 57 |
| **API tags (operation groups)** | 105 |
| **API *Resource definitions** | 233 |
| **TF Framework resources** | 76 |
| **TF SDK resources** | 24 |

**Missing path resources (no TF resource):** artifacts, build-information, configuration, dashboard, dashboardconfiguration, deploymenttargettags, deployments, deploymentsettings, deprecations, dynamic-extensions, events, externalgroups, externalsecuritygroupproviders, featuresconfiguration, githubissuetracker, icons, insights, interruptions, jiraintegration, jiraservicemanagement-integration, letsencryptconfiguration, licenses, machineroles, maintenanceconfiguration, migrations, observability, octopusservernodes, packages, permissions, progression, reporting, scheduler, search, serverconfiguration, serverstatus, serviceaccounts, servicenow-integration, signingkeyconfiguration, smtpconfiguration, subscriptions, tasks, telemetry, telemetryconfiguration, upgradeconfiguration, v1, workertaskleases, and a few others (see Section 2).

---

## 1. API paths (all endpoints)

| Path | Methods | Covered by TF |
|------|---------|----------------|
| `/` | GET | — |
| `/.well-known/jwks` | GET | — |
| `/.well-known/openid-configuration` | GET | — |
| `/accounts` | GET, POST | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/all` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/azureenvironments` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{accountId}` | PUT | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}` | DELETE, GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/pk` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/resourceGroups` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/storageAccounts` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/usages` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/v1` | DELETE | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/websites` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/accounts/{id}/{resourceGroupName}/websites/{webSiteName}/slots` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/actionTemplates/{id}` | DELETE | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actionTemplates/{id}/v1` | DELETE | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates` | GET, POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/all` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/categories` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/search` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}` | GET, PUT | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/actionsUpdate` | POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/actionsUpdate/bulk` | POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/logo` | GET, POST, PUT | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/usage` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/v1` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/versions` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{id}/versions/{version}` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/actiontemplates/{typeOrId}/versions/{version}/logo` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/artifacts` | GET, POST | — |
| `/artifacts/{id}` | DELETE, GET, PUT | — |
| `/artifacts/{id}/content` | GET, PUT | — |
| `/audit-stream` | GET, PUT | — |
| `/authentication` | GET | — |
| `/authentication/checklogininitiated` | POST | — |
| `/azuredevopsissuetracker/connectivitycheck` | POST | — |
| `/build-information` | GET, POST | — |
| `/build-information/bulk` | DELETE | — |
| `/build-information/{id}` | DELETE, GET | — |
| `/capabilities` | GET | — |
| `/capabilities/{capability}` | GET | — |
| `/certificates` | GET, POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/all` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/certificate-global` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/generate` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}` | DELETE, GET, PUT | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/archive` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/archive/v1` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/export` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/replace` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/unarchive` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/unarchive/v1` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/certificates/{id}/usages` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/channels` | GET, POST | octopusdeploy_channel |
| `/channels/all` | GET | octopusdeploy_channel |
| `/channels/rule-test` | GET, POST | octopusdeploy_channel |
| `/channels/rule-test/v1` | GET, POST | octopusdeploy_channel |
| `/channels/{id}` | DELETE, GET, PUT | octopusdeploy_channel |
| `/channels/{id}/releases` | GET | octopusdeploy_channel |
| `/cloudtemplate/{id}/metadata` | POST | — |
| `/communityactiontemplates` | GET | octopusdeploy_community_step_template |
| `/communityactiontemplates/{id}` | GET | octopusdeploy_community_step_template |
| `/communityactiontemplates/{id}/actiontemplate` | GET | octopusdeploy_community_step_template |
| `/communityactiontemplates/{id}/actiontemplate/{actiontemplatespaceId}` | GET | octopusdeploy_community_step_template |
| `/communityactiontemplates/{id}/installation` | POST, PUT | octopusdeploy_community_step_template |
| `/communityactiontemplates/{id}/installation/{actiontemplatespaceId}` | POST, PUT | octopusdeploy_community_step_template |
| `/communityactiontemplates/{id}/logo` | GET | octopusdeploy_community_step_template |
| `/configuration` | GET | — |
| `/configuration/certificates` | GET | — |
| `/configuration/certificates/{id}` | GET | — |
| `/configuration/certificates/{id}/public-cer` | GET | — |
| `/configuration/retention-default` | GET, PUT | — |
| `/configuration/versioncontrol/clear-cache` | POST | — |
| `/configuration/versioncontrol/clear-cache/v1` | POST | — |
| `/configuration/{id}` | GET | — |
| `/configuration/{id}/metadata` | GET | — |
| `/configuration/{id}/values` | GET, PUT | — |
| `/dashboard` | GET | — |
| `/dashboard/dynamic` | GET | — |
| `/dashboardconfiguration` | GET, PUT | — |
| `/deploymentTargetTags/{tag}` | GET | — |
| `/deploymentfreezes` | GET, POST | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project, octopusdeploy_deployment_freeze_tenant (+1 more) |
| `/deploymentfreezes/{id}` | DELETE, GET, PUT | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project, octopusdeploy_deployment_freeze_tenant (+1 more) |
| `/deploymentprocesses` | GET | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/deploymentprocesses/{deploymentProcessId}/template` | GET | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/deploymentprocesses/{id}` | GET, PUT | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/deployments` | GET, POST | — |
| `/deployments/override` | POST | — |
| `/deployments/v1` | POST | — |
| `/deployments/{id}` | DELETE, GET | — |
| `/deploymentsettings/{id}` | GET | — |
| `/deploymentsettings/{projectId}` | PUT | — |
| `/deploymenttargettags` | GET, POST | — |
| `/deploymenttargettags/{tag}` | DELETE | — |
| `/deprecations/toggle` | POST | — |
| `/deprecations/toggle/v1` | POST | — |
| `/dynamic-extensions/features/metadata` | GET | — |
| `/dynamic-extensions/features/values` | GET, PUT | — |
| `/environments` | GET, POST | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/all` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/all/v1` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/sortorder` | PUT | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/summary` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/v1` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/{environmentId}` | PUT | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/{environmentId}/metadata` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/{environmentId}/singlyScopedVariableDetails` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/{id}` | DELETE, GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/environments/{id}/machines` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/events` | GET | — |
| `/events/agents` | GET | — |
| `/events/archives` | GET | — |
| `/events/archives/v1` | GET | — |
| `/events/archives/{fileName}` | DELETE, GET | — |
| `/events/archives/{fileName}/v1` | DELETE | — |
| `/events/categories` | GET | — |
| `/events/documenttypes` | GET | — |
| `/events/groups` | GET | — |
| `/events/{id}` | GET | — |
| `/externalgroups/directoryServices` | GET | — |
| `/externalgroups/ldap` | GET | — |
| `/externalsecuritygroupproviders` | GET | — |
| `/externalusers/directoryServices` | GET | octopusdeploy_user |
| `/externalusers/ldap` | GET | octopusdeploy_user |
| `/featuresconfiguration` | GET, PUT | — |
| `/feeds` | GET, POST | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/all` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/stats` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/{feedId}/packages` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/{feedId}/packages/notes` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/{id}` | DELETE, GET, PUT | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/{id}/packages/search` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/feeds/{id}/packages/versions` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/githubissuetracker/connectivitycheck` | POST | — |
| `/icons/all` | GET | — |
| `/icons/categories` | GET | — |
| `/insights/reports/{reportId}/logo` | GET | — |
| `/integrated-challenge` | GET | — |
| `/interruptions` | GET | — |
| `/interruptions/{id}` | GET | — |
| `/interruptions/{id}/responsible` | GET, PUT | — |
| `/interruptions/{id}/submit` | POST | — |
| `/jiraintegration/connectivitycheck/connectapp` | POST | — |
| `/jiraintegration/connectivitycheck/jira` | POST | — |
| `/jiraservicemanagement-integration/connectivity-test` | POST | — |
| `/letsencryptconfiguration` | GET, PUT | — |
| `/libraryvariablesets` | GET, POST | octopusdeploy_library_variable_set |
| `/libraryvariablesets/all` | GET | octopusdeploy_library_variable_set |
| `/libraryvariablesets/all/v1` | GET, POST | octopusdeploy_library_variable_set |
| `/libraryvariablesets/{id}` | DELETE, GET, PUT | octopusdeploy_library_variable_set |
| `/libraryvariablesets/{id}/usages` | GET | octopusdeploy_library_variable_set |
| `/licenses/licenses-current` | GET, PUT | — |
| `/licenses/licenses-current-features` | GET | — |
| `/licenses/licenses-current-status` | GET | — |
| `/licenses/licenses-current-usage` | GET | — |
| `/lifecycles` | GET, POST | octopusdeploy_lifecycle |
| `/lifecycles/all` | GET | octopusdeploy_lifecycle |
| `/lifecycles/previews` | GET | octopusdeploy_lifecycle |
| `/lifecycles/{id}` | DELETE, GET, PUT | octopusdeploy_lifecycle |
| `/lifecycles/{id}/preview` | GET | octopusdeploy_lifecycle |
| `/lifecycles/{id}/projects` | GET | octopusdeploy_lifecycle |
| `/machinepolicies` | GET, POST | octopusdeploy_machine_policy |
| `/machinepolicies/all` | GET | octopusdeploy_machine_policy |
| `/machinepolicies/template` | GET | octopusdeploy_machine_policy |
| `/machinepolicies/{id}` | DELETE, GET, PUT | octopusdeploy_machine_policy |
| `/machinepolicies/{id}/machines` | GET | octopusdeploy_machine_policy |
| `/machinepolicies/{id}/v1` | DELETE | octopusdeploy_machine_policy |
| `/machinepolicies/{id}/workers` | GET | octopusdeploy_machine_policy |
| `/machineroles/all` | GET | — |
| `/machineroles/all/v1` | GET | — |
| `/machines` | GET, POST | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/all/v1` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/discover` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/operatingsystem/names/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/operatingsystem/shells/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{id}` | DELETE, GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{id}/connection` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{id}/latestdeployments` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{id}/tasks` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{id}/tasks/v1` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{machineId}/singlyScopedVariableDetails` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/machines/{machineid}` | PUT | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/maintenanceconfiguration` | GET, PUT | — |
| `/maintenanceconfiguration/v1` | PUT | — |
| `/migrations/import` | POST | — |
| `/migrations/partialexport` | POST | — |
| `/nuget/packages` | PUT | octopusdeploy_nuget_feed |
| `/octopusservernodes` | GET | — |
| `/octopusservernodes/all` | GET | — |
| `/octopusservernodes/ping` | GET | — |
| `/octopusservernodes/summary` | GET | — |
| `/octopusservernodes/{id}` | DELETE, GET, PUT | — |
| `/octopusservernodes/{id}/details` | GET | — |
| `/packages` | GET | — |
| `/packages/bulk` | DELETE | — |
| `/packages/bulk/v1` | DELETE | — |
| `/packages/notes` | GET | — |
| `/packages/raw` | POST | — |
| `/packages/{id}` | DELETE, GET | — |
| `/packages/{id}/raw` | GET | — |
| `/packages/{id}/v1` | DELETE | — |
| `/packages/{packageId}/{baseVersion}/delta` | POST | — |
| `/packages/{packageId}/{version}/delta-signature` | GET | — |
| `/performanceconfiguration` | GET, PUT | — |
| `/permissions/all` | GET | — |
| `/platformhub/accounts` | GET, POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/accounts/{id}` | DELETE, GET, PUT | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/git-credentials` | GET, POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/git-credentials/{id}` | DELETE, GET, PUT | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/git/branches` | GET, POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/git/tags` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/policies/{slug}/versions` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/policies/{slug}/versions/{version}/modify-status` | POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/processtemplates/{slug}/share` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/processtemplates/{slug}/versions` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/processtemplates/{slug}/{versionMask}` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/versioncontrol` | GET, PUT | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/policies` | GET, POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/policies/{slug}` | GET, PUT | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/policies/{slug}/publish` | POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/processtemplates` | GET, POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/processtemplates/{slug}` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/processtemplates/{slug}/share` | POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/processtemplates/{slug}/variables/names` | GET | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/processtemplates/{slug}/versions` | POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitRef}/projecttemplates/{slug}/share` | POST | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/platformhub/{gitref}/processtemplates/{slug}` | DELETE, PUT | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account (+7 more) |
| `/progression/runbooks/taskRuns` | GET | — |
| `/progression/runbooks/{runbookId}` | GET | — |
| `/progression/runbooks/{runbookId}/v1` | GET | — |
| `/progression/{projectId}` | GET | — |
| `/projectgroups` | GET, POST | octopusdeploy_project_group |
| `/projectgroups/all` | GET | octopusdeploy_project_group |
| `/projectgroups/{id}` | DELETE, GET, PUT | octopusdeploy_project_group |
| `/projectgroups/{id}/projects` | GET | octopusdeploy_project_group |
| `/projects` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/all` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{id}/logo` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/channels` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/channels/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/channels/{id}/releases` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/channels/{id}/v2` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/deploymentprocesses` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/deploymentprocesses/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/deploymentprocesses/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/deploymentsettings` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/branches` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/branches/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/branches/{branchName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/commits/{hash}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/connectivity-test` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/convert` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/refs/{refName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/tags` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/tags/{tagName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/git/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/logo` | POST, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/metadata` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/progression` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/progression/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/releases` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/releases/{id}/variables` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/releases/{version}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookProcesses` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookProcesses/{id}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookRuns` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookRuns/{id}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{idOrName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}` | PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/runbookRuns` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookSnapshots/{id}/variables` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookruns/{id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbookruns/{runbookRunId}/retry/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/all` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/all/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{id}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{id}/runbookSnapshots` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{runbookId}/run` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{runbookId}/runbookRuns/previews` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/runbooksnapshots/{id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/summary` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/summary/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/triggers` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/triggers/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/deploymentprocesses` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/deploymentprocesses/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/deploymentprocesses/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/deploymentsettings` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{id}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{runbookId}/run/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/runbooks/{runbookId}/runbookRuns/previews` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/summary` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitRef}/summary/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{gitref}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projects/{projectId}/{unusedGitRef}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/projecttriggers` | GET, POST | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger (+2 more) |
| `/projecttriggers/{id}` | DELETE, GET, PUT | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger (+2 more) |
| `/proxies` | GET, POST | octopusdeploy_machine_proxy |
| `/proxies/all` | GET | octopusdeploy_machine_proxy |
| `/proxies/{id}` | DELETE, GET, PUT | octopusdeploy_machine_proxy |
| `/releases` | GET, POST | — |
| `/releases/create/v1` | POST | — |
| `/releases/{id}` | DELETE, GET, PUT | — |
| `/releases/{id}/deployments/template` | GET | — |
| `/releases/{releaseId}/defects` | GET, POST | — |
| `/releases/{releaseId}/defects/resolve` | POST | — |
| `/releases/{releaseId}/deployments` | GET | — |
| `/releases/{releaseId}/deployments/preview/{environmentId}` | GET | — |
| `/releases/{releaseId}/deployments/preview/{environmentId}/{tenantId}` | GET | — |
| `/releases/{releaseId}/deployments/previews` | POST | — |
| `/releases/{releaseId}/progression` | GET | — |
| `/releases/{releaseId}/snapshot-variables` | POST | — |
| `/reporting/deployments-counted-by-week` | GET | — |
| `/reporting/deployments/xml` | GET | — |
| `/retentionpolicies` | GET | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy |
| `/retentionpolicies/{id}` | PUT | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy |
| `/runbook-runs/create/v1` | POST | octopusdeploy_runbook |
| `/runbookProcesses` | GET | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/runbookProcesses/{id}` | GET, PUT | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/runbookProcesses/{id}/runbookSnapshotTemplate` | GET | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/runbookRuns` | GET, POST | octopusdeploy_runbook |
| `/runbookRuns/{id}` | GET | octopusdeploy_runbook |
| `/runbookSnapshots` | GET, POST | octopusdeploy_runbook |
| `/runbookSnapshots/{id}` | GET, PUT | octopusdeploy_runbook |
| `/runbookSnapshots/{id}/runbookRuns` | GET | octopusdeploy_runbook |
| `/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | GET | octopusdeploy_runbook |
| `/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | GET | octopusdeploy_runbook |
| `/runbookSnapshots/{id}/runbookRuns/template` | GET | octopusdeploy_runbook |
| `/runbookSnapshots/{id}/snapshot-variables` | POST | octopusdeploy_runbook |
| `/runbookruns/{id}` | DELETE | octopusdeploy_runbook |
| `/runbooks` | GET, POST | octopusdeploy_runbook |
| `/runbooks/all` | GET | octopusdeploy_runbook |
| `/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_runbook |
| `/runbooks/{id}/environments` | GET | octopusdeploy_runbook |
| `/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_runbook |
| `/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_runbook |
| `/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_runbook |
| `/runbooks/{id}/runbookSnapshots` | GET | octopusdeploy_runbook |
| `/runbooks/{runbookId}/run` | POST | octopusdeploy_runbook |
| `/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_runbook |
| `/runbooksnapshots/{id}` | DELETE | octopusdeploy_runbook |
| `/scheduler` | GET | — |
| `/scheduler/start` | GET | — |
| `/scheduler/stop` | GET | — |
| `/scheduler/trigger` | GET | — |
| `/scheduler/{name}/logs` | GET | — |
| `/scheduler/{name}/logs/raw` | GET | — |
| `/scopeduserroles` | GET, POST | octopusdeploy_scoped_user_role |
| `/scopeduserroles/{id}` | DELETE, GET, PUT | octopusdeploy_scoped_user_role |
| `/serverconfiguration` | GET, PUT | — |
| `/serverconfiguration/settings` | GET | — |
| `/serverstatus` | GET | — |
| `/serverstatus/counts` | GET | — |
| `/serverstatus/gc-collect` | POST | — |
| `/serverstatus/gc-collect/v1` | POST | — |
| `/serverstatus/health` | GET | — |
| `/serverstatus/logs` | GET | — |
| `/serverstatus/system-info` | GET | — |
| `/serverstatus/system-report` | GET | — |
| `/serverstatus/timezones` | GET | — |
| `/serviceaccounts/{serviceAccountId}/oidcidentities/create/v1` | POST | — |
| `/serviceaccounts/{serviceAccountId}/oidcidentities/v1` | GET | — |
| `/serviceaccounts/{serviceAccountId}/oidcidentities/{id}/v1` | DELETE, GET, PUT | — |
| `/servicenow-integration/connectivity-test` | POST | — |
| `/signingkeyconfiguration` | GET, PUT | — |
| `/smtpconfiguration` | GET, PUT | — |
| `/smtpconfiguration/isconfigured` | GET | — |
| `/smtpconfiguration/isconfigured/v1` | GET | — |
| `/smtpconfiguration/v1` | GET | — |
| `/spaces` | GET, POST | octopusdeploy_space |
| `/spaces/all` | GET | octopusdeploy_space |
| `/spaces/v1` | POST | octopusdeploy_space |
| `/spaces/{id}` | DELETE, GET, PUT | octopusdeploy_space |
| `/spaces/{id}/logo` | GET, POST, PUT | — |
| `/spaces/{id}/search` | GET | — |
| `/spaces/{id}/v1` | DELETE | — |
| `/spaces/{spaceIdentifier}/accounts` | GET, POST | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/all` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{accountId}` | PUT | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}` | DELETE, GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/pk` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/resourceGroups` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/storageAccounts` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/usages` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/v1` | DELETE | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/websites` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/accounts/{id}/{resourceGroupName}/websites/{webSiteName}/slots` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/spaces/{spaceIdentifier}/actionTemplates/{id}` | DELETE | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actionTemplates/{id}/v1` | DELETE | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates` | GET, POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/all` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/categories` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/search` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}` | GET, PUT | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/actionsUpdate` | POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/actionsUpdate/bulk` | POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/logo` | GET, POST, PUT | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/usage` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/v1` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/versions` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{id}/versions/{version}` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/actiontemplates/{typeOrId}/versions/{version}/logo` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/spaces/{spaceIdentifier}/artifacts` | GET, POST | — |
| `/spaces/{spaceIdentifier}/artifacts/{id}` | DELETE, GET, PUT | — |
| `/spaces/{spaceIdentifier}/artifacts/{id}/content` | GET, PUT | — |
| `/spaces/{spaceIdentifier}/build-information` | GET, POST | — |
| `/spaces/{spaceIdentifier}/build-information/bulk` | DELETE | — |
| `/spaces/{spaceIdentifier}/build-information/{id}` | DELETE, GET | — |
| `/spaces/{spaceIdentifier}/certificates` | GET, POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/all` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/generate` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}` | DELETE, GET, PUT | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/archive` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/archive/v1` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/export` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/replace` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/unarchive` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/unarchive/v1` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/certificates/{id}/usages` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/spaces/{spaceIdentifier}/channels` | GET, POST | octopusdeploy_channel |
| `/spaces/{spaceIdentifier}/channels/all` | GET | octopusdeploy_channel |
| `/spaces/{spaceIdentifier}/channels/rule-test` | GET, POST | octopusdeploy_channel |
| `/spaces/{spaceIdentifier}/channels/rule-test/v1` | GET, POST | octopusdeploy_channel |
| `/spaces/{spaceIdentifier}/channels/{id}` | DELETE, GET, PUT | octopusdeploy_channel |
| `/spaces/{spaceIdentifier}/channels/{id}/releases` | GET | octopusdeploy_channel |
| `/spaces/{spaceIdentifier}/dashboard` | GET | — |
| `/spaces/{spaceIdentifier}/dashboard/dynamic` | GET | — |
| `/spaces/{spaceIdentifier}/dashboardconfiguration` | GET, PUT | — |
| `/spaces/{spaceIdentifier}/deploymentTargetTags/{tag}` | GET | — |
| `/spaces/{spaceIdentifier}/deploymentprocesses` | GET | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/spaces/{spaceIdentifier}/deploymentprocesses/{deploymentProcessId}/template` | GET | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/spaces/{spaceIdentifier}/deploymentprocesses/{id}` | GET, PUT | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/spaces/{spaceIdentifier}/deployments` | GET, POST | — |
| `/spaces/{spaceIdentifier}/deployments/create/tenanted/v1` | POST | — |
| `/spaces/{spaceIdentifier}/deployments/create/untenanted/v1` | POST | — |
| `/spaces/{spaceIdentifier}/deployments/override` | POST | — |
| `/spaces/{spaceIdentifier}/deployments/v1` | POST | — |
| `/spaces/{spaceIdentifier}/deployments/{id}` | DELETE, GET | — |
| `/spaces/{spaceIdentifier}/deploymentsettings/{id}` | GET | — |
| `/spaces/{spaceIdentifier}/deploymentsettings/{projectId}` | PUT | — |
| `/spaces/{spaceIdentifier}/deploymenttargettags` | GET, POST | — |
| `/spaces/{spaceIdentifier}/deploymenttargettags/{tag}` | DELETE | — |
| `/spaces/{spaceIdentifier}/environments` | GET, POST | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/all` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/all/v1` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/ephemeral/{id}/deprovision` | POST | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/sortorder` | PUT | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/summary` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/summary/v2` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/v1` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/v2` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/{environmentId}` | PUT | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/{environmentId}/metadata` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/{environmentId}/singlyScopedVariableDetails` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/{id}` | DELETE, GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/{id}/machines` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/environments/{id}/v2` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/events` | GET | — |
| `/spaces/{spaceIdentifier}/events/agents` | GET | — |
| `/spaces/{spaceIdentifier}/events/categories` | GET | — |
| `/spaces/{spaceIdentifier}/events/documenttypes` | GET | — |
| `/spaces/{spaceIdentifier}/events/groups` | GET | — |
| `/spaces/{spaceIdentifier}/events/{id}` | GET | — |
| `/spaces/{spaceIdentifier}/feeds` | GET, POST | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/all` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/stats` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/{feedId}/packages` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/{feedId}/packages/notes` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/{id}` | DELETE, GET, PUT | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/{id}/packages/search` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/feeds/{id}/packages/versions` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/spaces/{spaceIdentifier}/git-credentials` | GET, POST | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/spaces/{spaceIdentifier}/git-credentials/v1` | GET, POST | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/spaces/{spaceIdentifier}/git-credentials/{id}` | DELETE, GET, PUT | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/spaces/{spaceIdentifier}/git-credentials/{id}/usage` | GET | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/spaces/{spaceIdentifier}/git-credentials/{id}/usage/v1` | GET | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/spaces/{spaceIdentifier}/git-credentials/{id}/v1` | DELETE, GET, PUT | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/spaces/{spaceIdentifier}/insights/reports` | GET, POST | — |
| `/spaces/{spaceIdentifier}/insights/reports/v1` | POST | — |
| `/spaces/{spaceIdentifier}/insights/reports/{id}` | DELETE, GET, PUT | — |
| `/spaces/{spaceIdentifier}/insights/reports/{id}/v1` | DELETE, GET, PUT | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/deployments` | GET | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/deployments/csv` | GET | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/logo` | GET, POST | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/logo/icon` | POST | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/logo/icon/v1` | POST | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/metrics` | GET | — |
| `/spaces/{spaceIdentifier}/insights/reports/{reportId}/metrics/v1` | GET | — |
| `/spaces/{spaceIdentifier}/interruptions` | GET | — |
| `/spaces/{spaceIdentifier}/interruptions/{id}` | GET | — |
| `/spaces/{spaceIdentifier}/interruptions/{id}/responsible` | GET, PUT | — |
| `/spaces/{spaceIdentifier}/interruptions/{id}/submit` | POST | — |
| `/spaces/{spaceIdentifier}/libraryvariablesets` | GET, POST | octopusdeploy_library_variable_set |
| `/spaces/{spaceIdentifier}/libraryvariablesets/all` | GET | octopusdeploy_library_variable_set |
| `/spaces/{spaceIdentifier}/libraryvariablesets/all/v1` | GET, POST | octopusdeploy_library_variable_set |
| `/spaces/{spaceIdentifier}/libraryvariablesets/{id}` | DELETE, GET, PUT | octopusdeploy_library_variable_set |
| `/spaces/{spaceIdentifier}/libraryvariablesets/{id}/usages` | GET | octopusdeploy_library_variable_set |
| `/spaces/{spaceIdentifier}/lifecycles` | GET, POST | octopusdeploy_lifecycle |
| `/spaces/{spaceIdentifier}/lifecycles/all` | GET | octopusdeploy_lifecycle |
| `/spaces/{spaceIdentifier}/lifecycles/previews` | GET | octopusdeploy_lifecycle |
| `/spaces/{spaceIdentifier}/lifecycles/{id}` | DELETE, GET, PUT | octopusdeploy_lifecycle |
| `/spaces/{spaceIdentifier}/lifecycles/{id}/preview` | GET | octopusdeploy_lifecycle |
| `/spaces/{spaceIdentifier}/lifecycles/{id}/projects` | GET | octopusdeploy_lifecycle |
| `/spaces/{spaceIdentifier}/machinepolicies` | GET, POST | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machinepolicies/all` | GET | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machinepolicies/template` | GET | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machinepolicies/{id}` | DELETE, GET, PUT | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machinepolicies/{id}/machines` | GET | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machinepolicies/{id}/v1` | DELETE | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machinepolicies/{id}/workers` | GET | octopusdeploy_machine_policy |
| `/spaces/{spaceIdentifier}/machineroles/all` | GET | — |
| `/spaces/{spaceIdentifier}/machineroles/all/v1` | GET | — |
| `/spaces/{spaceIdentifier}/machines` | GET, POST | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/all/v1` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/discover` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/operatingsystem/names/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/operatingsystem/shells/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{id}` | DELETE, GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{id}/connection` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{id}/latestdeployments` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{id}/tasks` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{id}/tasks/v1` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{machineId}/singlyScopedVariableDetails` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/machines/{machineid}` | PUT | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/spaces/{spaceIdentifier}/observability/agents` | POST | — |
| `/spaces/{spaceIdentifier}/observability/events/sessions` | POST | — |
| `/spaces/{spaceIdentifier}/observability/events/sessions/{sessionId}` | GET | — |
| `/spaces/{spaceIdentifier}/observability/kubernetes-monitors` | POST | — |
| `/spaces/{spaceIdentifier}/observability/kubernetes-monitors/{id}` | DELETE, GET | — |
| `/spaces/{spaceIdentifier}/observability/logs/sessions` | POST | — |
| `/spaces/{spaceIdentifier}/observability/logs/sessions/{sessionId}` | GET | — |
| `/spaces/{spaceIdentifier}/packages` | GET | — |
| `/spaces/{spaceIdentifier}/packages/bulk` | DELETE | — |
| `/spaces/{spaceIdentifier}/packages/bulk/v1` | DELETE | — |
| `/spaces/{spaceIdentifier}/packages/notes` | GET | — |
| `/spaces/{spaceIdentifier}/packages/raw` | POST | — |
| `/spaces/{spaceIdentifier}/packages/{id}` | DELETE, GET | — |
| `/spaces/{spaceIdentifier}/packages/{id}/raw` | GET | — |
| `/spaces/{spaceIdentifier}/packages/{id}/v1` | DELETE | — |
| `/spaces/{spaceIdentifier}/packages/{packageId}/{baseVersion}/delta` | POST | — |
| `/spaces/{spaceIdentifier}/packages/{packageId}/{version}/delta-signature` | GET | — |
| `/spaces/{spaceIdentifier}/parentEnvironments` | POST | octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/parentEnvironments/{environmentId}` | PUT | octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/parentEnvironments/{id}` | DELETE, GET | octopusdeploy_parent_environment |
| `/spaces/{spaceIdentifier}/processtemplates/{slug}/{versionMask}` | GET | octopusdeploy_process |
| `/spaces/{spaceIdentifier}/processtemplates/{slug}/{versionMask}/icon` | GET | octopusdeploy_process |
| `/spaces/{spaceIdentifier}/progression/runbooks/taskRuns` | GET | — |
| `/spaces/{spaceIdentifier}/progression/runbooks/{runbookId}` | GET | — |
| `/spaces/{spaceIdentifier}/progression/runbooks/{runbookId}/v1` | GET | — |
| `/spaces/{spaceIdentifier}/progression/{projectId}` | GET | — |
| `/spaces/{spaceIdentifier}/projectgroups` | GET, POST | octopusdeploy_project_group |
| `/spaces/{spaceIdentifier}/projectgroups/all` | GET | octopusdeploy_project_group |
| `/spaces/{spaceIdentifier}/projectgroups/{id}` | DELETE, GET, PUT | octopusdeploy_project_group |
| `/spaces/{spaceIdentifier}/projectgroups/{id}/projects` | GET | octopusdeploy_project_group |
| `/spaces/{spaceIdentifier}/projects` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/all` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{id}/logo` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/channels` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/channels/{channelId}/git-reference-rule-validation/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/channels/{channelId}/git-resource-rule-validation/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/channels/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/channels/{id}/releases` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/channels/{id}/v2` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/deploymentprocesses` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/deploymentprocesses/resolved` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/deploymentprocesses/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/deploymentprocesses/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/deploymentsettings` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral/{environmentId}/deprovision` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral/{environmentId}/deprovisioning/mark-successful` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral/{environmentId}/deprovisioning/retry` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral/{environmentId}/provisioning/mark-successful` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral/{environmentId}/provisioning/retry` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/ephemeral/{id}/status` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/{environmentId}/tenants/{tenantId}/livestatus` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/{environmentId}/tenants/{tenantId}/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/{environmentId}/tenants/{tenantId}/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}/manifest` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/{environmentId}/untenanted/livestatus` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/{environmentId}/untenanted/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/environments/{environmentId}/untenanted/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}/manifest` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/featuretoggles` | GET, POST, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/featuretoggles/generate-client-identifier` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/featuretoggles/rotate-client-identifier-signing-key` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/featuretoggles/{Id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/featuretoggles/{slug}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/branches` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/branches/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/branches/{branchName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/commits/{hash}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/connectivity-test` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/convert` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/migrate-runbooks` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/migrate-variables` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/refs` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/refs/{refName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/tags` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/tags/{tagName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/git/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/insights/deployments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/insights/deployments/csv` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/insights/metrics` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/insights/metrics/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/logo` | POST, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/logo/icon` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/metadata` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/progression` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/progression/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/releases` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/releases/{id}/variables` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/releases/{version}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookProcesses` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookProcesses/{id}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookRuns` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookRuns/{id}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookRuns/{id}/details/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{idOrName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}` | PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookSnapshots/{id}/variables` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookruns/{id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbookruns/{runbookRunId}/retry/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/all` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/all/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}/environments/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{id}/runbookSnapshots` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{runbookId}/run` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{runbookId}/runbookRuns/previews` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/runbooksnapshots/{id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/summary` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/summary/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/triggers` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/triggers/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/variables` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/deploymentprocesses` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/deploymentprocesses/resolved` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/deploymentprocesses/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/deploymentprocesses/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/deploymentsettings` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{id}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{id}/environments/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{runbookId}/run/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/runbooks/{runbookId}/runbookRuns/previews` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/summary` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/summary/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitRef}/variables` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{gitref}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projects/{projectId}/{unusedGitRef}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/spaces/{spaceIdentifier}/projecttriggers` | GET, POST | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger (+2 more) |
| `/spaces/{spaceIdentifier}/projecttriggers/{id}` | DELETE, GET, PUT | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger (+2 more) |
| `/spaces/{spaceIdentifier}/proxies` | GET, POST | octopusdeploy_machine_proxy |
| `/spaces/{spaceIdentifier}/proxies/all` | GET | octopusdeploy_machine_proxy |
| `/spaces/{spaceIdentifier}/proxies/{id}` | DELETE, GET, PUT | octopusdeploy_machine_proxy |
| `/spaces/{spaceIdentifier}/releases` | GET, POST | — |
| `/spaces/{spaceIdentifier}/releases/create/v1` | POST | — |
| `/spaces/{spaceIdentifier}/releases/{id}` | DELETE, GET, PUT | — |
| `/spaces/{spaceIdentifier}/releases/{id}/deployments/template` | GET | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/defects` | GET, POST | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/defects/resolve` | POST | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/deployments` | GET | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/deployments/preview/{environmentId}` | GET | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/deployments/preview/{environmentId}/{tenantId}` | GET | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/deployments/previews` | POST | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/missingPackages` | GET | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/progression` | GET | — |
| `/spaces/{spaceIdentifier}/releases/{releaseId}/snapshot-variables` | POST | — |
| `/spaces/{spaceIdentifier}/reporting/deployments-counted-by-week` | GET | — |
| `/spaces/{spaceIdentifier}/reporting/deployments/xml` | GET | — |
| `/spaces/{spaceIdentifier}/retentionpolicies` | GET | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy |
| `/spaces/{spaceIdentifier}/retentionpolicies/{id}` | PUT | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy |
| `/spaces/{spaceIdentifier}/runbook-runs/create/v1` | POST | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookProcesses` | GET | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/spaces/{spaceIdentifier}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/spaces/{spaceIdentifier}/runbookProcesses/{id}/runbookSnapshotTemplate` | GET | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/spaces/{spaceIdentifier}/runbookRuns` | GET, POST | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookRuns/{id}` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots` | GET, POST | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots/{id}` | GET, PUT | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots/{id}/runbookRuns` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots/{id}/runbookRuns/template` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookSnapshots/{id}/snapshot-variables` | POST | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbookruns/{id}` | DELETE | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks` | GET, POST | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/all` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{id}/environments` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{id}/runbookSnapshots` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{runbookId}/run` | POST | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/runbooksnapshots/{id}` | DELETE | octopusdeploy_runbook |
| `/spaces/{spaceIdentifier}/scopeduserroles` | GET, POST | octopusdeploy_scoped_user_role |
| `/spaces/{spaceIdentifier}/scopeduserroles/{id}` | DELETE, GET, PUT | octopusdeploy_scoped_user_role |
| `/spaces/{spaceIdentifier}/spaces/{id}/search` | GET | octopusdeploy_space |
| `/spaces/{spaceIdentifier}/subscriptions` | GET, POST | — |
| `/spaces/{spaceIdentifier}/subscriptions/all` | GET | — |
| `/spaces/{spaceIdentifier}/subscriptions/{id}` | DELETE, GET, PUT | — |
| `/spaces/{spaceIdentifier}/tagsets` | GET, POST | octopusdeploy_tag, octopusdeploy_tag_set |
| `/spaces/{spaceIdentifier}/tagsets/all` | GET | octopusdeploy_tag, octopusdeploy_tag_set |
| `/spaces/{spaceIdentifier}/tagsets/sortorder` | PUT | octopusdeploy_tag, octopusdeploy_tag_set |
| `/spaces/{spaceIdentifier}/tagsets/{id}` | DELETE, GET, PUT | octopusdeploy_tag, octopusdeploy_tag_set |
| `/spaces/{spaceIdentifier}/tasks` | GET, POST | — |
| `/spaces/{spaceIdentifier}/tasks/rerun/{id}` | POST | — |
| `/spaces/{spaceIdentifier}/tasks/tasktypes` | GET | — |
| `/spaces/{spaceIdentifier}/tasks/{id}` | GET | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/cancel` | POST | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/details` | GET | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/prioritize` | POST | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/queued-behind` | GET | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/raw` | GET | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/state` | POST | — |
| `/spaces/{spaceIdentifier}/tasks/{id}/status/messages` | GET | — |
| `/spaces/{spaceIdentifier}/teammembership` | GET | octopusdeploy_team |
| `/spaces/{spaceIdentifier}/teammembership/previewteam` | POST | octopusdeploy_team |
| `/spaces/{spaceIdentifier}/teams` | GET, POST | octopusdeploy_team |
| `/spaces/{spaceIdentifier}/teams/all` | GET | octopusdeploy_team |
| `/spaces/{spaceIdentifier}/teams/{id}` | DELETE, GET, PUT | octopusdeploy_team |
| `/spaces/{spaceIdentifier}/teams/{id}/scopeduserroles` | GET | octopusdeploy_team |
| `/spaces/{spaceIdentifier}/tenants` | GET, POST | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/all` | GET | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/status` | GET | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/tag-test` | GET | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/variables-missing` | GET | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/{id}` | DELETE, GET, PUT | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/{id}/logo` | GET, POST, PUT | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/{id}/variables` | GET, POST, PUT | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/{tenantId}/commonvariables` | GET, POST, PUT | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/{tenantId}/logo/icon` | POST | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenants/{tenantId}/projectvariables` | GET, POST, PUT | octopusdeploy_tenant |
| `/spaces/{spaceIdentifier}/tenantvariables/all` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable |
| `/spaces/{spaceIdentifier}/users/invitations` | POST | octopusdeploy_user |
| `/spaces/{spaceIdentifier}/users/invitations/{id}` | GET | octopusdeploy_user |
| `/spaces/{spaceIdentifier}/users/{id}/permissions` | GET | octopusdeploy_user |
| `/spaces/{spaceIdentifier}/users/{id}/permissions/configuration` | GET | octopusdeploy_user |
| `/spaces/{spaceIdentifier}/users/{id}/permissions/export` | GET | octopusdeploy_user |
| `/spaces/{spaceIdentifier}/users/{userId}/teams` | GET | octopusdeploy_user |
| `/spaces/{spaceIdentifier}/variables/all` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/spaces/{spaceIdentifier}/variables/names` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/spaces/{spaceIdentifier}/variables/preview` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/spaces/{spaceIdentifier}/variables/{id}` | GET, PUT | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/spaces/{spaceIdentifier}/workerpools` | GET, POST | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/dynamicworkertypes` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/sortorder` | PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/summary` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/supportedtypes` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/{id}` | DELETE, GET, PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workerpools/{id}/workers` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/spaces/{spaceIdentifier}/workers` | GET, POST | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workers/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workers/discover` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workers/operatingsystem/names/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workers/operatingsystem/shells/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workers/{id}` | DELETE, GET, PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workers/{id}/connection` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/spaces/{spaceIdentifier}/workertaskleases` | GET | — |
| `/spaces/{spaceId}/logo/icon` | POST | — |
| `/subscriptions` | GET, POST | — |
| `/subscriptions/all` | GET | — |
| `/subscriptions/{id}` | DELETE, GET, PUT | — |
| `/tagsets` | GET, POST | octopusdeploy_tag, octopusdeploy_tag_set |
| `/tagsets/all` | GET | octopusdeploy_tag, octopusdeploy_tag_set |
| `/tagsets/sortorder` | PUT | octopusdeploy_tag, octopusdeploy_tag_set |
| `/tagsets/{id}` | DELETE, GET, PUT | octopusdeploy_tag, octopusdeploy_tag_set |
| `/tasks` | GET, POST | — |
| `/tasks/rerun/{id}` | POST | — |
| `/tasks/tasktypes` | GET | — |
| `/tasks/{id}` | GET | — |
| `/tasks/{id}/cancel` | POST | — |
| `/tasks/{id}/details` | GET | — |
| `/tasks/{id}/prioritize` | POST | — |
| `/tasks/{id}/queued-behind` | GET | — |
| `/tasks/{id}/raw` | GET | — |
| `/tasks/{id}/state` | POST | — |
| `/tasks/{id}/status/messages` | GET | — |
| `/teammembership` | GET | octopusdeploy_team |
| `/teammembership/previewteam` | POST | octopusdeploy_team |
| `/teams` | GET, POST | octopusdeploy_team |
| `/teams/all` | GET | octopusdeploy_team |
| `/teams/{id}` | DELETE, GET, PUT | octopusdeploy_team |
| `/teams/{id}/scopeduserroles` | GET | octopusdeploy_team |
| `/telemetry/download` | GET | — |
| `/telemetry/lastTask` | GET | — |
| `/telemetryconfiguration` | GET, PUT | — |
| `/tenants` | GET, POST | octopusdeploy_tenant |
| `/tenants/all` | GET | octopusdeploy_tenant |
| `/tenants/status` | GET | octopusdeploy_tenant |
| `/tenants/tag-test` | GET | octopusdeploy_tenant |
| `/tenants/variables-missing` | GET | octopusdeploy_tenant |
| `/tenants/{id}` | DELETE, GET, PUT | octopusdeploy_tenant |
| `/tenants/{id}/logo` | GET, POST, PUT | octopusdeploy_tenant |
| `/tenants/{id}/variables` | GET, POST, PUT | octopusdeploy_tenant |
| `/tenants/{tenantId}/commonvariables` | GET, POST, PUT | octopusdeploy_tenant |
| `/tenants/{tenantId}/projectvariables` | GET, POST, PUT | octopusdeploy_tenant |
| `/tenantvariables/all` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable |
| `/token/v1` | POST | octopusdeploy_token_account |
| `/upgradeconfiguration` | GET, PUT | — |
| `/userroles` | GET, POST | octopusdeploy_scoped_user_role, octopusdeploy_user_role |
| `/userroles/all` | GET | octopusdeploy_scoped_user_role, octopusdeploy_user_role |
| `/userroles/{id}` | DELETE, GET, PUT | octopusdeploy_scoped_user_role, octopusdeploy_user_role |
| `/users` | GET, POST | octopusdeploy_user |
| `/users/access-token` | POST | octopusdeploy_user |
| `/users/all` | GET | octopusdeploy_user |
| `/users/authenticate/AzureAD` | POST | octopusdeploy_user |
| `/users/authenticate/GenericOidc` | POST | octopusdeploy_user |
| `/users/authenticate/GoogleApps` | POST | octopusdeploy_user |
| `/users/authenticate/OctopusID` | POST | octopusdeploy_user |
| `/users/authenticate/Okta` | POST | octopusdeploy_user |
| `/users/authenticatedToken/AzureAD` | GET, POST | octopusdeploy_user |
| `/users/authenticatedToken/GenericOidc` | GET, POST | octopusdeploy_user |
| `/users/authenticatedToken/GoogleApps` | GET, POST | octopusdeploy_user |
| `/users/authenticatedToken/OctopusID` | GET, POST | octopusdeploy_user |
| `/users/authenticatedToken/Okta` | GET, POST | octopusdeploy_user |
| `/users/authentication` | GET | octopusdeploy_user |
| `/users/authentication/{userId}` | GET | octopusdeploy_user |
| `/users/external-search` | GET | octopusdeploy_user |
| `/users/identity-metadata` | GET | octopusdeploy_user |
| `/users/invitations` | POST | octopusdeploy_user |
| `/users/invitations/{id}` | GET | octopusdeploy_user |
| `/users/login` | POST | octopusdeploy_user |
| `/users/logout` | POST | octopusdeploy_user |
| `/users/me` | GET | octopusdeploy_user |
| `/users/register` | POST | octopusdeploy_user |
| `/users/{id}` | DELETE, GET, PUT | octopusdeploy_user |
| `/users/{id}/permissions` | GET | octopusdeploy_user |
| `/users/{id}/permissions/configuration` | GET | octopusdeploy_user |
| `/users/{id}/permissions/export` | GET | octopusdeploy_user |
| `/users/{id}/spaces` | GET | octopusdeploy_user |
| `/users/{userId}/apikeys` | GET, POST | octopusdeploy_user |
| `/users/{userId}/apikeys/v1` | GET | octopusdeploy_user |
| `/users/{userId}/apikeys/{id}` | DELETE, GET | octopusdeploy_user |
| `/users/{userId}/revoke-sessions` | PUT | octopusdeploy_user |
| `/users/{userId}/teams` | GET | octopusdeploy_user |
| `/variables/all` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/variables/names` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/variables/preview` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/variables/{id}` | GET, PUT | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/workerpools` | GET, POST | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/dynamicworkertypes` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/sortorder` | PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/summary` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/supportedtypes` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/{id}` | DELETE, GET, PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workerpools/{id}/workers` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/workers` | GET, POST | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/workers/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/workers/discover` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/workers/operatingsystem/names/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/workers/operatingsystem/shells/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/workers/{id}` | DELETE, GET, PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/workers/{id}/connection` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}` | GET | — |
| `/{spaceId}/accounts` | GET, POST | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/all` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{accountId}` | PUT | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}` | DELETE, GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/pk` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/resourceGroups` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/storageAccounts` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/usages` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/v1` | DELETE | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/websites` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/accounts/{id}/{resourceGroupName}/websites/{webSiteName}/slots` | GET | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account (+11 more) |
| `/{spaceId}/actionTemplates/{id}` | DELETE | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actionTemplates/{id}/v1` | DELETE | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates` | GET, POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/all` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/categories` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/search` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}` | GET, PUT | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/actionsUpdate` | POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/actionsUpdate/bulk` | POST | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/logo` | GET, POST, PUT | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/usage` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/v1` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/versions` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{id}/versions/{version}` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/actiontemplates/{typeOrId}/versions/{version}/logo` | GET | octopusdeploy_community_step_template, octopusdeploy_step_template |
| `/{spaceId}/artifacts` | GET, POST | — |
| `/{spaceId}/artifacts/{id}` | DELETE, GET, PUT | — |
| `/{spaceId}/artifacts/{id}/content` | GET, PUT | — |
| `/{spaceId}/build-information` | GET, POST | — |
| `/{spaceId}/build-information/bulk` | DELETE | — |
| `/{spaceId}/build-information/{id}` | DELETE, GET | — |
| `/{spaceId}/certificates` | GET, POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/all` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/generate` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}` | DELETE, GET, PUT | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/archive` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/archive/v1` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/export` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/replace` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/unarchive` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/unarchive/v1` | POST | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/certificates/{id}/usages` | GET | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| `/{spaceId}/channels` | GET, POST | octopusdeploy_channel |
| `/{spaceId}/channels/all` | GET | octopusdeploy_channel |
| `/{spaceId}/channels/rule-test` | GET, POST | octopusdeploy_channel |
| `/{spaceId}/channels/rule-test/v1` | GET, POST | octopusdeploy_channel |
| `/{spaceId}/channels/{id}` | DELETE, GET, PUT | octopusdeploy_channel |
| `/{spaceId}/channels/{id}/releases` | GET | octopusdeploy_channel |
| `/{spaceId}/dashboard` | GET | — |
| `/{spaceId}/dashboard/dynamic` | GET | — |
| `/{spaceId}/dashboardconfiguration` | GET, PUT | — |
| `/{spaceId}/deploymentTargetTags/{tag}` | GET | — |
| `/{spaceId}/deploymentprocesses` | GET | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/{spaceId}/deploymentprocesses/{deploymentProcessId}/template` | GET | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/{spaceId}/deploymentprocesses/{id}` | GET, PUT | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step (+2 more) |
| `/{spaceId}/deployments` | GET, POST | — |
| `/{spaceId}/deployments/create/tenanted/v1` | POST | — |
| `/{spaceId}/deployments/create/untenanted/v1` | POST | — |
| `/{spaceId}/deployments/override` | POST | — |
| `/{spaceId}/deployments/v1` | POST | — |
| `/{spaceId}/deployments/{id}` | DELETE, GET | — |
| `/{spaceId}/deploymentsettings/{id}` | GET | — |
| `/{spaceId}/deploymentsettings/{projectId}` | PUT | — |
| `/{spaceId}/deploymenttargettags` | GET, POST | — |
| `/{spaceId}/deploymenttargettags/{tag}` | DELETE | — |
| `/{spaceId}/environments` | GET, POST | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/all` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/all/v1` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/ephemeral/{id}/deprovision` | POST | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/sortorder` | PUT | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/summary` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/summary/v2` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/v1` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/v2` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/{environmentId}` | PUT | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/{environmentId}/metadata` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/{environmentId}/singlyScopedVariableDetails` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/{id}` | DELETE, GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/{id}/machines` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/environments/{id}/v2` | GET | octopusdeploy_environment, octopusdeploy_parent_environment |
| `/{spaceId}/events` | GET | — |
| `/{spaceId}/events/agents` | GET | — |
| `/{spaceId}/events/categories` | GET | — |
| `/{spaceId}/events/documenttypes` | GET | — |
| `/{spaceId}/events/groups` | GET | — |
| `/{spaceId}/events/{id}` | GET | — |
| `/{spaceId}/feeds` | GET, POST | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/all` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/stats` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/{feedId}/packages` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/{feedId}/packages/notes` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/{id}` | DELETE, GET, PUT | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/{id}/packages/search` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/feeds/{id}/packages/versions` | GET | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed (+10 more) |
| `/{spaceId}/git-credentials` | GET, POST | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/{spaceId}/git-credentials/v1` | GET, POST | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/{spaceId}/git-credentials/{id}` | DELETE, GET, PUT | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/{spaceId}/git-credentials/{id}/usage` | GET | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/{spaceId}/git-credentials/{id}/usage/v1` | GET | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/{spaceId}/git-credentials/{id}/v1` | DELETE, GET, PUT | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential |
| `/{spaceId}/insights/reports` | GET, POST | — |
| `/{spaceId}/insights/reports/v1` | POST | — |
| `/{spaceId}/insights/reports/{id}` | DELETE, GET, PUT | — |
| `/{spaceId}/insights/reports/{id}/v1` | DELETE, GET, PUT | — |
| `/{spaceId}/insights/reports/{reportId}/deployments` | GET | — |
| `/{spaceId}/insights/reports/{reportId}/deployments/csv` | GET | — |
| `/{spaceId}/insights/reports/{reportId}/logo` | GET, POST | — |
| `/{spaceId}/insights/reports/{reportId}/logo/icon` | POST | — |
| `/{spaceId}/insights/reports/{reportId}/logo/icon/v1` | POST | — |
| `/{spaceId}/insights/reports/{reportId}/metrics` | GET | — |
| `/{spaceId}/insights/reports/{reportId}/metrics/v1` | GET | — |
| `/{spaceId}/interruptions` | GET | — |
| `/{spaceId}/interruptions/{id}` | GET | — |
| `/{spaceId}/interruptions/{id}/responsible` | GET, PUT | — |
| `/{spaceId}/interruptions/{id}/submit` | POST | — |
| `/{spaceId}/libraryvariablesets` | GET, POST | octopusdeploy_library_variable_set |
| `/{spaceId}/libraryvariablesets/all` | GET | octopusdeploy_library_variable_set |
| `/{spaceId}/libraryvariablesets/all/v1` | GET, POST | octopusdeploy_library_variable_set |
| `/{spaceId}/libraryvariablesets/{id}` | DELETE, GET, PUT | octopusdeploy_library_variable_set |
| `/{spaceId}/libraryvariablesets/{id}/usages` | GET | octopusdeploy_library_variable_set |
| `/{spaceId}/lifecycles` | GET, POST | octopusdeploy_lifecycle |
| `/{spaceId}/lifecycles/all` | GET | octopusdeploy_lifecycle |
| `/{spaceId}/lifecycles/previews` | GET | octopusdeploy_lifecycle |
| `/{spaceId}/lifecycles/{id}` | DELETE, GET, PUT | octopusdeploy_lifecycle |
| `/{spaceId}/lifecycles/{id}/preview` | GET | octopusdeploy_lifecycle |
| `/{spaceId}/lifecycles/{id}/projects` | GET | octopusdeploy_lifecycle |
| `/{spaceId}/machinepolicies` | GET, POST | octopusdeploy_machine_policy |
| `/{spaceId}/machinepolicies/all` | GET | octopusdeploy_machine_policy |
| `/{spaceId}/machinepolicies/template` | GET | octopusdeploy_machine_policy |
| `/{spaceId}/machinepolicies/{id}` | DELETE, GET, PUT | octopusdeploy_machine_policy |
| `/{spaceId}/machinepolicies/{id}/machines` | GET | octopusdeploy_machine_policy |
| `/{spaceId}/machinepolicies/{id}/v1` | DELETE | octopusdeploy_machine_policy |
| `/{spaceId}/machinepolicies/{id}/workers` | GET | octopusdeploy_machine_policy |
| `/{spaceId}/machineroles/all` | GET | — |
| `/{spaceId}/machineroles/all/v1` | GET | — |
| `/{spaceId}/machines` | GET, POST | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/all/v1` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/discover` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/operatingsystem/names/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/operatingsystem/shells/all` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{id}` | DELETE, GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{id}/connection` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{id}/latestdeployments` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{id}/tasks` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{id}/tasks/v1` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{machineId}/singlyScopedVariableDetails` | GET | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/machines/{machineid}` | PUT | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target (+14 more) |
| `/{spaceId}/nuget/packages` | PUT | octopusdeploy_nuget_feed |
| `/{spaceId}/observability/agents` | POST | — |
| `/{spaceId}/observability/events/sessions` | POST | — |
| `/{spaceId}/observability/events/sessions/{sessionId}` | GET | — |
| `/{spaceId}/observability/kubernetes-monitors` | POST | — |
| `/{spaceId}/observability/kubernetes-monitors/{id}` | DELETE, GET | — |
| `/{spaceId}/observability/logs/sessions` | POST | — |
| `/{spaceId}/observability/logs/sessions/{sessionId}` | GET | — |
| `/{spaceId}/packages` | GET | — |
| `/{spaceId}/packages/bulk` | DELETE | — |
| `/{spaceId}/packages/bulk/v1` | DELETE | — |
| `/{spaceId}/packages/notes` | GET | — |
| `/{spaceId}/packages/raw` | POST | — |
| `/{spaceId}/packages/{id}` | DELETE, GET | — |
| `/{spaceId}/packages/{id}/raw` | GET | — |
| `/{spaceId}/packages/{id}/v1` | DELETE | — |
| `/{spaceId}/packages/{packageId}/{baseVersion}/delta` | POST | — |
| `/{spaceId}/packages/{packageId}/{version}/delta-signature` | GET | — |
| `/{spaceId}/parentEnvironments` | POST | octopusdeploy_parent_environment |
| `/{spaceId}/parentEnvironments/{environmentId}` | PUT | octopusdeploy_parent_environment |
| `/{spaceId}/parentEnvironments/{id}` | DELETE, GET | octopusdeploy_parent_environment |
| `/{spaceId}/processtemplates/{slug}/{versionMask}` | GET | octopusdeploy_process |
| `/{spaceId}/processtemplates/{slug}/{versionMask}/icon` | GET | octopusdeploy_process |
| `/{spaceId}/progression/runbooks/taskRuns` | GET | — |
| `/{spaceId}/progression/runbooks/{runbookId}` | GET | — |
| `/{spaceId}/progression/runbooks/{runbookId}/v1` | GET | — |
| `/{spaceId}/progression/{projectId}` | GET | — |
| `/{spaceId}/projectgroups` | GET, POST | octopusdeploy_project_group |
| `/{spaceId}/projectgroups/all` | GET | octopusdeploy_project_group |
| `/{spaceId}/projectgroups/{id}` | DELETE, GET, PUT | octopusdeploy_project_group |
| `/{spaceId}/projectgroups/{id}/projects` | GET | octopusdeploy_project_group |
| `/{spaceId}/projects` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/all` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{id}/logo` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/channels` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/channels/{channelId}/git-reference-rule-validation/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/channels/{channelId}/git-resource-rule-validation/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/channels/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/channels/{id}/releases` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/channels/{id}/v2` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/deploymentprocesses` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/deploymentprocesses/resolved` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/deploymentprocesses/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/deploymentprocesses/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/deploymentsettings` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral/{environmentId}/deprovision` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral/{environmentId}/deprovisioning/mark-successful` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral/{environmentId}/deprovisioning/retry` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral/{environmentId}/provisioning/mark-successful` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral/{environmentId}/provisioning/retry` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/ephemeral/{id}/status` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/{environmentId}/tenants/{tenantId}/livestatus` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/{environmentId}/tenants/{tenantId}/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/{environmentId}/tenants/{tenantId}/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}/manifest` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/{environmentId}/untenanted/livestatus` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/{environmentId}/untenanted/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/environments/{environmentId}/untenanted/machines/{sourceId}/resources/{desiredOrKubernetesMonitoredResourceId}/manifest` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/featuretoggles` | GET, POST, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/featuretoggles/generate-client-identifier` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/featuretoggles/rotate-client-identifier-signing-key` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/featuretoggles/{Id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/featuretoggles/{slug}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/branches` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/branches/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/branches/{branchName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/commits/{hash}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/connectivity-test` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/convert` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/migrate-runbooks` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/migrate-variables` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/refs` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/refs/{refName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/tags` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/tags/{tagName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/git/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/insights/deployments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/insights/deployments/csv` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/insights/metrics` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/insights/metrics/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/logo` | POST, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/logo/icon` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/metadata` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/progression` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/progression/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/releases` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/releases/{id}/variables` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/releases/{version}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookProcesses` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookProcesses/{id}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookRuns` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookRuns/{id}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookRuns/{id}/details/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{idOrName}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}` | PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookSnapshots/{id}/variables` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookruns/{id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbookruns/{runbookRunId}/retry/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/all` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/all/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}/environments/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{id}/runbookSnapshots` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{runbookId}/run` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{runbookId}/runbookRuns/previews` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/runbooksnapshots/{id}` | DELETE | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/summary` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/summary/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/triggers` | GET, POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/triggers/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/variables` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/deploymentprocesses` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/deploymentprocesses/resolved` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/deploymentprocesses/template` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/deploymentprocesses/validate` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/deploymentsettings` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/v2` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{id}/environments` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{id}/environments/v2` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{runbookId}/run/v1` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/runbooks/{runbookId}/runbookRuns/previews` | POST | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/summary` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/summary/v1` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitRef}/variables` | GET, PUT | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{gitref}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projects/{projectId}/{unusedGitRef}` | GET | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy |
| `/{spaceId}/projecttriggers` | GET, POST | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger (+2 more) |
| `/{spaceId}/projecttriggers/{id}` | DELETE, GET, PUT | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger (+2 more) |
| `/{spaceId}/proxies` | GET, POST | octopusdeploy_machine_proxy |
| `/{spaceId}/proxies/all` | GET | octopusdeploy_machine_proxy |
| `/{spaceId}/proxies/{id}` | DELETE, GET, PUT | octopusdeploy_machine_proxy |
| `/{spaceId}/releases` | GET, POST | — |
| `/{spaceId}/releases/create/v1` | POST | — |
| `/{spaceId}/releases/{id}` | DELETE, GET, PUT | — |
| `/{spaceId}/releases/{id}/deployments/template` | GET | — |
| `/{spaceId}/releases/{releaseId}/defects` | GET, POST | — |
| `/{spaceId}/releases/{releaseId}/defects/resolve` | POST | — |
| `/{spaceId}/releases/{releaseId}/deployments` | GET | — |
| `/{spaceId}/releases/{releaseId}/deployments/preview/{environmentId}` | GET | — |
| `/{spaceId}/releases/{releaseId}/deployments/preview/{environmentId}/{tenantId}` | GET | — |
| `/{spaceId}/releases/{releaseId}/deployments/previews` | POST | — |
| `/{spaceId}/releases/{releaseId}/missingPackages` | GET | — |
| `/{spaceId}/releases/{releaseId}/progression` | GET | — |
| `/{spaceId}/releases/{releaseId}/snapshot-variables` | POST | — |
| `/{spaceId}/reporting/deployments-counted-by-week` | GET | — |
| `/{spaceId}/reporting/deployments/xml` | GET | — |
| `/{spaceId}/retentionpolicies` | GET | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy |
| `/{spaceId}/retentionpolicies/{id}` | PUT | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy |
| `/{spaceId}/runbook-runs/create/v1` | POST | octopusdeploy_runbook |
| `/{spaceId}/runbookProcesses` | GET | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/{spaceId}/runbookProcesses/{id}` | GET, PUT | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/{spaceId}/runbookProcesses/{id}/runbookSnapshotTemplate` | GET | octopusdeploy_runbook, octopusdeploy_runbook_process |
| `/{spaceId}/runbookRuns` | GET, POST | octopusdeploy_runbook |
| `/{spaceId}/runbookRuns/{id}` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots` | GET, POST | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots/{id}` | GET, PUT | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots/{id}/runbookRuns` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots/{id}/runbookRuns/template` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbookSnapshots/{id}/snapshot-variables` | POST | octopusdeploy_runbook |
| `/{spaceId}/runbookruns/{id}` | DELETE | octopusdeploy_runbook |
| `/{spaceId}/runbooks` | GET, POST | octopusdeploy_runbook |
| `/{spaceId}/runbooks/all` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{id}` | DELETE, GET, PUT | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{id}/environments` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{id}/runbookRunTemplate` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{id}/runbookRuns/preview/{environment}` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{id}/runbookSnapshots` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{runbookId}/run` | POST | octopusdeploy_runbook |
| `/{spaceId}/runbooks/{runbookId}/runbookSnapshotTemplate` | GET | octopusdeploy_runbook |
| `/{spaceId}/runbooksnapshots/{id}` | DELETE | octopusdeploy_runbook |
| `/{spaceId}/scopeduserroles` | GET, POST | octopusdeploy_scoped_user_role |
| `/{spaceId}/scopeduserroles/{id}` | DELETE, GET, PUT | octopusdeploy_scoped_user_role |
| `/{spaceId}/spaces/{id}/search` | GET | octopusdeploy_space |
| `/{spaceId}/subscriptions` | GET, POST | — |
| `/{spaceId}/subscriptions/all` | GET | — |
| `/{spaceId}/subscriptions/{id}` | DELETE, GET, PUT | — |
| `/{spaceId}/tagsets` | GET, POST | octopusdeploy_tag, octopusdeploy_tag_set |
| `/{spaceId}/tagsets/all` | GET | octopusdeploy_tag, octopusdeploy_tag_set |
| `/{spaceId}/tagsets/sortorder` | PUT | octopusdeploy_tag, octopusdeploy_tag_set |
| `/{spaceId}/tagsets/{id}` | DELETE, GET, PUT | octopusdeploy_tag, octopusdeploy_tag_set |
| `/{spaceId}/tasks` | GET, POST | — |
| `/{spaceId}/tasks/rerun/{id}` | POST | — |
| `/{spaceId}/tasks/tasktypes` | GET | — |
| `/{spaceId}/tasks/{id}` | GET | — |
| `/{spaceId}/tasks/{id}/cancel` | POST | — |
| `/{spaceId}/tasks/{id}/details` | GET | — |
| `/{spaceId}/tasks/{id}/prioritize` | POST | — |
| `/{spaceId}/tasks/{id}/queued-behind` | GET | — |
| `/{spaceId}/tasks/{id}/raw` | GET | — |
| `/{spaceId}/tasks/{id}/state` | POST | — |
| `/{spaceId}/tasks/{id}/status/messages` | GET | — |
| `/{spaceId}/teammembership` | GET | octopusdeploy_team |
| `/{spaceId}/teammembership/previewteam` | POST | octopusdeploy_team |
| `/{spaceId}/teams` | GET, POST | octopusdeploy_team |
| `/{spaceId}/teams/all` | GET | octopusdeploy_team |
| `/{spaceId}/teams/{id}` | DELETE, GET, PUT | octopusdeploy_team |
| `/{spaceId}/teams/{id}/scopeduserroles` | GET | octopusdeploy_team |
| `/{spaceId}/tenants` | GET, POST | octopusdeploy_tenant |
| `/{spaceId}/tenants/all` | GET | octopusdeploy_tenant |
| `/{spaceId}/tenants/status` | GET | octopusdeploy_tenant |
| `/{spaceId}/tenants/tag-test` | GET | octopusdeploy_tenant |
| `/{spaceId}/tenants/variables-missing` | GET | octopusdeploy_tenant |
| `/{spaceId}/tenants/{id}` | DELETE, GET, PUT | octopusdeploy_tenant |
| `/{spaceId}/tenants/{id}/logo` | GET, POST, PUT | octopusdeploy_tenant |
| `/{spaceId}/tenants/{id}/variables` | GET, POST, PUT | octopusdeploy_tenant |
| `/{spaceId}/tenants/{tenantId}/commonvariables` | GET, POST, PUT | octopusdeploy_tenant |
| `/{spaceId}/tenants/{tenantId}/logo/icon` | POST | octopusdeploy_tenant |
| `/{spaceId}/tenants/{tenantId}/projectvariables` | GET, POST, PUT | octopusdeploy_tenant |
| `/{spaceId}/tenantvariables/all` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable |
| `/{spaceId}/users/invitations` | POST | octopusdeploy_user |
| `/{spaceId}/users/invitations/{id}` | GET | octopusdeploy_user |
| `/{spaceId}/users/{id}/permissions` | GET | octopusdeploy_user |
| `/{spaceId}/users/{id}/permissions/configuration` | GET | octopusdeploy_user |
| `/{spaceId}/users/{id}/permissions/export` | GET | octopusdeploy_user |
| `/{spaceId}/users/{userId}/teams` | GET | octopusdeploy_user |
| `/{spaceId}/variables/all` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/{spaceId}/variables/names` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/{spaceId}/variables/preview` | GET | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/{spaceId}/variables/{id}` | GET, PUT | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable |
| `/{spaceId}/workerpools` | GET, POST | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/dynamicworkertypes` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/sortorder` | PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/summary` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/supportedtypes` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/{id}` | DELETE, GET, PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workerpools/{id}/workers` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool |
| `/{spaceId}/workers` | GET, POST | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workers/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workers/discover` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workers/operatingsystem/names/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workers/operatingsystem/shells/all` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workers/{id}` | DELETE, GET, PUT | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workers/{id}/connection` | GET | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor (+3 more) |
| `/{spaceId}/workertaskleases` | GET | — |

## 2. Path resources (unique API areas) and TF coverage

| API path resource | TF resources | Status |
|-------------------|--------------|--------|
| .well-known | — | ❌ Missing |
| accounts | octopusdeploy_amazon_web_services_account, octopusdeploy_aws_openid_connect_account, octopusdeploy_azure_subscription_account, octopusdeploy_gcp_account, octopusdeploy_platform_hub_aws_account ... | ✅ |
| actionTemplates | octopusdeploy_community_step_template, octopusdeploy_step_template | ✅ |
| artifacts | — | ❌ Missing |
| audit-stream | — | ❌ Missing |
| authentication | — | ❌ Missing |
| azuredevopsissuetracker | — | ❌ Missing |
| build-information | — | ❌ Missing |
| capabilities | — | ❌ Missing |
| certificates | octopusdeploy_certificate, octopusdeploy_tentacle_certificate | ✅ |
| channels | octopusdeploy_channel | ✅ |
| cloudtemplate | — | ❌ Missing |
| communityactiontemplates | octopusdeploy_community_step_template | ✅ |
| configuration | — | ❌ Missing |
| dashboard | — | ❌ Missing |
| dashboardconfiguration | — | ❌ Missing |
| deploymentTargetTags | — | ❌ Missing |
| deploymentfreezes | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project, octopusdeploy_deployment_freeze_tenant, octopusdeploy_project_deployment_freeze | ✅ |
| deploymentprocesses | octopusdeploy_deployment_process, octopusdeploy_process, octopusdeploy_process_child_step, octopusdeploy_process_step, octopusdeploy_runbook_process | ✅ |
| deployments | — | ❌ Missing |
| deploymentsettings | — | ❌ Missing |
| deprecations | — | ❌ Missing |
| dynamic-extensions | — | ❌ Missing |
| environments | octopusdeploy_environment, octopusdeploy_parent_environment | ✅ |
| events | — | ❌ Missing |
| externalgroups | — | ❌ Missing |
| externalsecuritygroupproviders | — | ❌ Missing |
| externalusers | octopusdeploy_user | ✅ |
| featuresconfiguration | — | ❌ Missing |
| feeds | octopusdeploy_artifactory_generic_feed, octopusdeploy_aws_elastic_container_registry_feed, octopusdeploy_azure_container_registry_feed, octopusdeploy_docker_container_registry_feed, octopusdeploy_git_hub_repository_feed ... | ✅ |
| git-credentials | octopusdeploy_git_credential, octopusdeploy_platform_hub_git_credential | ✅ |
| githubissuetracker | — | ❌ Missing |
| icons | — | ❌ Missing |
| insights | — | ❌ Missing |
| integrated-challenge | — | ❌ Missing |
| interruptions | — | ❌ Missing |
| jiraintegration | — | ❌ Missing |
| jiraservicemanagement-integration | — | ❌ Missing |
| letsencryptconfiguration | — | ❌ Missing |
| libraryvariablesets | octopusdeploy_library_variable_set | ✅ |
| licenses | — | ❌ Missing |
| lifecycles | octopusdeploy_lifecycle | ✅ |
| logo | — | ❌ Missing |
| machinepolicies | octopusdeploy_machine_policy | ✅ |
| machineroles | — | ❌ Missing |
| machines | octopusdeploy_azure_cloud_service_deployment_target, octopusdeploy_azure_service_fabric_cluster_deployment_target, octopusdeploy_azure_web_app_deployment_target, octopusdeploy_cloud_region_deployment_target, octopusdeploy_dynamic_worker_pool ... | ✅ |
| maintenanceconfiguration | — | ❌ Missing |
| migrations | — | ❌ Missing |
| nuget | octopusdeploy_nuget_feed | ✅ |
| observability | — | ❌ Missing |
| octopusservernodes | — | ❌ Missing |
| packages | — | ❌ Missing |
| parentEnvironments | octopusdeploy_parent_environment | ✅ |
| performanceconfiguration | — | ❌ Missing |
| permissions | — | ❌ Missing |
| platformhub | octopusdeploy_platform_hub_aws_account, octopusdeploy_platform_hub_aws_open_i_d_connect_account, octopusdeploy_platform_hub_azure_oidc_account, octopusdeploy_platform_hub_azure_service_principal_account, octopusdeploy_platform_hub_gcp_account ... | ✅ |
| processtemplates | octopusdeploy_process | ✅ |
| progression | — | ❌ Missing |
| projectgroups | octopusdeploy_project_group | ✅ |
| projects | octopusdeploy_project, octopusdeploy_project_auto_create_release, octopusdeploy_project_versioning_strategy | ✅ |
| projecttriggers | octopusdeploy_built_in_trigger, octopusdeploy_external_feed_create_release_trigger, octopusdeploy_git_trigger, octopusdeploy_project_deployment_target_trigger, octopusdeploy_project_scheduled_trigger | ✅ |
| proxies | octopusdeploy_machine_proxy | ✅ |
| releases | — | ❌ Missing |
| reporting | — | ❌ Missing |
| retentionpolicies | octopusdeploy_space_default_lifecycle_release_retention_policy, octopusdeploy_space_default_lifecycle_tentacle_retention_policy, octopusdeploy_space_default_runbook_retention_policy | ✅ |
| runbook-runs | octopusdeploy_runbook | ✅ |
| runbookProcesses | octopusdeploy_runbook, octopusdeploy_runbook_process | ✅ |
| runbookSnapshots | octopusdeploy_runbook | ✅ |
| runbooks | octopusdeploy_runbook | ✅ |
| scheduler | — | ❌ Missing |
| scopeduserroles | octopusdeploy_scoped_user_role | ✅ |
| search | — | ❌ Missing |
| serverconfiguration | — | ❌ Missing |
| serverstatus | — | ❌ Missing |
| serviceaccounts | — | ❌ Missing |
| servicenow-integration | — | ❌ Missing |
| signingkeyconfiguration | — | ❌ Missing |
| smtpconfiguration | — | ❌ Missing |
| spaces | octopusdeploy_space | ✅ |
| subscriptions | — | ❌ Missing |
| tagsets | octopusdeploy_tag, octopusdeploy_tag_set | ✅ |
| tasks | — | ❌ Missing |
| teammembership | octopusdeploy_team | ✅ |
| teams | octopusdeploy_team | ✅ |
| telemetry | — | ❌ Missing |
| telemetryconfiguration | — | ❌ Missing |
| tenants | octopusdeploy_tenant | ✅ |
| tenantvariables | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable | ✅ |
| token | octopusdeploy_token_account | ✅ |
| upgradeconfiguration | — | ❌ Missing |
| userroles | octopusdeploy_scoped_user_role, octopusdeploy_user_role | ✅ |
| users | octopusdeploy_user | ✅ |
| v1 | — | ❌ Missing |
| variables | octopusdeploy_tenant_common_variable, octopusdeploy_tenant_project_variable, octopusdeploy_variable | ✅ |
| workerpools | octopusdeploy_dynamic_worker_pool, octopusdeploy_static_worker_pool | ✅ |
| workers | octopusdeploy_dynamic_worker_pool, octopusdeploy_kubernetes_agent_worker, octopusdeploy_kubernetes_monitor, octopusdeploy_listening_tentacle_worker, octopusdeploy_ssh_connection_worker ... | ✅ |
| workertaskleases | — | ❌ Missing |

## 3. API tags (operation groups)

| Tag |
|-----|
| AccessTokens |
| Accounts |
| ActionTemplates |
| ApiKeys |
| Artifacts |
| AuditStream |
| Authentication |
| AzureDevOps |
| Branches |
| BuildInformation |
| Capabilities |
| Certificates |
| Channels |
| CloudTemplate |
| CommunityActionTemplates |
| CompliancePolicies |
| Configuration |
| Dashboard |
| DashboardConfiguration |
| DeploymentFreeze |
| DeploymentProcesses |
| DeploymentSettings |
| DeploymentTargetTags |
| DeploymentTargets |
| Deployments |
| Deprecations |
| DirectoryServices |
| DynamicExtensions |
| Environments |
| EphemeralEnvironments |
| EventRetention |
| Events |
| ExternalSecurityGroupProviders |
| FeaturesConfiguration |
| Feeds |
| GitHub |
| Home |
| Icons |
| Insights |
| IntegratedAuthentication |
| Interruptions |
| Invitations |
| JiraIntegration |
| JsonWebKeys |
| LDAP |
| LetsEncrypt |
| LibraryVariableSets |
| Licenses |
| Lifecycles |
| MachinePolicies |
| MachineRoles |
| Machines |
| MaintenanceConfiguration |
| Migrations |
| Nuget |
| Observability |
| OctopusServerNodes |
| OpenIDConnect |
| OpenIdConnect |
| Packages |
| ParentEnvironments |
| Performance |
| Permissions |
| PlatformHub |
| ProcessTemplates |
| Progression |
| ProjectGroups |
| ProjectTemplates |
| ProjectTriggers |
| Projects |
| Proxies |
| Releases |
| Reporting |
| Retention |
| RunbookProcesses |
| RunbookRuns |
| RunbookSnapshots |
| Runbooks |
| ScheduledJobs |
| ScopedUserRoles |
| Server |
| ServerStatus |
| ServiceAccountOidcIdentities |
| Signing |
| Smtp |
| Spaces |
| Subscriptions |
| TagSets |
| Tasks |
| TeamMemberships |
| Teams |
| Telemetry |
| Tenants |
| Toggles |
| TokenExchange |
| Upgrade |
| UserPermissions |
| UserRoles |
| Users |
| Variables |
| VersionControl |
| Web |
| WorkerPools |
| WorkerTaskLeases |
| Workers |

## 4. API *Resource definitions – property coverage

For each API resource type, we list properties that appear in the API but not in the mapped Terraform schema (missing or different name).

### AccountResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AccountType (`account_type`) | — |
| Description (`description`) | ✅ |
| EnvironmentIds (`environment_ids`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |
| TenantIds (`tenant_ids`) | ❌ |
| TenantTags (`tenant_tags`) | ✅ |
| TenantedDeploymentParticipation (`tenanted_deployment_participation`) | ✅ |

### ActionTemplateCategoryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DisplayOrder (`display_order`) | ❌ |
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### ActionTemplateParameterResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DefaultValue (`default_value`) | ❌ |
| DisplaySettings (`display_settings`) | ❌ |
| HelpText (`help_text`) | ❌ |
| Id (`id`) | ❌ |
| Label (`label`) | ❌ |
| Name (`name`) | ❌ |

### ActionTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionType (`action_type`) | ✅ |
| CommunityActionTemplateId (`community_action_template_id`) | ✅ |
| Description (`description`) | ✅ |
| GitDependencies (`git_dependencies`) | ✅ |
| Id (`id`) | ✅ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| Packages (`packages`) | ✅ |
| Parameters (`parameters`) | ✅ |
| Properties (`properties`) | ✅ |
| SpaceId (`space_id`) | ✅ |
| Version (`version`) | ✅ |

### ActionUpdateResultResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| ManualMergeRequiredReasonsByPropertyName (`manual_merge_required_reasons_by_property_name`) | ❌ |
| NamesOfNewParametersMissingDefaultValue (`names_of_new_parameters_missing_default_value`) | ❌ |
| Outcome (`outcome`) | ❌ |
| RemovedPackageUsages (`removed_package_usages`) | ❌ |

### ActionsUpdateProcessResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionIds (`action_ids`) | ❌ |
| GitRef (`git_ref`) | ❌ |
| ProcessId (`process_id`) | ❌ |
| ProcessType (`process_type`) | ❌ |
| ProjectId (`project_id`) | ❌ |

### ApiKeyCreatedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ApiKey (`api_key`) | ❌ |
| Created (`created`) | ❌ |
| Expires (`expires`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Purpose (`purpose`) | ❌ |
| UserId (`user_id`) | ❌ |

### ApiKeyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ApiKey (`api_key`) | ❌ |
| Created (`created`) | ❌ |
| Expires (`expires`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Purpose (`purpose`) | ❌ |
| UserId (`user_id`) | ❌ |

### ArchivedEventFileResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CreatedDate (`created_date`) | ❌ |
| FileBytes (`file_bytes`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ModifiedDate (`modified_date`) | ❌ |
| Name (`name`) | ❌ |

### ArtifactResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Created (`created`) | ❌ |
| Filename (`filename`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| LogCorrelationId (`log_correlation_id`) | ❌ |
| ServerTaskId (`server_task_id`) | ❌ |
| Source (`source`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### AuditStreamConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Active (`active`) | ❌ |
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| StreamConfigurationResource (`stream_configuration_resource`) | ❌ |

### AuthenticationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AnyAuthenticationProvidersSupportPasswordManagement (`any_authentication_providers_support_password_management`) | ❌ |
| ApiKeyDefaultExpiryDays (`api_key_default_expiry_days`) | ❌ |
| ApiKeyMaxExpiryDays (`api_key_max_expiry_days`) | ❌ |
| AuthenticationProviders (`authentication_providers`) | ❌ |
| AutoLoginEnabled (`auto_login_enabled`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| RememberMeEnabled (`remember_me_enabled`) | ❌ |
| UserApiKeysEnabled (`user_api_keys_enabled`) | ❌ |

### AutoDeployReleaseOverrideResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EnvironmentId (`environment_id`) | ❌ |
| ReleaseId (`release_id`) | ❌ |
| TenantId (`tenant_id`) | ❌ |

### AutomaticDeprovisioningRuleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ExpiryDays (`expiry_days`) | ❌ |
| ExpiryHours (`expiry_hours`) | ❌ |

### BaseEnvironmentV2Resource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| EnvironmentTags (`environment_tags`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Type (`type`) | — |

### BuiltInFeedStatsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| TotalPackages (`total_packages`) | ❌ |

### CertificateConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| SignatureAlgorithm (`signature_algorithm`) | ❌ |
| Thumbprint (`thumbprint`) | ❌ |

### CertificateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Archived (`archived`) | ✅ |
| CertificateChain (`certificate_chain`) | ❌ |
| CertificateData (`certificate_data`) | ✅ |
| CertificateDataFormat (`certificate_data_format`) | ✅ |
| EnvironmentIds (`environment_ids`) | ❌ |
| HasPrivateKey (`has_private_key`) | ✅ |
| Id (`id`) | ❌ |
| IsExpired (`is_expired`) | ✅ |
| IssuerCommonName (`issuer_common_name`) | ✅ |
| IssuerDistinguishedName (`issuer_distinguished_name`) | ✅ |
| IssuerOrganization (`issuer_organization`) | ✅ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| NotAfter (`not_after`) | ✅ |
| NotBefore (`not_before`) | ✅ |
| Notes (`notes`) | ✅ |
| Password (`password`) | ✅ |
| ReplacedBy (`replaced_by`) | ✅ |
| SelfSigned (`self_signed`) | ✅ |
| SerialNumber (`serial_number`) | ✅ |
| SignatureAlgorithmName (`signature_algorithm_name`) | ✅ |
| SpaceId (`space_id`) | ✅ |
| SubjectAlternativeNames (`subject_alternative_names`) | ✅ |
| SubjectCommonName (`subject_common_name`) | ✅ |
| SubjectDistinguishedName (`subject_distinguished_name`) | ✅ |
| SubjectOrganization (`subject_organization`) | ✅ |
| TenantIds (`tenant_ids`) | ❌ |
| TenantTags (`tenant_tags`) | ✅ |
| TenantedDeploymentParticipation (`tenanted_deployment_participation`) | ✅ |
| Thumbprint (`thumbprint`) | ✅ |
| Version (`version`) | ✅ |

### CertificateUsageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentTargetUsages (`deployment_target_usages`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LibraryVariableSetUsages (`library_variable_set_usages`) | ❌ |
| Links (`links`) | ❌ |
| ProjectUsages (`project_usages`) | ❌ |
| TenantUsages (`tenant_usages`) | ❌ |

### ChangeDetailsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Differences (`differences`) | ❌ |
| DocumentContext (`document_context`) | ❌ |

### ChannelCustomFieldDefinitionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| FieldName (`field_name`) | ❌ |

### ChannelGitResourceRuleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| GitDependencyActions (`git_dependency_actions`) | ❌ |
| Id (`id`) | ❌ |
| Rules (`rules`) | ❌ |

### ChannelResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AutomaticEphemeralEnvironmentDeployments (`automatic_ephemeral_environment_deployments`) | ❌ |
| CustomFieldDefinitions (`custom_field_definitions`) | ❌ |
| Description (`description`) | ✅ |
| EphemeralEnvironmentNameTemplate (`ephemeral_environment_name_template`) | ✅ |
| GitReferenceRules (`git_reference_rules`) | ❌ |
| GitResourceRules (`git_resource_rules`) | ❌ |
| Id (`id`) | ❌ |
| IsDefault (`is_default`) | ✅ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LifecycleId (`lifecycle_id`) | ✅ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| ParentEnvironmentId (`parent_environment_id`) | ✅ |
| ProjectId (`project_id`) | ✅ |
| Rules (`rules`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |
| TenantTags (`tenant_tags`) | ✅ |
| Type (`type`) | ✅ |

### ChannelVersionRuleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionPackages (`action_packages`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Tag (`tag`) | ❌ |
| VersionRange (`version_range`) | ❌ |

### CommunityActionTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Author (`author`) | ❌ |
| Description (`description`) | ❌ |
| HistoryUrl (`history_url`) | ❌ |
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| Packages (`packages`) | ❌ |
| Parameters (`parameters`) | — |
| Properties (`properties`) | — |
| Type (`type`) | ❌ |
| Version (`version`) | ❌ |
| Website (`website`) | ❌ |

### CommunityActionTemplateSnapshotResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| Parameters (`parameters`) | ❌ |
| Version (`version`) | ❌ |

### CompliancePolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ConditionsRego (`conditions_rego`) | ❌ |
| Description (`description`) | ❌ |
| GitRef (`git_ref`) | ❌ |
| Name (`name`) | ❌ |
| ScopeRego (`scope_rego`) | ❌ |
| Slug (`slug`) | ❌ |
| ViolationAction (`violation_action`) | ❌ |
| ViolationReason (`violation_reason`) | ❌ |

### CompliancePolicyVersionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| GitCommit (`git_commit`) | ❌ |
| GitRef (`git_ref`) | ❌ |
| Id (`id`) | ❌ |
| IsActive (`is_active`) | ❌ |
| Name (`name`) | ❌ |
| PublishedDate (`published_date`) | ❌ |
| RegoConditions (`rego_conditions`) | ❌ |
| RegoScope (`rego_scope`) | ❌ |
| Slug (`slug`) | ❌ |
| Version (`version`) | ❌ |
| ViolationAction (`violation_action`) | ❌ |
| ViolationReason (`violation_reason`) | ❌ |

### ContainerLogLineResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Message (`message`) | ❌ |
| Timestamp (`timestamp`) | ❌ |

### DashboardConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IncludedEnvironmentIds (`included_environment_ids`) | ❌ |
| IncludedEnvironmentTags (`included_environment_tags`) | ❌ |
| IncludedProjectGroupIds (`included_project_group_ids`) | ❌ |
| IncludedProjectIds (`included_project_ids`) | ❌ |
| IncludedProjectTags (`included_project_tags`) | ❌ |
| IncludedTenantIds (`included_tenant_ids`) | ❌ |
| IncludedTenantTags (`included_tenant_tags`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectLimit (`project_limit`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### DashboardEnvironmentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### DashboardItemResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ChannelId (`channel_id`) | ❌ |
| CompletedTime (`completed_time`) | ❌ |
| Created (`created`) | ❌ |
| DeploymentId (`deployment_id`) | ❌ |
| Duration (`duration`) | ❌ |
| EnvironmentId (`environment_id`) | ❌ |
| ErrorMessage (`error_message`) | ❌ |
| HasPendingInterruptions (`has_pending_interruptions`) | ❌ |
| HasWarningsOrErrors (`has_warnings_or_errors`) | ❌ |
| Id (`id`) | ❌ |
| IsCompleted (`is_completed`) | ❌ |
| IsCurrent (`is_current`) | ❌ |
| IsPrevious (`is_previous`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| QueueTime (`queue_time`) | ❌ |
| ReleaseId (`release_id`) | ❌ |
| ReleaseVersion (`release_version`) | ❌ |
| StartTime (`start_time`) | ❌ |
| State (`state`) | ❌ |
| TaskId (`task_id`) | ❌ |
| TenantId (`tenant_id`) | ❌ |

### DashboardProjectGroupResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EnvironmentIds (`environment_ids`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### DashboardProjectResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanPerformUntenantedDeployment (`can_perform_untenanted_deployment`) | ❌ |
| EnvironmentIds (`environment_ids`) | ❌ |
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| ProjectGroupId (`project_group_id`) | ❌ |
| Slug (`slug`) | ❌ |
| TenantedDeploymentMode (`tenanted_deployment_mode`) | ❌ |

### DashboardResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Environments (`environments`) | ❌ |
| Id (`id`) | ❌ |
| IsFiltered (`is_filtered`) | ❌ |
| Items (`items`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectGroups (`project_groups`) | ❌ |
| ProjectLimit (`project_limit`) | ❌ |
| Projects (`projects`) | ❌ |
| Tenants (`tenants`) | ❌ |

### DashboardTenantResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| ProjectEnvironments (`project_environments`) | ❌ |
| TenantTags (`tenant_tags`) | ❌ |

### DefectResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Status (`status`) | ❌ |

### DeploymentActionContainerResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Dockerfile (`dockerfile`) | ❌ |
| FeedId (`feed_id`) | ❌ |
| GitUrl (`git_url`) | ❌ |
| Image (`image`) | ❌ |

### DeploymentActionGitDependencyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentActionSlug (`deployment_action_slug`) | ❌ |
| GitDependencyName (`git_dependency_name`) | ❌ |

### DeploymentActionPackageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentAction (`deployment_action`) | ❌ |
| PackageReference (`package_reference`) | ❌ |

### DeploymentActionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionType (`action_type`) | ❌ |
| AvailableStepPackageVersions (`available_step_package_versions`) | ❌ |
| CanBeUsedForProjectVersioning (`can_be_used_for_project_versioning`) | ❌ |
| Channels (`channels`) | — |
| ChannelsVariable (`channels_variable`) | ❌ |
| CommunityActionTemplateSnapshot (`community_action_template_snapshot`) | ❌ |
| Condition (`condition`) | ❌ |
| Container (`container`) | ❌ |
| Environments (`environments`) | — |
| EnvironmentsVariable (`environments_variable`) | ❌ |
| ExcludedEnvironments (`excluded_environments`) | — |
| ExcludedEnvironmentsVariable (`excluded_environments_variable`) | ❌ |
| GitDependencies (`git_dependencies`) | ❌ |
| Id (`id`) | ❌ |
| Inputs (`inputs`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| IsRequired (`is_required`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastSavedStepPackageVersion (`last_saved_step_package_version`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| Notes (`notes`) | ❌ |
| Packages (`packages`) | ❌ |
| Properties (`properties`) | — |
| Slug (`slug`) | ❌ |
| StepPackageVersion (`step_package_version`) | ❌ |
| TenantTags (`tenant_tags`) | — |
| TenantTagsVariable (`tenant_tags_variable`) | ❌ |
| WorkerPoolId (`worker_pool_id`) | ❌ |
| WorkerPoolVariable (`worker_pool_variable`) | ❌ |

### DeploymentFreezeResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| End (`end`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| ProjectEnvironmentScope (`project_environment_scope`) | ❌ |
| RecurringSchedule (`recurring_schedule`) | ❌ |
| Start (`start`) | ❌ |
| TenantProjectEnvironmentScope (`tenant_project_environment_scope`) | ❌ |

### DeploymentPreviewResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Changes (`changes`) | ❌ |
| ChangesMarkdown (`changes_markdown`) | ❌ |
| Form (`form`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| StepsToExecute (`steps_to_execute`) | ❌ |
| UseGuidedFailureModeByDefault (`use_guided_failure_mode_by_default`) | ❌ |

### DeploymentProcessResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastSnapshotId (`last_snapshot_id`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Steps (`steps`) | ❌ |
| Version (`version`) | ❌ |

### DeploymentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ChangeRequestSettings (`change_request_settings`) | ❌ |
| Changes (`changes`) | ❌ |
| ChangesMarkdown (`changes_markdown`) | ❌ |
| ChannelId (`channel_id`) | ❌ |
| Comments (`comments`) | ❌ |
| Created (`created`) | ❌ |
| DebugMode (`debug_mode`) | ❌ |
| DeployedBy (`deployed_by`) | ❌ |
| DeployedById (`deployed_by_id`) | ❌ |
| DeployedToMachineIds (`deployed_to_machine_ids`) | ❌ |
| DeploymentProcessId (`deployment_process_id`) | ❌ |
| EnvironmentId (`environment_id`) | ❌ |
| ExcludedMachineIds (`excluded_machine_ids`) | ❌ |
| ExecutionPlanLogContext (`execution_plan_log_context`) | ❌ |
| FailTargetDiscovery (`fail_target_discovery`) | ❌ |
| FailureEncountered (`failure_encountered`) | ❌ |
| ForcePackageDownload (`force_package_download`) | ❌ |
| ForcePackageRedeployment (`force_package_redeployment`) | ❌ |
| FormValues (`form_values`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ManifestVariableSetId (`manifest_variable_set_id`) | ❌ |
| Name (`name`) | ❌ |
| Priority (`priority`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| QueueTime (`queue_time`) | ❌ |
| QueueTimeExpiry (`queue_time_expiry`) | ❌ |
| ReleaseId (`release_id`) | ❌ |
| SkipActions (`skip_actions`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| SpecificMachineIds (`specific_machine_ids`) | ❌ |
| TaskId (`task_id`) | ❌ |
| TenantId (`tenant_id`) | ❌ |
| TentacleRetentionPeriod (`tentacle_retention_period`) | ❌ |
| UseGuidedFailure (`use_guided_failure`) | ❌ |

### DeploymentSettingsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ConnectivityPolicy (`connectivity_policy`) | ❌ |
| DefaultGuidedFailureMode (`default_guided_failure_mode`) | ❌ |
| DefaultToSkipIfAlreadyInstalled (`default_to_skip_if_already_installed`) | ❌ |
| DeploymentChangesTemplate (`deployment_changes_template`) | ❌ |
| FailTargetDiscovery (`fail_target_discovery`) | ❌ |
| ForcePackageDownload (`force_package_download`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| ReleaseNotesTemplate (`release_notes_template`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| VersioningStrategy (`versioning_strategy`) | ❌ |

### DeploymentStepResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Actions (`actions`) | ❌ |
| Condition (`condition`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| PackageRequirement (`package_requirement`) | ❌ |
| Properties (`properties`) | ❌ |
| Slug (`slug`) | ❌ |
| StartTrigger (`start_trigger`) | ❌ |
| Type (`type`) | — |

### DeploymentTargetTagResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| SpaceId (`space_id`) | ❌ |
| Tag (`tag`) | ❌ |

### DeploymentTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentNotes (`deployment_notes`) | ❌ |
| Id (`id`) | ❌ |
| IsDeploymentProcessModified (`is_deployment_process_modified`) | ❌ |
| IsLibraryVariableSetModified (`is_library_variable_set_modified`) | ❌ |
| IsVariableSetModified (`is_variable_set_modified`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| PromoteTo (`promote_to`) | ❌ |
| TenantPromotions (`tenant_promotions`) | ❌ |

### DocumentTypeResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |

### DynamicExtensionsFeaturesMetadataResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Features (`features`) | ❌ |

### DynamicExtensionsFeaturesValuesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Values (`values`) | ❌ |

### EndpointResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CommunicationStyle (`communication_style`) | — |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### EnvironmentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AllowDynamicInfrastructure (`allow_dynamic_infrastructure`) | ✅ |
| Description (`description`) | ✅ |
| EnvironmentTags (`environment_tags`) | ✅ |
| ExtensionSettings (`extension_settings`) | — |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| Slug (`slug`) | ✅ |
| SortOrder (`sort_order`) | ✅ |
| SpaceId (`space_id`) | ✅ |
| UseGuidedFailure (`use_guided_failure`) | ✅ |

### EnvironmentSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentTargetSummaries (`deployment_target_summaries`) | ❌ |
| Environment (`environment`) | ❌ |
| MachineEndpointSummaries (`machine_endpoint_summaries`) | ❌ |
| MachineHealthStatusSummaries (`machine_health_status_summaries`) | ❌ |
| MachineIdsForCalamariUpgrade (`machine_ids_for_calamari_upgrade`) | ❌ |
| MachineIdsForTentacleUpgrade (`machine_ids_for_tentacle_upgrade`) | ❌ |
| MachineTenantSummaries (`machine_tenant_summaries`) | ❌ |
| MachineTenantTagSummaries (`machine_tenant_tag_summaries`) | ❌ |
| TentacleUpgradesRequired (`tentacle_upgrades_required`) | ❌ |
| TotalDisabledMachines (`total_disabled_machines`) | ❌ |
| TotalMachines (`total_machines`) | ❌ |

### EnvironmentSummaryV2Resource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentTargetSummaries (`deployment_target_summaries`) | ❌ |
| Environment (`environment`) | ❌ |
| MachineEndpointSummaries (`machine_endpoint_summaries`) | ❌ |
| MachineHealthStatusSummaries (`machine_health_status_summaries`) | ❌ |
| MachineIdsForCalamariUpgrade (`machine_ids_for_calamari_upgrade`) | ❌ |
| MachineIdsForTentacleUpgrade (`machine_ids_for_tentacle_upgrade`) | ❌ |
| MachineTenantSummaries (`machine_tenant_summaries`) | ❌ |
| MachineTenantTagSummaries (`machine_tenant_tag_summaries`) | ❌ |
| TentacleUpgradesRequired (`tentacle_upgrades_required`) | ❌ |
| TotalDisabledMachines (`total_disabled_machines`) | ❌ |
| TotalMachines (`total_machines`) | ❌ |

### EnvironmentsSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentTargetSummaries (`deployment_target_summaries`) | ❌ |
| EnvironmentSummaries (`environment_summaries`) | ❌ |
| MachineEndpointSummaries (`machine_endpoint_summaries`) | ❌ |
| MachineHealthStatusSummaries (`machine_health_status_summaries`) | ❌ |
| MachineIdsForCalamariUpgrade (`machine_ids_for_calamari_upgrade`) | ❌ |
| MachineIdsForTentacleUpgrade (`machine_ids_for_tentacle_upgrade`) | ❌ |
| MachineTenantSummaries (`machine_tenant_summaries`) | ❌ |
| MachineTenantTagSummaries (`machine_tenant_tag_summaries`) | ❌ |
| TentacleUpgradesRequired (`tentacle_upgrades_required`) | ❌ |
| TotalDisabledMachines (`total_disabled_machines`) | ❌ |
| TotalMachines (`total_machines`) | ❌ |

### EventAgentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### EventCategoryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### EventGroupResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EventCategories (`event_categories`) | ❌ |
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### EventNotificationSubscriptionFilterResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DocumentTypes (`document_types`) | ❌ |
| Environments (`environments`) | ❌ |
| EventAgents (`event_agents`) | ❌ |
| EventCategories (`event_categories`) | ❌ |
| EventGroups (`event_groups`) | ❌ |
| ProjectGroups (`project_groups`) | ❌ |
| Projects (`projects`) | ❌ |
| Tags (`tags`) | ❌ |
| Tenants (`tenants`) | ❌ |
| Users (`users`) | ❌ |

### EventNotificationSubscriptionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EmailDigestLastProcessed (`email_digest_last_processed`) | ❌ |
| EmailDigestLastProcessedEventAutoId (`email_digest_last_processed_event_auto_id`) | ❌ |
| EmailFrequencyPeriod (`email_frequency_period`) | ❌ |
| EmailPriority (`email_priority`) | ❌ |
| EmailShowDatesInTimeZoneId (`email_show_dates_in_time_zone_id`) | ❌ |
| EmailTeams (`email_teams`) | ❌ |
| Filter (`filter`) | ❌ |
| WebhookHeaderKey (`webhook_header_key`) | ❌ |
| WebhookHeaderValue (`webhook_header_value`) | ❌ |
| WebhookLastProcessed (`webhook_last_processed`) | ❌ |
| WebhookLastProcessedEventAutoId (`webhook_last_processed_event_auto_id`) | ❌ |
| WebhookTeams (`webhook_teams`) | ❌ |
| WebhookTimeout (`webhook_timeout`) | ❌ |
| WebhookURI (`webhook_u_r_i`) | ❌ |

### EventResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Category (`category`) | ❌ |
| ChangeDetails (`change_details`) | ❌ |
| Comments (`comments`) | ❌ |
| Details (`details`) | ❌ |
| Id (`id`) | ❌ |
| IdentityEstablishedWith (`identity_established_with`) | ❌ |
| IpAddress (`ip_address`) | ❌ |
| IsService (`is_service`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Message (`message`) | ❌ |
| MessageHtml (`message_html`) | ❌ |
| MessageReferences (`message_references`) | ❌ |
| Occurred (`occurred`) | ❌ |
| RelatedDocumentIds (`related_document_ids`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| UserAgent (`user_agent`) | ❌ |
| UserId (`user_id`) | ❌ |
| Username (`username`) | ❌ |

### ExternalLinkResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Label (`label`) | ❌ |
| Uri (`uri`) | ❌ |

### FeatureToggleEnvironmentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentEnvironmentId (`deployment_environment_id`) | ❌ |
| ExcludedTenantIds (`excluded_tenant_ids`) | ❌ |
| ExcludedTenantTags (`excluded_tenant_tags`) | ❌ |
| FeatureToggleId (`feature_toggle_id`) | ❌ |
| IsEnabled (`is_enabled`) | ❌ |
| MinimumVersion (`minimum_version`) | ❌ |
| RolloutPercentage (`rollout_percentage`) | ❌ |
| Segments (`segments`) | ❌ |
| TenantIds (`tenant_ids`) | ❌ |
| TenantTags (`tenant_tags`) | ❌ |
| TenantTargetingStrategy (`tenant_targeting_strategy`) | ❌ |
| ToggleStrategy (`toggle_strategy`) | ❌ |

### FeatureToggleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DefaultIsEnabled (`default_is_enabled`) | ❌ |
| Description (`description`) | ❌ |
| Environments (`environments`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### FeaturesConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DefaultPowerShellEdition (`default_power_shell_edition`) | ❌ |
| HelpSidebarSupportLink (`help_sidebar_support_link`) | ❌ |
| HelpSidebarSupportLinkLabel (`help_sidebar_support_link_label`) | ❌ |
| Id (`id`) | ❌ |
| IsAutomaticStepUpdatesEnabled (`is_automatic_step_updates_enabled`) | ❌ |
| IsBuiltInWorkerEnabled (`is_built_in_worker_enabled`) | ❌ |
| IsCommunityActionTemplatesEnabled (`is_community_action_templates_enabled`) | ❌ |
| IsCompositeDockerHubRegistryFeedEnabled (`is_composite_docker_hub_registry_feed_enabled`) | ❌ |
| IsConfigureFeedsWithLocalOrSmbPathsEnabled (`is_configure_feeds_with_local_or_smb_paths_enabled`) | ❌ |
| IsExperimentalUIFeatureEnabled (`is_experimental_u_i_feature_enabled`) | ❌ |
| IsHelpSidebarEnabled (`is_help_sidebar_enabled`) | ❌ |
| IsKubernetesCloudTargetDiscoveryEnabled (`is_kubernetes_cloud_target_discovery_enabled`) | ❌ |
| IsProjectsPageOnboardingEnabled (`is_projects_page_onboarding_enabled`) | ❌ |
| IsProjectsPageOptimizationEnabled (`is_projects_page_optimization_enabled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### FeedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| FeedType (`feed_type`) | — |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| PackageAcquisitionLocationOptions (`package_acquisition_location_options`) | ✅ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |

### GitBranchResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanonicalName (`canonical_name`) | ❌ |
| Id (`id`) | ❌ |
| IsProtected (`is_protected`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### GitCommitResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanonicalName (`canonical_name`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### GitCredentialRepositoryRestrictionsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AllowedRepositories (`allowed_repositories`) | ❌ |
| Enabled (`enabled`) | ❌ |

### GitCredentialResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Details (`details`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| RepositoryRestrictions (`repository_restrictions`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### GitDependencyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DefaultBranch (`default_branch`) | ❌ |
| FilePathFilters (`file_path_filters`) | ❌ |
| GitCredentialId (`git_credential_id`) | ❌ |
| GitCredentialType (`git_credential_type`) | ❌ |
| GitHubConnectionId (`git_hub_connection_id`) | ❌ |
| Name (`name`) | ❌ |
| RepositoryUri (`repository_uri`) | ❌ |
| StepPackageInputsReferenceId (`step_package_inputs_reference_id`) | ❌ |

### GitNamedRefByNameResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanonicalName (`canonical_name`) | ❌ |
| Id (`id`) | ❌ |
| IsProtected (`is_protected`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### GitPersistenceSettingsConversionStateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| RunbooksAreInGit (`runbooks_are_in_git`) | ❌ |
| VariablesAreInGit (`variables_are_in_git`) | ❌ |

### GitPersistenceSettingsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| BasePath (`base_path`) | ❌ |
| ConversionState (`conversion_state`) | ❌ |
| Credentials (`credentials`) | ❌ |
| DefaultBranch (`default_branch`) | ❌ |
| ProtectedBranchNamePatterns (`protected_branch_name_patterns`) | ❌ |
| ProtectedDefaultBranch (`protected_default_branch`) | ❌ |
| Type (`type`) | ❌ |
| Url (`url`) | ❌ |

### GitReferenceResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| GitCommit (`git_commit`) | ❌ |
| GitRef (`git_ref`) | ❌ |

### GitTagResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanonicalName (`canonical_name`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### IconResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Color (`color`) | ❌ |
| Id (`id`) | ❌ |

### IdentityClaimResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| IsIdentifyingClaim (`is_identifying_claim`) | ❌ |
| Value (`value`) | ❌ |

### IdentityMetadataResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ClaimDescriptors (`claim_descriptors`) | ❌ |
| IdentityProviderName (`identity_provider_name`) | ❌ |
| Links (`links`) | ❌ |
| ScimEnabled (`scim_enabled`) | ❌ |

### IdentityResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Claims (`claims`) | — |
| IdentityProviderName (`identity_provider_name`) | ❌ |

### InsightsDataSeriesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Intervals (`intervals`) | ❌ |
| Name (`name`) | ❌ |

### InsightsEnvironmentGroupResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Environments (`environments`) | ❌ |
| Name (`name`) | ❌ |

### InsightsFailureRateMetricResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentFailure (`deployment_failure`) | ❌ |
| Failed (`failed`) | — |
| Rate (`rate`) | — |
| Successful (`successful`) | ❌ |
| SuccessfulButHadGuidedFailure (`successful_but_had_guided_failure`) | ❌ |
| Total (`total`) | — |

### InsightsReportResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ChannelIds (`channel_ids`) | ❌ |
| Description (`description`) | ❌ |
| EnvironmentGroups (`environment_groups`) | ❌ |
| IconColor (`icon_color`) | ❌ |
| IconId (`icon_id`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| ProjectGroupIds (`project_group_ids`) | ❌ |
| ProjectIds (`project_ids`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| TenantIds (`tenant_ids`) | ❌ |
| TenantMode (`tenant_mode`) | ❌ |
| TenantTags (`tenant_tags`) | ❌ |
| TimeZone (`time_zone`) | ❌ |

### InterruptionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanTakeResponsibility (`can_take_responsibility`) | ❌ |
| CorrelationId (`correlation_id`) | ❌ |
| Created (`created`) | ❌ |
| Form (`form`) | ❌ |
| HasResponsibility (`has_responsibility`) | ❌ |
| Id (`id`) | ❌ |
| IsLinkedToOtherInterruption (`is_linked_to_other_interruption`) | ❌ |
| IsPending (`is_pending`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| RelatedDocumentIds (`related_document_ids`) | ❌ |
| ResponsibleTeamIds (`responsible_team_ids`) | ❌ |
| ResponsibleUserId (`responsible_user_id`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| TaskId (`task_id`) | ❌ |
| Title (`title`) | ❌ |
| Type (`type`) | ❌ |

### InvitationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AddToTeamIds (`add_to_team_ids`) | ❌ |
| Expires (`expires`) | ❌ |
| Id (`id`) | ❌ |
| InvitationCode (`invitation_code`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### KubernetesEventResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Action (`action`) | ❌ |
| Count (`count`) | ❌ |
| FirstObservedTime (`first_observed_time`) | ❌ |
| LastObservedTime (`last_observed_time`) | ❌ |
| Manifest (`manifest`) | ❌ |
| Note (`note`) | ❌ |
| Reason (`reason`) | ❌ |
| ReportingController (`reporting_controller`) | ❌ |
| ReportingInstance (`reporting_instance`) | ❌ |
| Type (`type`) | ❌ |

### KubernetesLiveStatusDetailedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Children (`children`) | ❌ |
| DesiredResourceId (`desired_resource_id`) | ❌ |
| ExternalLink (`external_link`) | ❌ |
| HealthStatus (`health_status`) | ❌ |
| Kind (`kind`) | ❌ |
| LastUpdated (`last_updated`) | ❌ |
| ManifestSummary (`manifest_summary`) | ❌ |
| Name (`name`) | ❌ |
| Namespace (`namespace`) | ❌ |
| ResourceId (`resource_id`) | ❌ |
| ResourceSourceId (`resource_source_id`) | ❌ |
| SourceType (`source_type`) | ❌ |
| SyncStatus (`sync_status`) | ❌ |

### KubernetesLiveStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Children (`children`) | ❌ |
| DesiredResourceId (`desired_resource_id`) | ❌ |
| Group (`group`) | ❌ |
| HealthStatus (`health_status`) | ❌ |
| Kind (`kind`) | ❌ |
| LastUpdated (`last_updated`) | ❌ |
| MachineId (`machine_id`) | ❌ |
| Name (`name`) | ❌ |
| Namespace (`namespace`) | ❌ |
| ResourceId (`resource_id`) | ❌ |
| ResourceSourceId (`resource_source_id`) | ❌ |
| SourceType (`source_type`) | ❌ |
| SyncStatus (`sync_status`) | ❌ |

### KubernetesMachineLiveStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| MachineId (`machine_id`) | ❌ |
| Resources (`resources`) | ❌ |
| Status (`status`) | ❌ |

### LetsEncryptConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AcceptLetsEncryptTermsOfService (`accept_lets_encrypt_terms_of_service`) | ❌ |
| CertificateExpiryDate (`certificate_expiry_date`) | ❌ |
| CertificateThumbprint (`certificate_thumbprint`) | ❌ |
| DnsName (`dns_name`) | ❌ |
| Enabled (`enabled`) | ❌ |
| HttpsPort (`https_port`) | ❌ |
| IPAddress (`i_p_address`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Path (`path`) | ❌ |
| RegistrationEmailAddress (`registration_email_address`) | ❌ |

### LibraryVariableSetResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ContentType (`content_type`) | ❌ |
| Description (`description`) | ✅ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| SpaceId (`space_id`) | ✅ |
| Templates (`templates`) | ❌ |
| VariableSetId (`variable_set_id`) | ✅ |
| Version (`version`) | ❌ |

### LicenseLimitStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CurrentUsage (`current_usage`) | ❌ |
| Disposition (`disposition`) | ❌ |
| EffectiveLimit (`effective_limit`) | ❌ |
| EffectiveLimitDescription (`effective_limit_description`) | ❌ |
| IsUnlimited (`is_unlimited`) | ❌ |
| LicenseLimitDescription (`license_limit_description`) | ❌ |
| LicensedLimit (`licensed_limit`) | ❌ |
| LimitStatus (`limit_status`) | ❌ |
| Message (`message`) | ❌ |
| Name (`name`) | ❌ |
| TargetTypes (`target_types`) | ❌ |

### LicenseLimitUsageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CurrentUsage (`current_usage`) | ❌ |
| Disposition (`disposition`) | ❌ |
| EffectiveLimit (`effective_limit`) | ❌ |
| EffectiveLimitDescription (`effective_limit_description`) | ❌ |
| IsUnlimited (`is_unlimited`) | ❌ |
| LicenseLimitDescription (`license_limit_description`) | ❌ |
| LicensedLimit (`licensed_limit`) | ❌ |
| LimitStatus (`limit_status`) | ❌ |
| LimitUsageDescription (`limit_usage_description`) | ❌ |
| Message (`message`) | ❌ |
| Name (`name`) | ❌ |
| TargetTypes (`target_types`) | ❌ |

### LicenseMessageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Disposition (`disposition`) | ❌ |
| Message (`message`) | ❌ |
| MessagePolicy (`message_policy`) | ❌ |

### LicenseResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LicenseText (`license_text`) | ❌ |
| Links (`links`) | ❌ |
| SerialNumber (`serial_number`) | ❌ |

### LicenseStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ComplianceSummary (`compliance_summary`) | ❌ |
| DaysToEffectiveExpiryDate (`days_to_effective_expiry_date`) | ❌ |
| DoesExpiryBlockKeyActivities (`does_expiry_block_key_activities`) | ❌ |
| EffectiveClusterTaskLimit (`effective_cluster_task_limit`) | ❌ |
| EffectiveExpiryDate (`effective_expiry_date`) | ❌ |
| EffectiveNodeTaskLimit (`effective_node_task_limit`) | ❌ |
| EffectiveStartDate (`effective_start_date`) | ❌ |
| HostingEnvironment (`hosting_environment`) | ❌ |
| Id (`id`) | ❌ |
| IsClusterTaskLimitControlledByLicense (`is_cluster_task_limit_controlled_by_license`) | ❌ |
| IsCompliant (`is_compliant`) | ❌ |
| IsInitialisationLicense (`is_initialisation_license`) | ❌ |
| IsNodeTaskLimitControlledByLicense (`is_node_task_limit_controlled_by_license`) | ❌ |
| IsPtm (`is_ptm`) | ❌ |
| IsTrial (`is_trial`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Limits (`limits`) | ❌ |
| Links (`links`) | ❌ |
| Messages (`messages`) | ❌ |
| PermissionsMode (`permissions_mode`) | ❌ |
| SerialNumber (`serial_number`) | ❌ |

### LicenseUsageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsPtm (`is_ptm`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Limits (`limits`) | ❌ |
| Links (`links`) | ❌ |
| SpacesUsage (`spaces_usage`) | ❌ |

### LifecycleProgressionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| NextDeployments (`next_deployments`) | ❌ |
| NextDeploymentsMinimumRequired (`next_deployments_minimum_required`) | ❌ |
| Phases (`phases`) | ❌ |

### LifecycleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| Phases (`phases`) | — |
| ReleaseRetentionPolicy (`release_retention_policy`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| TentacleRetentionPolicy (`tentacle_retention_policy`) | ❌ |

### LiveStatusSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| LastUpdated (`last_updated`) | ❌ |
| Status (`status`) | ❌ |

### LoadBalancerPingResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsInMaintenanceMode (`is_in_maintenance_mode`) | ❌ |
| IsOffline (`is_offline`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastSeen (`last_seen`) | ❌ |
| Links (`links`) | ❌ |
| MaxConcurrentTasks (`max_concurrent_tasks`) | ❌ |
| Name (`name`) | ❌ |
| Version (`version`) | ❌ |

### LoginInitiatedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProviderName (`provider_name`) | ❌ |
| WasLoginInitiated (`was_login_initiated`) | ❌ |

### MachineBasedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Architecture (`architecture`) | ❌ |
| Endpoint (`endpoint`) | ❌ |
| HasLatestCalamari (`has_latest_calamari`) | ❌ |
| HealthStatus (`health_status`) | ❌ |
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| IsInProcess (`is_in_process`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MachinePolicyId (`machine_policy_id`) | ❌ |
| Name (`name`) | ❌ |
| OperatingSystem (`operating_system`) | ❌ |
| ShellName (`shell_name`) | ❌ |
| ShellVersion (`shell_version`) | ❌ |
| SkipInitialHealthCheck (`skip_initial_health_check`) | ❌ |
| Slug (`slug`) | ❌ |
| StatusSummary (`status_summary`) | ❌ |
| Thumbprint (`thumbprint`) | ❌ |
| Uri (`uri`) | ❌ |

### MachineCleanupPolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeleteMachinesBehavior (`delete_machines_behavior`) | ❌ |
| DeleteMachinesElapsedTimeSpan (`delete_machines_elapsed_time_span`) | ❌ |

### MachineConnectivityPolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| MachineConnectivityBehavior (`machine_connectivity_behavior`) | ❌ |

### MachineHealthCheckPolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| BashHealthCheckPolicy (`bash_health_check_policy`) | ❌ |
| HealthCheckCron (`health_check_cron`) | ❌ |
| HealthCheckCronTimezone (`health_check_cron_timezone`) | ❌ |
| HealthCheckInterval (`health_check_interval`) | ❌ |
| HealthCheckType (`health_check_type`) | ❌ |
| PowerShellHealthCheckPolicy (`power_shell_health_check_policy`) | ❌ |

### MachinePackageCacheRetentionPolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| PackageUnit (`package_unit`) | ❌ |
| QuantityOfPackagesToKeep (`quantity_of_packages_to_keep`) | ❌ |
| QuantityOfVersionsToKeep (`quantity_of_versions_to_keep`) | ❌ |
| Strategy (`strategy`) | ❌ |
| VersionUnit (`version_unit`) | ❌ |

### MachinePolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ConnectionConnectTimeout (`connection_connect_timeout`) | ❌ |
| ConnectionRetryCountLimit (`connection_retry_count_limit`) | ❌ |
| ConnectionRetrySleepInterval (`connection_retry_sleep_interval`) | ❌ |
| ConnectionRetryTimeLimit (`connection_retry_time_limit`) | ❌ |
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| IsDefault (`is_default`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MachineCleanupPolicy (`machine_cleanup_policy`) | ❌ |
| MachineConnectivityPolicy (`machine_connectivity_policy`) | ❌ |
| MachineHealthCheckPolicy (`machine_health_check_policy`) | ❌ |
| MachinePackageCacheRetentionPolicy (`machine_package_cache_retention_policy`) | ❌ |
| MachineRpcCallRetryPolicy (`machine_rpc_call_retry_policy`) | ❌ |
| MachineUpdatePolicy (`machine_update_policy`) | ❌ |
| Name (`name`) | ❌ |
| PollingRequestQueueTimeout (`polling_request_queue_timeout`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### MachineResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Architecture (`architecture`) | ❌ |
| Endpoint (`endpoint`) | ❌ |
| EnvironmentIds (`environment_ids`) | ❌ |
| HasLatestCalamari (`has_latest_calamari`) | ❌ |
| HealthStatus (`health_status`) | ❌ |
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| IsInProcess (`is_in_process`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MachinePolicyId (`machine_policy_id`) | ❌ |
| Name (`name`) | ❌ |
| OperatingSystem (`operating_system`) | ❌ |
| Roles (`roles`) | ❌ |
| ShellName (`shell_name`) | ❌ |
| ShellVersion (`shell_version`) | ❌ |
| SkipInitialHealthCheck (`skip_initial_health_check`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| StatusSummary (`status_summary`) | ❌ |
| TenantIds (`tenant_ids`) | ❌ |
| TenantTags (`tenant_tags`) | ❌ |
| TenantedDeploymentParticipation (`tenanted_deployment_participation`) | ❌ |
| Thumbprint (`thumbprint`) | ❌ |
| Uri (`uri`) | ❌ |

### MachineRpcCallRetryPolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Enabled (`enabled`) | ❌ |
| HealthCheckRetryDuration (`health_check_retry_duration`) | ❌ |
| RetryDuration (`retry_duration`) | ❌ |

### MachineScriptPolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| RunType (`run_type`) | ❌ |
| ScriptBody (`script_body`) | ❌ |

### MachineUpdatePolicyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CalamariUpdateBehavior (`calamari_update_behavior`) | ❌ |
| KubernetesAgentUpdateBehavior (`kubernetes_agent_update_behavior`) | ❌ |
| TentacleUpdateAccountId (`tentacle_update_account_id`) | ❌ |
| TentacleUpdateBehavior (`tentacle_update_behavior`) | ❌ |

### MaintenanceConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsInMaintenanceMode (`is_in_maintenance_mode`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### ManifestSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Annotations (`annotations`) | ❌ |
| CreationTimestamp (`creation_timestamp`) | ❌ |
| Kind (`kind`) | ❌ |
| Labels (`labels`) | ❌ |

### MigrationImportResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeletePackageOnCompletion (`delete_package_on_completion`) | ❌ |
| FailureCallbackUri (`failure_callback_uri`) | ❌ |
| Id (`id`) | ❌ |
| IsDryRun (`is_dry_run`) | ❌ |
| IsEncryptedPackage (`is_encrypted_package`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| OverwriteExisting (`overwrite_existing`) | ❌ |
| PackageFeedSpaceId (`package_feed_space_id`) | ❌ |
| PackageId (`package_id`) | ❌ |
| PackageVersion (`package_version`) | ❌ |
| Password (`password`) | ❌ |
| SuccessCallbackUri (`success_callback_uri`) | ❌ |
| TaskId (`task_id`) | ❌ |

### MigrationPartialExportResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DestinationApiKey (`destination_api_key`) | ❌ |
| DestinationPackageFeed (`destination_package_feed`) | ❌ |
| DestinationPackageFeedSpaceId (`destination_package_feed_space_id`) | ❌ |
| EncryptPackage (`encrypt_package`) | ❌ |
| FailureCallbackUri (`failure_callback_uri`) | ❌ |
| Id (`id`) | ❌ |
| IgnoreCertificates (`ignore_certificates`) | ❌ |
| IgnoreDeployments (`ignore_deployments`) | ❌ |
| IgnoreMachines (`ignore_machines`) | ❌ |
| IgnoreTenants (`ignore_tenants`) | ❌ |
| IncludeTaskLogs (`include_task_logs`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| PackageId (`package_id`) | ❌ |
| PackageVersion (`package_version`) | ❌ |
| Password (`password`) | ❌ |
| Projects (`projects`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| SuccessCallbackUri (`success_callback_uri`) | ❌ |
| TaskId (`task_id`) | ❌ |

### MissingVariableResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EnvironmentId (`environment_id`) | ❌ |
| LibraryVariableSetId (`library_variable_set_id`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| VariableTemplateId (`variable_template_id`) | ❌ |
| VariableTemplateName (`variable_template_name`) | ❌ |

### MonitorErrorResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ErrorCode (`error_code`) | ❌ |
| ErrorMessage (`error_message`) | ❌ |

### MultiTenancyStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Enabled (`enabled`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### OctopusPackageVersionBuildInformationMappedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Branch (`branch`) | ❌ |
| BuildEnvironment (`build_environment`) | ❌ |
| BuildNumber (`build_number`) | ❌ |
| BuildUrl (`build_url`) | ❌ |
| Commits (`commits`) | ❌ |
| Created (`created`) | ❌ |
| Id (`id`) | ❌ |
| IncompleteDataWarning (`incomplete_data_warning`) | ❌ |
| IssueTrackerName (`issue_tracker_name`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| PackageId (`package_id`) | ❌ |
| VcsCommitNumber (`vcs_commit_number`) | ❌ |
| VcsCommitUrl (`vcs_commit_url`) | ❌ |
| VcsRoot (`vcs_root`) | ❌ |
| VcsType (`vcs_type`) | ❌ |
| Version (`version`) | ❌ |
| WorkItems (`work_items`) | ❌ |

### OctopusServerClusterSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Links (`links`) | ❌ |
| Nodes (`nodes`) | ❌ |

### OctopusServerNodeDetailsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| RunningTasks (`running_tasks`) | ❌ |

### OctopusServerNodeResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsInMaintenanceMode (`is_in_maintenance_mode`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MaxConcurrentTasks (`max_concurrent_tasks`) | ❌ |
| Name (`name`) | ❌ |

### OctopusServerNodeSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsInMaintenanceMode (`is_in_maintenance_mode`) | ❌ |
| IsOffline (`is_offline`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastSeen (`last_seen`) | ❌ |
| Links (`links`) | ❌ |
| MaxConcurrentTasks (`max_concurrent_tasks`) | ❌ |
| MaxSqlConnectionPoolSize (`max_sql_connection_pool_size`) | ❌ |
| Name (`name`) | ❌ |
| RecommendedMaxSqlConnectionPoolSize (`recommended_max_sql_connection_pool_size`) | — |
| RunningTaskCount (`running_task_count`) | ❌ |
| Version (`version`) | ❌ |

### PackageDescriptionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| LatestVersion (`latest_version`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### PackageFromBuiltInFeedResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| FeedId (`feed_id`) | ❌ |
| FileExtension (`file_extension`) | ❌ |
| Hash (`hash`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| NuGetFeedId (`nu_get_feed_id`) | ❌ |
| NuGetPackageId (`nu_get_package_id`) | ❌ |
| PackageId (`package_id`) | ❌ |
| PackageSizeBytes (`package_size_bytes`) | ❌ |
| PackageVersionBuildInformation (`package_version_build_information`) | ❌ |
| Published (`published`) | ❌ |
| ReleaseNotes (`release_notes`) | ❌ |
| Summary (`summary`) | ❌ |
| Title (`title`) | ❌ |
| Version (`version`) | ❌ |

### PackageNoteListResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Packages (`packages`) | ❌ |

### PackageReferenceResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AcquisitionLocation (`acquisition_location`) | ❌ |
| FeedId (`feed_id`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| PackageId (`package_id`) | ❌ |
| Properties (`properties`) | ❌ |
| StepPackageInputsReferenceId (`step_package_inputs_reference_id`) | ❌ |
| Version (`version`) | ❌ |

### PackageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| FeedId (`feed_id`) | ❌ |
| FileExtension (`file_extension`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| NuGetFeedId (`nu_get_feed_id`) | ❌ |
| NuGetPackageId (`nu_get_package_id`) | ❌ |
| PackageId (`package_id`) | ❌ |
| PackageVersionBuildInformation (`package_version_build_information`) | ❌ |
| Published (`published`) | ❌ |
| ReleaseNotes (`release_notes`) | ❌ |
| Summary (`summary`) | ❌ |
| Title (`title`) | ❌ |
| Version (`version`) | ❌ |

### PackageSignatureResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| BaseVersion (`base_version`) | ❌ |
| Signature (`signature`) | ❌ |

### PackageVersionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| FeedId (`feed_id`) | ❌ |
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| PackageId (`package_id`) | ❌ |
| Published (`published`) | ❌ |
| ReleaseNotes (`release_notes`) | ❌ |
| SizeBytes (`size_bytes`) | ❌ |
| Title (`title`) | ❌ |
| Version (`version`) | ❌ |

### PerformanceConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DefaultDashboardRenderMode (`default_dashboard_render_mode`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### PhaseDeploymentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Deployment (`deployment`) | ❌ |
| Task (`task`) | ❌ |

### PhaseProgressionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AutomaticDeploymentTargets (`automatic_deployment_targets`) | ❌ |
| Blocked (`blocked`) | ❌ |
| Deployments (`deployments`) | ❌ |
| Id (`id`) | ❌ |
| IsOptionalPhase (`is_optional_phase`) | ❌ |
| IsPriorityPhase (`is_priority_phase`) | ❌ |
| MinimumEnvironmentsBeforePromotion (`minimum_environments_before_promotion`) | ❌ |
| Name (`name`) | ❌ |
| OptionalDeploymentTargets (`optional_deployment_targets`) | ❌ |
| Progress (`progress`) | ❌ |

### PhaseResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AutomaticDeploymentTargets (`automatic_deployment_targets`) | ❌ |
| Id (`id`) | ❌ |
| IsOptionalPhase (`is_optional_phase`) | ❌ |
| IsPriorityPhase (`is_priority_phase`) | ❌ |
| MinimumEnvironmentsBeforePromotion (`minimum_environments_before_promotion`) | ❌ |
| Name (`name`) | ❌ |
| OptionalDeploymentTargets (`optional_deployment_targets`) | ❌ |
| ReleaseRetentionPolicy (`release_retention_policy`) | ❌ |
| TentacleRetentionPolicy (`tentacle_retention_policy`) | ❌ |

### PlatformHubAccountResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Details (`details`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| Slug (`slug`) | ❌ |

### PlatformHubGitCredentialResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Details (`details`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| RepositoryRestrictions (`repository_restrictions`) | ❌ |

### PlatformHubVersionControlSettingsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| BasePath (`base_path`) | ❌ |
| Credentials (`credentials`) | ❌ |
| DefaultBranch (`default_branch`) | ❌ |
| Url (`url`) | ❌ |

### ProcessTemplateParameterResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DisplaySettings (`display_settings`) | ❌ |
| HelpText (`help_text`) | ❌ |
| IsOptional (`is_optional`) | ❌ |
| Label (`label`) | ❌ |
| Name (`name`) | ❌ |
| Values (`values`) | ❌ |

### ProcessTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| GitRef (`git_ref`) | ❌ |
| Icon (`icon`) | ❌ |
| Id (`id`) | — |
| Name (`name`) | ❌ |
| Parameters (`parameters`) | — |
| Slug (`slug`) | ❌ |
| Steps (`steps`) | — |

### ProcessTemplateVersionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| GitCommit (`git_commit`) | ❌ |
| GitRef (`git_ref`) | ❌ |
| Icon (`icon`) | ❌ |
| Id (`id`) | ❌ |
| IsPreRelease (`is_pre_release`) | ❌ |
| Name (`name`) | ❌ |
| Parameters (`parameters`) | ❌ |
| PublishedDate (`published_date`) | ❌ |
| Slug (`slug`) | ❌ |
| Steps (`steps`) | ❌ |
| Version (`version`) | ❌ |

### ProgressionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ChannelEnvironments (`channel_environments`) | ❌ |
| Environments (`environments`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LifecycleEnvironments (`lifecycle_environments`) | ❌ |
| Links (`links`) | ❌ |
| Releases (`releases`) | ❌ |

### ProjectGroupResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ✅ |
| EnvironmentIds (`environment_ids`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| RetentionPolicyId (`retention_policy_id`) | ❌ |
| Slug (`slug`) | ✅ |
| SpaceId (`space_id`) | ✅ |

### ProjectResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AllowIgnoreChannelRules (`allow_ignore_channel_rules`) | ❌ |
| AutoCreateRelease (`auto_create_release`) | ❌ |
| AutoDeployReleaseOverrides (`auto_deploy_release_overrides`) | — |
| ClonedFromProjectId (`cloned_from_project_id`) | ❌ |
| CombineHealthAndSyncStatusInDashboardLiveStatus (`combine_health_and_sync_status_in_dashboard_live_status`) | ❌ |
| DefaultGuidedFailureMode (`default_guided_failure_mode`) | ❌ |
| DefaultPowerShellEdition (`default_power_shell_edition`) | ❌ |
| DefaultToSkipIfAlreadyInstalled (`default_to_skip_if_already_installed`) | ❌ |
| DeploymentChangesTemplate (`deployment_changes_template`) | ❌ |
| DeploymentProcessId (`deployment_process_id`) | ❌ |
| DeprovisioningRunbookId (`deprovisioning_runbook_id`) | ❌ |
| Description (`description`) | ❌ |
| DiscreteChannelRelease (`discrete_channel_release`) | ❌ |
| ExecuteDeploymentsOnResilientPipeline (`execute_deployments_on_resilient_pipeline`) | ❌ |
| ExtensionSettings (`extension_settings`) | — |
| ForcePackageDownload (`force_package_download`) | ❌ |
| Icon (`icon`) | ❌ |
| Id (`id`) | ❌ |
| IncludedLibraryVariableSetIds (`included_library_variable_set_ids`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| IsVersionControlled (`is_version_controlled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LifecycleId (`lifecycle_id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| PersistenceSettings (`persistence_settings`) | ❌ |
| ProjectConnectivityPolicy (`project_connectivity_policy`) | ❌ |
| ProjectGroupId (`project_group_id`) | ❌ |
| ProjectTags (`project_tags`) | ❌ |
| ProjectTemplateDetails (`project_template_details`) | ❌ |
| ProvisioningRunbookId (`provisioning_runbook_id`) | ❌ |
| ReleaseCreationStrategy (`release_creation_strategy`) | ❌ |
| ReleaseNotesTemplate (`release_notes_template`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Templates (`templates`) | ❌ |
| TenantedDeploymentMode (`tenanted_deployment_mode`) | ❌ |
| VariableSetId (`variable_set_id`) | ❌ |
| VersioningStrategy (`versioning_strategy`) | ❌ |

### ProjectTriggerResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Action (`action`) | ❌ |
| Description (`description`) | ❌ |
| Filter (`filter`) | ❌ |
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### ProjectVariablesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| OwnerId (`owner_id`) | ❌ |
| ScopeValues (`scope_values`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Variables (`variables`) | ❌ |
| Version (`version`) | ❌ |

### PropertyValueResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| IsSensitive (`is_sensitive`) | ❌ |
| SensitiveValue (`sensitive_value`) | ❌ |
| Value (`value`) | ❌ |

### ProxyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Host (`host`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| Password (`password`) | ❌ |
| Port (`port`) | ❌ |
| ProxyType (`proxy_type`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Username (`username`) | ❌ |

### RecurringScheduleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EndAfterOccurrences (`end_after_occurrences`) | ❌ |
| EndDate (`end_date`) | ❌ |
| EndOnDate (`end_on_date`) | ❌ |
| EndType (`end_type`) | ❌ |
| StartDate (`start_date`) | ❌ |
| Type (`type`) | — |
| Unit (`unit`) | ❌ |
| UserUtcOffsetInMinutes (`user_utc_offset_in_minutes`) | ❌ |

### ReleaseChangesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| BuildInformation (`build_information`) | ❌ |
| Commits (`commits`) | ❌ |
| ReleaseNotes (`release_notes`) | ❌ |
| Version (`version`) | ❌ |
| WorkItems (`work_items`) | ❌ |

### ReleaseCreationStrategyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ChannelId (`channel_id`) | ❌ |
| ReleaseCreationPackage (`release_creation_package`) | ❌ |

### ReleasePackageVersionBuildInformationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Branch (`branch`) | ❌ |
| BuildEnvironment (`build_environment`) | ❌ |
| BuildNumber (`build_number`) | ❌ |
| BuildUrl (`build_url`) | ❌ |
| Commits (`commits`) | ❌ |
| IssueTrackerName (`issue_tracker_name`) | ❌ |
| PackageId (`package_id`) | ❌ |
| VcsCommitNumber (`vcs_commit_number`) | ❌ |
| VcsCommitUrl (`vcs_commit_url`) | ❌ |
| VcsRoot (`vcs_root`) | ❌ |
| VcsType (`vcs_type`) | ❌ |
| Version (`version`) | ❌ |
| WorkItems (`work_items`) | ❌ |

### ReleaseProgressionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Channel (`channel`) | ❌ |
| Deployments (`deployments`) | ❌ |
| HasUnresolvedDefect (`has_unresolved_defect`) | ❌ |
| NextDeployments (`next_deployments`) | ❌ |
| Release (`release`) | ❌ |
| ReleaseRetentionPeriod (`release_retention_period`) | ❌ |
| TentacleRetentionPeriod (`tentacle_retention_period`) | ❌ |

### ReleaseResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Assembled (`assembled`) | ❌ |
| BuildInformation (`build_information`) | ❌ |
| ChannelId (`channel_id`) | ❌ |
| CustomFields (`custom_fields`) | ❌ |
| Id (`id`) | ❌ |
| IgnoreChannelRules (`ignore_channel_rules`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LibraryVariableSetSnapshotIds (`library_variable_set_snapshot_ids`) | ❌ |
| Links (`links`) | ❌ |
| ProjectDeploymentProcessSnapshotId (`project_deployment_process_snapshot_id`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| ProjectVariableSetSnapshotId (`project_variable_set_snapshot_id`) | ❌ |
| ReleaseNotes (`release_notes`) | ❌ |
| SelectedGitResources (`selected_git_resources`) | ❌ |
| SelectedPackages (`selected_packages`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Version (`version`) | ❌ |
| VersionControlReference (`version_control_reference`) | ❌ |

### ReleaseTemplateGitResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionName (`action_name`) | ❌ |
| DefaultBranch (`default_branch`) | ❌ |
| FilePathFilters (`file_path_filters`) | ❌ |
| GitCredentialId (`git_credential_id`) | ❌ |
| GitHubConnectionId (`git_hub_connection_id`) | ❌ |
| GitResourceSelectedLastRelease (`git_resource_selected_last_release`) | ❌ |
| IsResolvable (`is_resolvable`) | ❌ |
| Name (`name`) | ❌ |
| RepositoryUri (`repository_uri`) | ❌ |

### ReleaseTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DeploymentProcessId (`deployment_process_id`) | ❌ |
| GitResources (`git_resources`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastReleaseVersion (`last_release_version`) | ❌ |
| Links (`links`) | ❌ |
| NextVersionIncrement (`next_version_increment`) | ❌ |
| Packages (`packages`) | ❌ |
| VersioningPackageReferenceName (`versioning_package_reference_name`) | ❌ |
| VersioningPackageStepName (`versioning_package_step_name`) | ❌ |

### ReportDeploymentCountOverTimeResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| ReportData (`report_data`) | ❌ |

### ResolvedProjectTemplateDetailsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| IsShared (`is_shared`) | ❌ |
| Slug (`slug`) | ❌ |
| Version (`version`) | ❌ |
| VersionMask (`version_mask`) | ❌ |

### RetentionDefaultConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| RetentionDays (`retention_days`) | ❌ |

### RootResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ApiVersion (`api_version`) | ❌ |
| Application (`application`) | ❌ |
| HasLongTermSupport (`has_long_term_support`) | — |
| Id (`id`) | ❌ |
| InstallationId (`installation_id`) | ❌ |
| IsEarlyAccessProgram (`is_early_access_program`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Version (`version`) | ❌ |

### RunbookProcessResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastSnapshotId (`last_snapshot_id`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| RunbookId (`runbook_id`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Steps (`steps`) | — |
| Version (`version`) | ❌ |

### RunbookResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ConnectivityPolicy (`connectivity_policy`) | ✅ |
| DefaultGuidedFailureMode (`default_guided_failure_mode`) | ✅ |
| Description (`description`) | ✅ |
| EnvironmentScope (`environment_scope`) | ✅ |
| Environments (`environments`) | ✅ |
| FailTargetDiscovery (`fail_target_discovery`) | ❌ |
| ForcePackageDownload (`force_package_download`) | ✅ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MultiTenancyMode (`multi_tenancy_mode`) | ✅ |
| Name (`name`) | ✅ |
| ProjectId (`project_id`) | ✅ |
| PublishedRunbookSnapshotId (`published_runbook_snapshot_id`) | ✅ |
| RunRetentionPolicy (`run_retention_policy`) | ❌ |
| RunbookProcessId (`runbook_process_id`) | ✅ |
| RunbookTags (`runbook_tags`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |

### RunbookRunPreviewResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Form (`form`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| StepsToExecute (`steps_to_execute`) | ❌ |
| UseGuidedFailureModeByDefault (`use_guided_failure_mode_by_default`) | ❌ |

### RunbookRunResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ChangeRequestSettings (`change_request_settings`) | ❌ |
| Comments (`comments`) | ❌ |
| Created (`created`) | ❌ |
| DebugMode (`debug_mode`) | ❌ |
| DeployedBy (`deployed_by`) | ❌ |
| DeployedById (`deployed_by_id`) | ❌ |
| DeployedToMachineIds (`deployed_to_machine_ids`) | ❌ |
| EnvironmentId (`environment_id`) | ❌ |
| ExcludedMachineIds (`excluded_machine_ids`) | ❌ |
| ExecutionPlanLogContext (`execution_plan_log_context`) | ❌ |
| FailTargetDiscovery (`fail_target_discovery`) | ❌ |
| FailureEncountered (`failure_encountered`) | ❌ |
| ForcePackageDownload (`force_package_download`) | ❌ |
| FormValues (`form_values`) | ❌ |
| FrozenRunbookProcessId (`frozen_runbook_process_id`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ManifestVariableSetId (`manifest_variable_set_id`) | ❌ |
| Name (`name`) | ❌ |
| Priority (`priority`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| QueueTime (`queue_time`) | ❌ |
| QueueTimeExpiry (`queue_time_expiry`) | ❌ |
| RunbookId (`runbook_id`) | ❌ |
| RunbookName (`runbook_name`) | ❌ |
| RunbookSnapshotId (`runbook_snapshot_id`) | ❌ |
| SkipActions (`skip_actions`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| SpecificMachineIds (`specific_machine_ids`) | ❌ |
| TaskId (`task_id`) | ❌ |
| TenantId (`tenant_id`) | ❌ |
| TentacleRetentionPeriod (`tentacle_retention_period`) | ❌ |
| UseGuidedFailure (`use_guided_failure`) | ❌ |

### RunbookRunTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsLibraryVariableSetModified (`is_library_variable_set_modified`) | ❌ |
| IsRunbookProcessModified (`is_runbook_process_modified`) | ❌ |
| IsVariableSetModified (`is_variable_set_modified`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| PromoteTo (`promote_to`) | ❌ |
| TenantPromotions (`tenant_promotions`) | ❌ |

### RunbookSnapshotResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Assembled (`assembled`) | ❌ |
| FrozenProjectVariableSetId (`frozen_project_variable_set_id`) | ❌ |
| FrozenRunbookProcessId (`frozen_runbook_process_id`) | ❌ |
| GitReference (`git_reference`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LibraryVariableSetSnapshotIds (`library_variable_set_snapshot_ids`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| Notes (`notes`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| ProjectVariableSetSnapshotId (`project_variable_set_snapshot_id`) | ❌ |
| RunbookId (`runbook_id`) | ❌ |
| SelectedGitResources (`selected_git_resources`) | ❌ |
| SelectedPackages (`selected_packages`) | ❌ |
| SpaceId (`space_id`) | ❌ |

### RunbookSnapshotTemplateResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| GitResources (`git_resources`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| NextNameIncrement (`next_name_increment`) | ❌ |
| Packages (`packages`) | ❌ |
| RunbookId (`runbook_id`) | ❌ |
| RunbookProcessId (`runbook_process_id`) | ❌ |

### RunbooksDashboardItemResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CompletedTime (`completed_time`) | ❌ |
| Created (`created`) | ❌ |
| Duration (`duration`) | ❌ |
| EnvironmentId (`environment_id`) | ❌ |
| ErrorMessage (`error_message`) | ❌ |
| GitReference (`git_reference`) | ❌ |
| HasPendingInterruptions (`has_pending_interruptions`) | ❌ |
| HasWarningsOrErrors (`has_warnings_or_errors`) | ❌ |
| Id (`id`) | ❌ |
| IsCompleted (`is_completed`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| QueueTime (`queue_time`) | ❌ |
| RunBy (`run_by`) | ❌ |
| RunName (`run_name`) | ❌ |
| RunbookId (`runbook_id`) | ❌ |
| RunbookSnapshotId (`runbook_snapshot_id`) | ❌ |
| RunbookSnapshotName (`runbook_snapshot_name`) | ❌ |
| RunbookSnapshotNotes (`runbook_snapshot_notes`) | ❌ |
| StartTime (`start_time`) | ❌ |
| State (`state`) | ❌ |
| TaskId (`task_id`) | ❌ |
| TenantId (`tenant_id`) | ❌ |

### RunbooksProgressionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Environments (`environments`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| RunbookRuns (`runbook_runs`) | ❌ |

### ScheduledTaskDetailsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActivityLog (`activity_log`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### ScheduledTaskStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsEnabled (`is_enabled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### SchedulerStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsRunning (`is_running`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| TaskStatus (`task_status`) | ❌ |

### ScopedUserRoleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EnvironmentIds (`environment_ids`) | ✅ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ProjectGroupIds (`project_group_ids`) | ✅ |
| ProjectIds (`project_ids`) | ✅ |
| SpaceId (`space_id`) | ✅ |
| TeamId (`team_id`) | ✅ |
| TenantIds (`tenant_ids`) | ✅ |
| UserRoleId (`user_role_id`) | ✅ |

### SelectedGitResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionName (`action_name`) | ❌ |
| GitReferenceResource (`git_reference_resource`) | ❌ |
| GitResourceReferenceName (`git_resource_reference_name`) | ❌ |

### ServerConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ServerUri (`server_uri`) | ❌ |

### ServerConfigurationSettingsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ConfigurationSet (`configuration_set`) | ❌ |
| ConfigurationValues (`configuration_values`) | ❌ |

### ServerConfigurationValueResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Key (`key`) | ❌ |
| Value (`value`) | ❌ |

### ServerStatusHealthResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| IsCompliantWithLicense (`is_compliant_with_license`) | ❌ |
| IsEntireClusterDrainingTasks (`is_entire_cluster_draining_tasks`) | ❌ |
| IsEntireClusterReadOnly (`is_entire_cluster_read_only`) | ❌ |
| IsOperatingNormally (`is_operating_normally`) | — |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### ServerStatusResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsDatabaseEncrypted (`is_database_encrypted`) | ❌ |
| IsInMaintenanceMode (`is_in_maintenance_mode`) | ❌ |
| IsMajorMinorUpgrade (`is_major_minor_upgrade`) | ❌ |
| IsUpgradeAvailable (`is_upgrade_available`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MaintenanceExpires (`maintenance_expires`) | ❌ |
| MaximumAvailableVersion (`maximum_available_version`) | ❌ |
| MaximumAvailableVersionCoveredByLicense (`maximum_available_version_covered_by_license`) | ❌ |

### ServerTaskStatusMessageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Category (`category`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Message (`message`) | ❌ |
| Title (`title`) | ❌ |

### ServerTimezoneResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsLocal (`is_local`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### ServiceAccountOidcIdentityResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Audience (`audience`) | ❌ |
| Id (`id`) | ✅ |
| Issuer (`issuer`) | ✅ |
| Name (`name`) | ✅ |
| ServiceAccountId (`service_account_id`) | ✅ |
| Subject (`subject`) | ✅ |

### SigningKeyConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ExpireAfterDays (`expire_after_days`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| PublicKeyHostingLocation (`public_key_hosting_location`) | ❌ |
| RevokeAfterDays (`revoke_after_days`) | ❌ |

### SmtpConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Details (`details`) | ❌ |
| EnableSsl (`enable_ssl`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| SendEmailFrom (`send_email_from`) | ❌ |
| SmtpHost (`smtp_host`) | ❌ |
| SmtpPort (`smtp_port`) | ❌ |
| Timeout (`timeout`) | ❌ |

### SmtpCredentialDetailsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CredentialType (`credential_type`) | ❌ |

### SmtpIsConfiguredResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsConfigured (`is_configured`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### SnapshotGitReferenceResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| GitCommit (`git_commit`) | ❌ |
| GitRef (`git_ref`) | ❌ |
| VariablesGitCommit (`variables_git_commit`) | ❌ |

### SpaceLicenseUsageResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| MachinesCount (`machines_count`) | ❌ |
| ProjectsCount (`projects_count`) | ❌ |
| SpaceName (`space_name`) | ❌ |
| TenantsCount (`tenants_count`) | ❌ |

### SpaceResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ✅ |
| ExtensionSettings (`extension_settings`) | — |
| Icon (`icon`) | ❌ |
| Id (`id`) | ❌ |
| IsDefault (`is_default`) | ✅ |
| IsPrivate (`is_private`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| Slug (`slug`) | ✅ |
| SpaceManagersTeamMembers (`space_managers_team_members`) | ✅ |
| SpaceManagersTeams (`space_managers_teams`) | ✅ |
| TaskQueueStopped (`task_queue_stopped`) | ❌ |

### SpaceRootResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### StepPackageInputsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Value (`value`) | ❌ |

### SubscriptionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| EventNotificationSubscription (`event_notification_subscription`) | ❌ |
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Type (`type`) | ❌ |

### SystemInfoResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ClrVersion (`clr_version`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MinThreadPoolCount (`min_thread_pool_count`) | ❌ |
| OSVersion (`o_s_version`) | ❌ |
| ThreadCount (`thread_count`) | ❌ |
| Uptime (`uptime`) | ❌ |
| Version (`version`) | ❌ |
| WorkingSetBytes (`working_set_bytes`) | ❌ |

### TagSetResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Description (`description`) | ✅ |
| Id (`id`) | ✅ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| Scopes (`scopes`) | ✅ |
| SortOrder (`sort_order`) | ✅ |
| SpaceId (`space_id`) | ✅ |
| Tags (`tags`) | ❌ |
| Type (`type`) | ✅ |

### TaskDetailsResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActivityLogs (`activity_logs`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| PhysicalLogSize (`physical_log_size`) | ❌ |
| Progress (`progress`) | ❌ |
| Task (`task`) | ❌ |

### TaskResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Arguments (`arguments`) | ❌ |
| CanRerun (`can_rerun`) | ❌ |
| Completed (`completed`) | ❌ |
| CompletedTime (`completed_time`) | ❌ |
| Description (`description`) | ❌ |
| Duration (`duration`) | ❌ |
| ErrorMessage (`error_message`) | ❌ |
| EstimatedRemainingQueueDurationSeconds (`estimated_remaining_queue_duration_seconds`) | ❌ |
| FinishedSuccessfully (`finished_successfully`) | — |
| HasBeenPickedUpByProcessor (`has_been_picked_up_by_processor`) | ❌ |
| HasPendingInterruptions (`has_pending_interruptions`) | ❌ |
| HasWarningsOrErrors (`has_warnings_or_errors`) | ❌ |
| Id (`id`) | ❌ |
| IsCompleted (`is_completed`) | — |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LastUpdatedTime (`last_updated_time`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| ProjectId (`project_id`) | ❌ |
| QueueTime (`queue_time`) | ❌ |
| QueueTimeExpiry (`queue_time_expiry`) | ❌ |
| ServerNode (`server_node`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| StartTime (`start_time`) | ❌ |
| State (`state`) | ❌ |
| UnmetPreconditions (`unmet_preconditions`) | ❌ |

### TaskTypeResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### TeamResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanBeDeleted (`can_be_deleted`) | ✅ |
| CanBeRenamed (`can_be_renamed`) | ✅ |
| CanChangeMembers (`can_change_members`) | ✅ |
| CanChangeRoles (`can_change_roles`) | ✅ |
| Description (`description`) | ✅ |
| ExternalSecurityGroups (`external_security_groups`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MemberUserIds (`member_user_ids`) | ❌ |
| Name (`name`) | ✅ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |

### TelemetryConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Enabled (`enabled`) | ❌ |
| Id (`id`) | ❌ |
| IsTelemetryEnforced (`is_telemetry_enforced`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| ShowAsNewUntil (`show_as_new_until`) | ❌ |

### TemplateParameterValueResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Scope (`scope`) | ❌ |
| Value (`value`) | ❌ |

### TenantResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ClonedFromTenantId (`cloned_from_tenant_id`) | ✅ |
| CustomFields (`custom_fields`) | ❌ |
| Description (`description`) | ✅ |
| Icon (`icon`) | ❌ |
| Id (`id`) | ✅ |
| IsDisabled (`is_disabled`) | ✅ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ✅ |
| ProjectEnvironments (`project_environments`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |
| TenantTags (`tenant_tags`) | ✅ |

### TenantVariableResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ConcurrencyToken (`concurrency_token`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| LibraryVariables (`library_variables`) | ❌ |
| Links (`links`) | ❌ |
| ProjectVariables (`project_variables`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| TenantId (`tenant_id`) | ❌ |
| TenantName (`tenant_name`) | ❌ |

### TenantsMissingVariablesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Links (`links`) | — |
| MissingVariables (`missing_variables`) | ❌ |
| TenantId (`tenant_id`) | ❌ |

### TriggerActionResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| ActionType (`action_type`) | — |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### TriggerFilterResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| FilterType (`filter_type`) | — |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |

### UpgradeConfigurationResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| AllowChecking (`allow_checking`) | ❌ |
| Id (`id`) | ❌ |
| IncludeStatistics (`include_statistics`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| NotificationMode (`notification_mode`) | ❌ |

### UserPermissionSetResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| IsPermissionsComplete (`is_permissions_complete`) | ❌ |
| IsTeamsComplete (`is_teams_complete`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| SpacePermissions (`space_permissions`) | ❌ |
| SystemPermissions (`system_permissions`) | ❌ |
| Teams (`teams`) | ❌ |

### UserResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanPasswordBeEdited (`can_password_be_edited`) | ✅ |
| Created (`created`) | ❌ |
| DisplayName (`display_name`) | ✅ |
| EmailAddress (`email_address`) | ✅ |
| Id (`id`) | ❌ |
| Identities (`identities`) | ❌ |
| IsActive (`is_active`) | ✅ |
| IsRequestor (`is_requestor`) | ✅ |
| IsService (`is_service`) | ✅ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Password (`password`) | ✅ |
| Username (`username`) | ✅ |

### UserRoleResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanBeDeleted (`can_be_deleted`) | ❌ |
| Description (`description`) | ❌ |
| GrantedSpacePermissions (`granted_space_permissions`) | ❌ |
| GrantedSystemPermissions (`granted_system_permissions`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| SpacePermissionDescriptions (`space_permission_descriptions`) | ❌ |
| SupportedRestrictions (`supported_restrictions`) | ❌ |
| SystemPermissionDescriptions (`system_permission_descriptions`) | ❌ |

### ValidatedGitReferenceResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanonicalName (`canonical_name`) | ❌ |
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |

### VariableSetResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| OwnerId (`owner_id`) | ❌ |
| ScopeValues (`scope_values`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| Variables (`variables`) | ❌ |
| Version (`version`) | ❌ |

### VariablesScopedToDocumentResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| HasUnauthorizedLibraryVariableSetVariables (`has_unauthorized_library_variable_set_variables`) | ❌ |
| HasUnauthorizedProjectVariables (`has_unauthorized_project_variables`) | ❌ |
| VariableMap (`variable_map`) | ❌ |

### VersioningStrategyResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| DonorPackage (`donor_package`) | ❌ |
| Template (`template`) | ❌ |

### WorkerPoolDynamicWorkerTypesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | — |
| Links (`links`) | ❌ |
| WorkerTypes (`worker_types`) | ❌ |

### WorkerPoolResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| CanAddWorkers (`can_add_workers`) | ❌ |
| Description (`description`) | ❌ |
| Id (`id`) | ❌ |
| IsDefault (`is_default`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| Name (`name`) | ❌ |
| Slug (`slug`) | ❌ |
| SortOrder (`sort_order`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| WorkerPoolType (`worker_pool_type`) | ❌ |

### WorkerPoolSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| MachineEndpointSummaries (`machine_endpoint_summaries`) | ❌ |
| MachineHealthStatusSummaries (`machine_health_status_summaries`) | ❌ |
| MachineIdsForCalamariUpgrade (`machine_ids_for_calamari_upgrade`) | ❌ |
| MachineIdsForTentacleUpgrade (`machine_ids_for_tentacle_upgrade`) | ❌ |
| TentacleUpgradesRequired (`tentacle_upgrades_required`) | ❌ |
| TotalDisabledMachines (`total_disabled_machines`) | ❌ |
| TotalMachines (`total_machines`) | ❌ |
| WorkerPool (`worker_pool`) | ❌ |

### WorkerPoolSupportedTypesResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Id (`id`) | — |
| Links (`links`) | ❌ |
| SupportedPoolTypes (`supported_pool_types`) | ❌ |

### WorkerPoolsSummaryResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| MachineEndpointSummaries (`machine_endpoint_summaries`) | ❌ |
| MachineHealthStatusSummaries (`machine_health_status_summaries`) | ❌ |
| MachineIdsForCalamariUpgrade (`machine_ids_for_calamari_upgrade`) | ❌ |
| MachineIdsForTentacleUpgrade (`machine_ids_for_tentacle_upgrade`) | ❌ |
| TentacleUpgradesRequired (`tentacle_upgrades_required`) | ❌ |
| TotalDisabledMachines (`total_disabled_machines`) | ❌ |
| TotalMachines (`total_machines`) | ❌ |
| WorkerPoolSummaries (`worker_pool_summaries`) | ❌ |

### WorkerResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Architecture (`architecture`) | ❌ |
| Endpoint (`endpoint`) | ❌ |
| HasLatestCalamari (`has_latest_calamari`) | ❌ |
| HealthStatus (`health_status`) | ❌ |
| Id (`id`) | ❌ |
| IsDisabled (`is_disabled`) | ✅ |
| IsInProcess (`is_in_process`) | ❌ |
| LastModifiedBy (`last_modified_by`) | ❌ |
| LastModifiedOn (`last_modified_on`) | ❌ |
| Links (`links`) | ❌ |
| MachinePolicyId (`machine_policy_id`) | ✅ |
| Name (`name`) | ✅ |
| OperatingSystem (`operating_system`) | ❌ |
| ShellName (`shell_name`) | ❌ |
| ShellVersion (`shell_version`) | ❌ |
| SkipInitialHealthCheck (`skip_initial_health_check`) | ❌ |
| Slug (`slug`) | ❌ |
| SpaceId (`space_id`) | ✅ |
| StatusSummary (`status_summary`) | ❌ |
| Thumbprint (`thumbprint`) | ✅ |
| Uri (`uri`) | ✅ |
| WorkerPoolIds (`worker_pool_ids`) | ✅ |

### WorkerTaskLeaseResource

| API property (snake_case) | In TF? |
|--------------------------|--------|
| Exclusive (`exclusive`) | ❌ |
| Id (`id`) | ❌ |
| Name (`name`) | ❌ |
| ServerTaskId (`server_task_id`) | ❌ |
| SpaceId (`space_id`) | ❌ |
| TakenAt (`taken_at`) | ❌ |
| WorkerId (`worker_id`) | ❌ |
| WorkerPoolId (`worker_pool_id`) | ❌ |

## 5. Terraform resources (Framework + SDK)

| Source | Resource name |
|--------|---------------|
| Framework | `octopusdeploy_certificate` |
| Framework | `octopusdeploy_channel` |
| Framework | `octopusdeploy_space` |
| Framework | `octopusdeploy_project_group` |
| Framework | `octopusdeploy_maven_feed` |
| Framework | `octopusdeploy_o_c_i_registry_feed` |
| Framework | `octopusdeploy_s3_feed` |
| Framework | `octopusdeploy_google_container_registry_feed` |
| Framework | `octopusdeploy_azure_container_registry_feed` |
| Framework | `octopusdeploy_amazon_web_services_account` |
| Framework | `octopusdeploy_azure_subscription_account` |
| Framework | `octopusdeploy_lifecycle` |
| Framework | `octopusdeploy_environment` |
| Framework | `octopusdeploy_parent_environment` |
| Framework | `octopusdeploy_step_template` |
| Framework | `octopusdeploy_community_step_template` |
| Framework | `octopusdeploy_git_credential` |
| Framework | `octopusdeploy_platform_hub_git_credential` |
| Framework | `octopusdeploy_platform_hub_aws_account` |
| Framework | `octopusdeploy_platform_hub_aws_open_i_d_connect_account` |
| Framework | `octopusdeploy_platform_hub_azure_oidc_account` |
| Framework | `octopusdeploy_platform_hub_azure_service_principal_account` |
| Framework | `octopusdeploy_platform_hub_gcp_account` |
| Framework | `octopusdeploy_platform_hub_generic_oidc_account` |
| Framework | `octopusdeploy_platform_hub_username_password_account` |
| Framework | `octopusdeploy_helm_feed` |
| Framework | `octopusdeploy_artifactory_generic_feed` |
| Framework | `octopusdeploy_git_hub_repository_feed` |
| Framework | `octopusdeploy_aws_elastic_container_registry_feed` |
| Framework | `octopusdeploy_nuget_feed` |
| Framework | `octopusdeploy_npm_feed` |
| Framework | `octopusdeploy_tenant_project` |
| Framework | `octopusdeploy_tenant_project_variable` |
| Framework | `octopusdeploy_tenant_common_variable` |
| Framework | `octopusdeploy_library_variable_set_feed` |
| Framework | `octopusdeploy_variable` |
| Framework | `octopusdeploy_project` |
| Framework | `octopusdeploy_project_versioning_strategy` |
| Framework | `octopusdeploy_machine_proxy` |
| Framework | `octopusdeploy_tag` |
| Framework | `octopusdeploy_docker_container_registry_feed` |
| Framework | `octopusdeploy_tag_set` |
| Framework | `octopusdeploy_username_password_account` |
| Framework | `octopusdeploy_runbook` |
| Framework | `octopusdeploy_tenant` |
| Framework | `octopusdeploy_tentacle_certificate` |
| Framework | `octopusdeploy_listening_tentacle_worker` |
| Framework | `octopusdeploy_s_s_h_connection_worker` |
| Framework | `octopusdeploy_script_module` |
| Framework | `octopusdeploy_user` |
| Framework | `octopusdeploy_deployment_freeze` |
| Framework | `octopusdeploy_deployment_freeze_project` |
| Framework | `octopusdeploy_generic_oidc` |
| Framework | `octopusdeploy_deployment_freeze_tenant` |
| Framework | `octopusdeploy_git_trigger` |
| Framework | `octopusdeploy_built_in_trigger` |
| Framework | `octopusdeploy_process` |
| Framework | `octopusdeploy_process_step` |
| Framework | `octopusdeploy_process_steps_order` |
| Framework | `octopusdeploy_process_child_step` |
| Framework | `octopusdeploy_process_child_steps_order` |
| Framework | `octopusdeploy_process_templated_step` |
| Framework | `octopusdeploy_process_templated_child_step` |
| Framework | `octopusdeploy_project_deployment_freeze` |
| Framework | `octopusdeploy_project_auto_create_release` |
| Framework | `octopusdeploy_kubernetes_monitor` |
| Framework | `octopusdeploy_team` |
| Framework | `octopusdeploy_scoped_user_role` |
| Framework | `octopusdeploy_space_default_lifecycle_release_retention_policy` |
| Framework | `octopusdeploy_space_default_lifecycle_tentacle_retention_policy` |
| Framework | `octopusdeploy_space_default_runbook_retention_policy` |
| Framework | `octopusdeploy_platform_hub_version_control_username_password_settings` |
| Framework | `octopusdeploy_platform_hub_version_control_anonymous_settings` |
| SDK | `octopusdeploy_aws_openid_connect_account` |
| SDK | `octopusdeploy_azure_cloud_service_deployment_target` |
| SDK | `octopusdeploy_azure_service_fabric_cluster_deployment_target` |
| SDK | `octopusdeploy_azure_service_principal` |
| SDK | `octopusdeploy_azure_openid_connect` |
| SDK | `octopusdeploy_azure_web_app_deployment_target` |
| SDK | `octopusdeploy_cloud_region_deployment_target` |
| SDK | `octopusdeploy_deployment_process` |
| SDK | `octopusdeploy_dynamic_worker_pool` |
| SDK | `octopusdeploy_gcp_account` |
| SDK | `octopusdeploy_kubernetes_agent_deployment_target` |
| SDK | `octopusdeploy_kubernetes_agent_worker` |
| SDK | `octopusdeploy_kubernetes_cluster_deployment_target` |
| SDK | `octopusdeploy_listening_tentacle_deployment_target` |
| SDK | `octopusdeploy_machine_policy` |
| SDK | `octopusdeploy_offline_package_drop_deployment_target` |
| SDK | `octopusdeploy_polling_tentacle_deployment_target` |
| SDK | `octopusdeploy_polling_subscription_id` |
| SDK | `octopusdeploy_project_deployment_target_trigger` |
| SDK | `octopusdeploy_external_feed_create_release_trigger` |
| SDK | `octopusdeploy_project_scheduled_trigger` |
| SDK | `octopusdeploy_runbook_process` |
| SDK | `octopusdeploy_ssh_connection_deployment_target` |
| SDK | `octopusdeploy_ssh_key_account` |
| SDK | `octopusdeploy_static_worker_pool` |
| SDK | `octopusdeploy_token_account` |
| SDK | `octopusdeploy_user_role` |