# Terraform Provider vs Octopus Server API (Source Code)

Comparison against the **Octopus Server C# source** (controllers + Resource/Command/Request contracts).

---

## 1. Resources in the Server/API not in Terraform

Server `*Resource` contracts that have **no corresponding Terraform resource** (no mapping or no schema/provider registration).

- **AccountDetailsResource**
- **AccountResource**
- **AccountUsageResource**
- **ActionTemplateCategoryResource**
- **ActionTemplateParameterResource**
- **ActionTemplateSearchResource**
- **ActionTemplateUsageResource**
- **ActionUpdateResultResource**
- **ActionsUpdateProcessResource**
- **AggregatedPerformanceTelemetryResource**
- **AnnuallyRecurringScheduleResource**
- **AnonymousGitCredentialUsageResource**
- **ApiKeyCreatedResource**
- **ApiKeyResource**
- **ArcFeedFilterResource**
- **ArchiveLimitConfigurationResource**
- **ArchivedEventFileResource**
- **ArtifactResource**
- **AuditStreamConfigurationResource**
- **AuthenticationConfigResource**
- **AuthenticationResource**
- **AutoDeployActionResource**
- **AutoDeployReleaseOverrideResource**
- **AutomaticDeprovisioningRuleResource**
- **AwsElasticContainerRegistryFeedResource**
- **AzureContainerRegistryFeedResource**
- **AzureEnvironmentResource**
- **AzureResourceGroupResource**
- **AzureRootResource**
- **AzureSmtpConfigurationResource**
- **AzureStorageAccountResource**
- **AzureWebSiteResource**
- **AzureWebSiteSlotResource**
- **BackupConfigurationResource**
- **BaseEnvironmentV2Resource**
- **BaseProcessTemplateUsageParameterValueResource**
- **BehaviouralTelemetryConfigurationResource**
- **BuiltInFeedResource**
- **BuiltInFeedStatsResource**
- **CertificateConfigurationResource**
- **CertificateUsageResource**
- **ChangeResource**
- **ChannelCustomFieldDefinitionResource**
- **ChannelGitResourceRuleResource**
- **ChannelVersionRuleResource**
- **CloudRegionEndpointResource**
- **CloudTemplateResource**
- **CommunityActionTemplateResource**
- **CommunityActionTemplateSnapshotResource**
- **CompliancePolicyBffResource**
- **CompliancePolicyResource**
- **CompliancePolicySummaryResource**
- **CompliancePolicyVersionResource**
- **ConnectProjectToTenantsTaskResource**
- **ContinueFirstProjectForOnboardingResource**
- **ContinuousDailyScheduledTriggerFilterResource**
- **CreateEndpointFromEndpointResource**
- **CreateReleaseActionResource**
- **CronScheduledTriggerFilterResource**
- **DailyRecurringScheduleResource**
- **DailyScheduledTriggerFilterResource**
- **DashboardConfigurationResource**
- **DashboardEnvironmentResource**
- **DashboardItemResource**
- **DashboardProjectGroupResource**
- **DashboardProjectResource**
- **DashboardResource**
- **DashboardTenantResource**
- **DatabaseMaintenanceResource**
- **DatabasePersistenceSettingsResource**
- **DaysPerMonthScheduledTriggerFilterResource**
- **DaysPerWeekScheduledTriggerFilterResource**
- **DefectResource**
- **DeployLatestReleaseActionResource**
- **DeployLatestReleaseToEnvironmentActionResource**
- **DeployNewReleaseActionResource**
- **DeploymentActionContainerResource**
- **DeploymentActionGitDependencyResource**
- **DeploymentActionPackageResource**
- **DeploymentActionResource**
- **DeploymentFreezeResource**
- **DeploymentPreviewBaseResource**
- **DeploymentPreviewResource**
- **DeploymentProcessResource**
- **DeploymentProcessV2Resource**
- **DeploymentResource**
- **DeploymentSettingsResource**
- **DeploymentStepResource**
- **DeploymentTemplateBaseResource**
- **DeploymentTemplateResource**
- **DeploymentVariablesResource**
- **DeprecationsConfigurationResource**
- **DetailedProjectsDashboardBffEnvironmentDeploymentsResource**
- **DetailedProjectsDashboardBffProjectGroupEnvironmentResource**
- **DetailedProjectsDashboardBffProjectGroupResource**
- **DetailedProjectsDashboardBffProjectResource**
- **DockerFeedResource**
- **DocumentTypeResource**
- **DynamicEnvironmentResource**
- **DynamicWorkerPoolResource**
- **EndpointResource**
- **EnvironmentSummaryResource**
- **EnvironmentSummaryV2Resource**
- **EnvironmentsSummaryResource**
- **EphemeralEnvironmentV2Resource**
- **EventAgentResource**
- **EventCategoryResource**
- **EventCsvResource**
- **EventGroupResource**
- **EventNotificationSubscriptionFilterResource**
- **EventNotificationSubscriptionResource**
- **EventResource**
- **EventRetentionConfigurationResource**
- **ExtensionsInfoResource**
- **FeatureToggleEnvironmentResource**
- **FeatureToggleResource**
- **FeaturesConfigurationResource**
- **FeedFilterResource**
- **FeedResource**
- **FreezeScheduleResource**
- **GetServerTaskApprovalByTaskIdResponseResource**
- **GitBranchResource**
- **GitCommitResource**
- **GitCredentialDetailsResource**
- **GitCredentialRepositoryRestrictionsResource**
- **GitCredentialResource**
- **GitCredentialUsageResource**
- **GitDependencyCollectionResource**
- **GitDependencyResource**
- **GitFilterResource**
- **GitHubAppConnectionResource**
- **GitHubAppInstallationResource**
- **GitHubConfigurationResource**
- **GitHubFeedResource**
- **GitHubGitCredentialUsageResource**
- **GitNamedRefByNameResource**
- **GitPersistenceSettingsConversionStateResource**
- **GitPersistenceSettingsResource**
- **GitReferenceResource**
- **GitTagResource**
- **GoogleContainerRegistryFeedResource**
- **GoogleSmtpConfigurationResource**
- **HostedServerStatusExternalHealthResource**
- **HostedServerStatusInternalHealthResource**
- **IconResource**
- **IdentityClaimResource**
- **IdentityResource**
- **InsightsDataSeriesResource**
- **InsightsEnvironmentGroupResource**
- **InsightsReportResource**
- **InterruptionResource**
- **InterruptionSummaryResource**
- **InvitationResource**
- **KubernetesAgentDetailsResource**
- **KubernetesResourceManifestDiffResource**
- **KubernetesTaskResourceStatusResource**
- **KubernetesTentacleEndpointResource**
- **LatestReleaseResource**
- **LetsEncryptConfigurationResource**
- **LibraryVariableSetUsageResource**
- **LibraryVariablesResource**
- **LicenseLimitStatusResource**
- **LicenseLimitUsageResource**
- **LicenseMessageResource**
- **LicenseResource**
- **LicenseStatusResource**
- **LicenseUsageResource**
- **LifecycleProgressionResource**
- **ListeningTentacleEndpointConfigurationResource**
- **ListeningTentacleEndpointResource**
- **LoadBalancerPingResource**
- **LoginInitiatedResource**
- **MachineBasedResource**
- **MachineCleanupPolicyResource**
- **MachineConnectivityPolicyResource**
- **MachineFilterResource**
- **MachineHealthCheckPolicyResource**
- **MachinePackageCacheRetentionPolicyResource**
- **MachinePolicyResource**
- **MachineResource**
- **MachineRpcCallRetryPolicyResource**
- **MachineScriptPolicyResource**
- **MachineUpdatePolicyResource**
- **MaintenanceConfigurationResource**
- **MigrationImportResource**
- **MigrationImportTaskArgumentsResource**
- **MigrationPartialExportResource**
- **MigrationPartialExportTaskArgumentsResource**
- **MissingVariableResource**
- **MonthlyRecurringScheduleResource**
- **MultiTenancyStatusResource**
- **NsbEndpointStatusResource**
- **NuGetFeedResource**
- **OctopusPackageVersionBuildInformationMappedResource**
- **OctopusPackageVersionBuildInformationResource**
- **OctopusProjectFeedResource**
- **OctopusServerClusterSummaryResource**
- **OctopusServerNodeDetailsResource**
- **OctopusServerNodeResource**
- **OctopusServerNodeSummaryResource**
- **OfflineDropDestinationResource**
- **OfflineDropEndpointResource**
- **OidcConfigResource**
- **OnboardingResource**
- **OnboardingTaskResource**
- **OnceDailyScheduledTriggerFilterResource**
- **OpenTelemetryAuditStreamConfigurationResource**
- **PackageAcquisitionLocationResource**
- **PackageDescriptionResource**
- **PackageFromBuiltInFeedResource**
- **PackageNoteListResource**
- **PackageReferenceCollectionResource**
- **PackageReferenceResource**
- **PackageResource**
- **PackageSignatureResource**
- **PackageVersionResource**
- **ParentEnvironmentV2Resource**
- **PerformanceConfigurationResource**
- **PersistenceSettingsResource**
- **PhaseDeploymentResource**
- **PhaseProgressionResource**
- **PhaseResource**
- **PlatformHubAccountResource**
- **PlatformHubGitCredentialResource**
- **PlatformHubVersionControlSettingsResource**
- **PollingTentacleEndpointConfigurationResource**
- **PollingTentacleEndpointResource**
- **ProcessElementResource**
- **ProcessTemplateParameterResource**
- **ProcessTemplateResource**
- **ProcessTemplateSummaryResource**
- **ProcessTemplateUsagePackageParameterValueResource**
- **ProcessTemplateUsageParameterValueResource**
- **ProcessTemplateUsageResource**
- **ProcessTemplateUsageSensitiveParameterValueResource**
- **ProcessTemplateVersionResource**
- **ProgressionResource**
- **ProjectIntentsResource**
- **ProjectTemplateDeploymentProcessResource**
- **ProjectTemplateParameterResource**
- **ProjectTemplateParameterSetResource**
- **ProjectTemplateParameterValueResource**
- **ProjectTemplateParametersResource**
- **ProjectTemplateVariableResource**
- **ProjectTemplateVersionResource**
- **ProjectTriggerExecutionDetailResource**
- **ProjectTriggerExecutionResource**
- **ProjectTriggerResource**
- **ProjectVariablesResource**
- **ProjectsPageDeploymentResource**
- **PropertyValueResource**
- **ProxyResource**
- **RecurringScheduleResource**
- **RedisStatusResource**
- **ReferenceGitCredentialUsageResource**
- **ReferencePlatformHubGitCredentialUsageResource**
- **ReleaseChangesResource**
- **ReleaseCreationStrategyResource**
- **ReleaseNoteOptionsResource**
- **ReleasePackageVersionBuildInformationResource**
- **ReleaseProgressionResource**
- **ReleaseResource**
- **ReleaseTemplateBaseResource**
- **ReleaseTemplateGitResource**
- **ReleaseTemplateResource**
- **ReportDeploymentCountOverTimeResource**
- **ResolvedProjectTemplateDetailsResource**
- **Resource**
- **RetentionDefaultConfigurationResource**
- **RetentionPoliciesConfigurationResource**
- **RetentionTypeResource**
- **RootResource**
- **RunRunbookActionResource**
- **RunbookProcessResource**
- **RunbookProcessV2Resource**
- **RunbookRunPreviewResource**
- **RunbookRunResource**
- **RunbookRunTemplateResource**
- **RunbookRunVariablesResource**
- **RunbookSnapshotProgressionResource**
- **RunbookSnapshotResource**
- **RunbookSnapshotTemplateResource**
- **RunbooksDashboardItemResource**
- **RunbooksProgressionResource**
- **ScheduledTaskDetailsResource**
- **ScheduledTaskStatusResource**
- **ScheduledTriggerFilterResource**
- **SchedulerStatusResource**
- **ScopedDeploymentActionResource**
- **ScopedUserRoleResource**
- **SelectedGitResource**
- **ServerActivitiesResource**
- **ServerActivityResource**
- **ServerConfigurationResource**
- **ServerConfigurationSettingsResource**
- **ServerConfigurationValueResource**
- **ServerStatusHealthResource**
- **ServerStatusResource**
- **ServerTaskApprovalResource**
- **ServerTaskStatusMessageResource**
- **ServerTimezoneResource**
- **SigningKeyConfigurationResource**
- **SigningKeyResource**
- **SmtpConfigurationResource**
- **SmtpCredentialDetailsResource**
- **SmtpIsConfiguredResource**
- **SnapshotGitReferenceResource**
- **SpaceLicenseUsageResource**
- **SpaceRootResource**
- **SplunkAuditStreamConfigurationResource**
- **SshEndpointResource**
- **SshKeyPairAccountResource**
- **StaticEnvironmentV2Resource**
- **StaticWorkerPoolResource**
- **StatusResource**
- **StepPackageDeploymentTargetTypeResource**
- **StepPackageEndpointResource**
- **StepPackageInputsResource**
- **SubscriptionResource**
- **SumoLogicAuditStreamConfigurationResource**
- **SystemInfoResource**
- **TaskDetailsResource**
- **TaskResource**
- **TaskSummaryResource**
- **TaskTypeResource**
- **TeamNameResource**
- **TelemetryConfigurationResource**
- **TemplatedProjectParameterValuesResource**
- **TenantProjectEnvironmentMappingResource**
- **TenantVariableResource**
- **TenantsMissingVariablesResource**
- **TentacleCommunicationModeResource**
- **TentacleDetailsResource**
- **TentacleEndpointConfigurationResource**
- **TentacleEndpointResource**
- **TokenAccountResource**
- **TriggerActionResource**
- **TriggerFilterResource**
- **UpgradeConfigurationResource**
- **UserOnboardingResource**
- **UserPermissionSetResource**
- **UserRoleResource**
- **UsernamePasswordGitCredentialDetailsResource**
- **UsernamePasswordGitCredentialUsageResource**
- **UsernamePasswordSmtpConfigurationResource**
- **ValidatedGitReferenceResource**
- **VariableResource**
- **VariableSetResource**
- **VariablesScopedToDocumentResource**
- **VersioningStrategyResource**
- **WeeklyRecurringScheduleResource**
- **WorkerPoolDynamicWorkerTypesResource**
- **WorkerPoolResource**
- **WorkerPoolSummaryResource**
- **WorkerPoolSupportedTypesResource**
- **WorkerPoolsSummaryResource**
- **WorkerTaskLeaseResource**

---

## 2. Resources in Terraform not in the Server/API

Terraform resources that are **not the target of any Server contract** in the current mapping (Framework + SDK).

- **octopusdeploy_aws_elastic_container_registry**
- **octopusdeploy_aws_openid_connect_account**
- **octopusdeploy_azure_cloud_service_deployment_target**
- **octopusdeploy_azure_container_registry**
- **octopusdeploy_azure_openid_connect**
- **octopusdeploy_azure_service_fabric_cluster_deployment_target**
- **octopusdeploy_azure_service_principal**
- **octopusdeploy_azure_web_app_deployment_target**
- **octopusdeploy_built_in_trigger**
- **octopusdeploy_cloud_region_deployment_target**
- **octopusdeploy_community_step_template**
- **octopusdeploy_deployment_freeze_project**
- **octopusdeploy_deployment_freeze_tenant**
- **octopusdeploy_deployment_process**
- **octopusdeploy_docker_container_registry**
- **octopusdeploy_dynamic_worker_pool**
- **octopusdeploy_external_feed_create_release_trigger**
- **octopusdeploy_gcp_account**
- **octopusdeploy_git_trigger**
- **octopusdeploy_github_repository_feed**
- **octopusdeploy_google_container_registry**
- **octopusdeploy_kubernetes_agent_deployment_target**
- **octopusdeploy_kubernetes_agent_worker**
- **octopusdeploy_kubernetes_cluster_deployment_target**
- **octopusdeploy_listening_tentacle_deployment_target**
- **octopusdeploy_machine_policy**
- **octopusdeploy_nuget_feed**
- **octopusdeploy_offline_package_drop_deployment_target**
- **octopusdeploy_parent_environment**
- **octopusdeploy_platform_hub_version_control_anonymous_settings**
- **octopusdeploy_platform_hub_version_control_username_password_settings**
- **octopusdeploy_polling_subscription_id**
- **octopusdeploy_polling_tentacle_deployment_target**
- **octopusdeploy_process**
- **octopusdeploy_process_child_step**
- **octopusdeploy_process_child_steps_order**
- **octopusdeploy_process_step**
- **octopusdeploy_process_steps_order**
- **octopusdeploy_process_templated_child_step**
- **octopusdeploy_process_templated_step**
- **octopusdeploy_project_auto_create_release**
- **octopusdeploy_project_deployment_target_trigger**
- **octopusdeploy_project_scheduled_trigger**
- **octopusdeploy_project_versioning_strategy**
- **octopusdeploy_runbook_process**
- **octopusdeploy_script_module**
- **octopusdeploy_space_default_lifecycle_release_retention_policy**
- **octopusdeploy_space_default_lifecycle_tentacle_retention_policy**
- **octopusdeploy_space_default_runbook_retention_policy**
- **octopusdeploy_ssh_connection_deployment_target**
- **octopusdeploy_ssh_key_account**
- **octopusdeploy_static_worker_pool**
- **octopusdeploy_tenant_project**
- **octopusdeploy_tentacle_certificate**
- **octopusdeploy_token_account**
- **octopusdeploy_user_role**

---

## 3. Property differences (Server/API vs Terraform)

Resources that exist in both: for each, properties that **differ** are listed with whether the property is **missing from Terraform** (API has it, TF does not) or **missing from Server/API** (TF has it, API does not).

**Note:** `*_ids` and the corresponding resource name are treated as the same (e.g. `environment_ids` ↔ `environments`, `tenant_ids` ↔ `tenants`) and are not listed as differing. Properties ending in `_type` (e.g. `account_type`, `feed_type`) are type discriminators that determine the concrete resource in Terraform and are not listed as missing. **Read-only properties** (those that exist on the Resource but not on the corresponding `ModifyXxxCommand`) are excluded from "Missing from Terraform" so only writable API properties are compared.

### ActionTemplateResource ↔ **octopusdeploy_step_template**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `acquisition_location` |  | ✓ |
| `default_branch` |  | ✓ |
| `default_sensitive_value` |  | ✓ |
| `default_value` |  | ✓ |
| `display_settings` |  | ✓ |
| `feed_id` |  | ✓ |
| `file_path_filters` |  | ✓ |
| `git_credential_id` |  | ✓ |
| `help_text` |  | ✓ |
| `inputs` | ✓ |  |
| `label` |  | ✓ |
| `package_id` |  | ✓ |
| `repository_uri` |  | ✓ |
| `step_package_id` |  | ✓ |
| `step_package_version` | ✓ |  |
| `step_template` |  | ✓ |

### AmazonWebServicesAccountResource ↔ **octopusdeploy_aws_account**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `details` | ✓ |  |

### AmazonWebServicesOidcAccountResource ↔ **octopusdeploy_aws_account**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `access_key` |  | ✓ |
| `account_test_subject_keys` | ✓ |  |
| `custom_claims` | ✓ |  |
| `deployment_subject_keys` | ✓ |  |
| `details` | ✓ |  |
| `health_check_subject_keys` | ✓ |  |
| `role_arn` | ✓ |  |
| `secret_key` |  | ✓ |
| `session_duration` | ✓ |  |

### AzureOidcAccountResource ↔ **octopusdeploy_azure_subscription_account**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `account_test_subject_keys` | ✓ |  |
| `active_directory_endpoint_base_uri` | ✓ |  |
| `audience` | ✓ |  |
| `certificate` |  | ✓ |
| `certificate_thumbprint` |  | ✓ |
| `client_id` | ✓ |  |
| `custom_claims` | ✓ |  |
| `deployment_subject_keys` | ✓ |  |
| `details` | ✓ |  |
| `health_check_subject_keys` | ✓ |  |
| `management_endpoint` |  | ✓ |
| `resource_management_endpoint_base_uri` | ✓ |  |
| `storage_endpoint_suffix` |  | ✓ |
| `subscription_id` |  | ✓ |
| `subscription_number` | ✓ |  |
| `tenant_id` | ✓ |  |

### AzureServicePrincipalAccountResource ↔ **octopusdeploy_azure_subscription_account**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `active_directory_endpoint_base_uri` | ✓ |  |
| `certificate` |  | ✓ |
| `certificate_thumbprint` |  | ✓ |
| `client_id` | ✓ |  |
| `details` | ✓ |  |
| `management_endpoint` |  | ✓ |
| `password` | ✓ |  |
| `resource_management_endpoint_base_uri` | ✓ |  |
| `storage_endpoint_suffix` |  | ✓ |
| `subscription_id` |  | ✓ |
| `subscription_number` | ✓ |  |
| `tenant_id` | ✓ |  |

### ChannelResource ↔ **octopusdeploy_channel**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `action_package` |  | ✓ |
| `automatic_ephemeral_environment_deployments` | ✓ |  |
| `custom_field_definitions` | ✓ |  |
| `deployment_action` |  | ✓ |
| `git_reference_rules` | ✓ |  |
| `git_resource_rules` | ✓ |  |
| `package_reference` |  | ✓ |
| `rule` |  | ✓ |
| `rules` | ✓ |  |
| `tag` |  | ✓ |
| `version_range` |  | ✓ |

### EnvironmentResource ↔ **octopusdeploy_environment**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `environments` |  | ✓ |
| `ids` |  | ✓ |
| `is_enabled` |  | ✓ |
| `jira_extension_settings` |  | ✓ |
| `jira_service_management_extension_settings` |  | ✓ |
| `partial_name` |  | ✓ |
| `servicenow_extension_settings` |  | ✓ |
| `skip` |  | ✓ |
| `take` |  | ✓ |

### GcsStorageFeedResource ↔ **octopusdeploy_gcs_storage_feed**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `audience` |  | ✓ |
| `package_acquisition_location_options` | ✓ |  |
| `subject_keys` |  | ✓ |

### GenericOidcAccountResource ↔ **octopusdeploy_generic_oidc_account**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `account_test_subject_keys` | ✓ |  |
| `deployment_subject_keys` | ✓ |  |
| `details` | ✓ |  |
| `execution_subject_keys` |  | ✓ |
| `health_check_subject_keys` | ✓ |  |

### KubernetesMonitorResource ↔ **octopusdeploy_kubernetes_monitor**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `authentication_token` |  | ✓ |
| `certificate_thumbprint` |  | ✓ |
| `preserve_authentication_token` |  | ✓ |

### LibraryVariableSetResource ↔ **octopusdeploy_library_variable_set**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `help_text` |  | ✓ |
| `ids` |  | ✓ |
| `label` |  | ✓ |
| `library_variable_sets` |  | ✓ |
| `partial_name` |  | ✓ |
| `skip` |  | ✓ |
| `take` |  | ✓ |
| `template` |  | ✓ |
| `template_ids` |  | ✓ |
| `templates` | ✓ |  |
| `version` | ✓ |  |

### LifecycleResource ↔ **octopusdeploy_lifecycle**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `automatic_deployment_targets` |  | ✓ |
| `ids` |  | ✓ |
| `is_optional_phase` |  | ✓ |
| `is_priority_phase` |  | ✓ |
| `minimum_environments_before_promotion` |  | ✓ |
| `optional_deployment_targets` |  | ✓ |
| `partial_name` |  | ✓ |
| `phases` | ✓ |  |
| `quantity_to_keep` |  | ✓ |
| `release_retention_policy` | ✓ |  |
| `should_keep_forever` |  | ✓ |
| `skip` |  | ✓ |
| `strategy` |  | ✓ |
| `take` |  | ✓ |
| `tentacle_retention_policy` | ✓ |  |
| `unit` |  | ✓ |

### OciRegistryFeedResource ↔ **octopusdeploy_oci_registry_feed**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `package_acquisition_location_options` | ✓ |  |

### ProjectGroupResource ↔ **octopusdeploy_project_group**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `ids` |  | ✓ |
| `partial_name` |  | ✓ |
| `project_groups` |  | ✓ |
| `skip` |  | ✓ |
| `take` |  | ✓ |

### ProjectResource ↔ **octopusdeploy_project**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `allow_deployments_to_no_targets` |  | ✓ |
| `allow_ignore_channel_rules` | ✓ |  |
| `base_path` |  | ✓ |
| `channel_id` |  | ✓ |
| `combine_health_and_sync_status_in_dashboard_live_status` | ✓ |  |
| `connection_id` |  | ✓ |
| `connectivity_policy` |  | ✓ |
| `default_branch` |  | ✓ |
| `default_power_shell_edition` | ✓ |  |
| `default_value` |  | ✓ |
| `deployment_action` |  | ✓ |
| `display_settings` |  | ✓ |
| `donor_package` |  | ✓ |
| `donor_package_step_id` |  | ✓ |
| `environment_id` |  | ✓ |
| `exclude_unhealthy_targets` |  | ✓ |
| `execute_deployments_on_resilient_pipeline` | ✓ |  |
| `force_package_download` | ✓ |  |
| `git_anonymous_persistence_settings` |  | ✓ |
| `git_credential_id` |  | ✓ |
| `git_github_app_persistence_settings` |  | ✓ |
| `git_library_persistence_settings` |  | ✓ |
| `git_username_password_persistence_settings` |  | ✓ |
| `github_connection_id` |  | ✓ |
| `help_text` |  | ✓ |
| `ids` |  | ✓ |
| `included_library_variable_set_ids` | ✓ |  |
| `included_library_variable_sets` |  | ✓ |
| `is_clone` |  | ✓ |
| `is_discrete_channel_release` |  | ✓ |
| `is_enabled` |  | ✓ |
| `is_state_automatically_transitioned` |  | ✓ |
| `jira_service_management_extension_settings` |  | ✓ |
| `label` |  | ✓ |
| `package_reference` |  | ✓ |
| `partial_name` |  | ✓ |
| `password` |  | ✓ |
| `persistence_settings` | ✓ |  |
| `project_connectivity_policy` | ✓ |  |
| `project_template_details` | ✓ |  |
| `protected_branches` |  | ✓ |
| `release_creation_package` |  | ✓ |
| `release_creation_package_step_id` |  | ✓ |
| `release_id` |  | ✓ |
| `service_desk_project_name` |  | ✓ |
| `servicenow_extension_settings` |  | ✓ |
| `skip` |  | ✓ |
| `skip_machine_behavior` |  | ✓ |
| `standard_change_template_name` |  | ✓ |
| `take` |  | ✓ |
| `target_roles` |  | ✓ |
| `template` |  | ✓ |
| `templates` | ✓ |  |
| `tenant_id` |  | ✓ |
| `tenanted_deployment_mode` | ✓ |  |
| `tenanted_deployment_participation` |  | ✓ |
| `url` |  | ✓ |
| `username` |  | ✓ |

### RunbookResource ↔ **octopusdeploy_runbook**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `allow_deployments_to_no_targets` |  | ✓ |
| `exclude_unhealthy_targets` |  | ✓ |
| `fail_target_discovery` | ✓ |  |
| `retention_policy` |  | ✓ |
| `retention_policy_with_strategy` |  | ✓ |
| `run_retention_policy` | ✓ |  |
| `skip_machine_behaviour` |  | ✓ |
| `target_roles` |  | ✓ |

### S3FeedResource ↔ **octopusdeploy_s3_feed**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `package_acquisition_location_options` | ✓ |  |
| `password` |  | ✓ |
| `region` | ✓ |  |
| `username` |  | ✓ |

### ServiceAccountOidcIdentityResource ↔ **octopusdeploy_service_account_oidc_identity**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `count` | ✓ |  |
| `external_id` | ✓ |  |
| `oidc_identities` | ✓ |  |
| `server_url` | ✓ |  |

### SpaceResource ↔ **octopusdeploy_space**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `icon` | ✓ |  |
| `is_private` | ✓ |  |
| `is_task_queue_stopped` |  | ✓ |
| `task_queue_stopped` | ✓ |  |

### TagResource ↔ **octopusdeploy_tag**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `tag_set_id` |  | ✓ |
| `tag_set_space_id` |  | ✓ |

### TagSetResource ↔ **octopusdeploy_tag_set**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `canonical_tag_name` |  | ✓ |
| `color` |  | ✓ |
| `ids` |  | ✓ |
| `partial_name` |  | ✓ |
| `skip` |  | ✓ |
| `tag_sets` |  | ✓ |
| `take` |  | ✓ |

### TeamResource ↔ **octopusdeploy_team**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `display_id_and_name` |  | ✓ |
| `display_name` |  | ✓ |
| `environment_ids` |  | ✓ |
| `external_security_group` |  | ✓ |
| `external_security_groups` | ✓ |  |
| `member_user_ids` | ✓ |  |
| `project_group_ids` |  | ✓ |
| `project_ids` |  | ✓ |
| `team_id` |  | ✓ |
| `tenant_ids` |  | ✓ |
| `user_role` |  | ✓ |
| `user_role_id` |  | ✓ |
| `users` |  | ✓ |

### TenantResource ↔ **octopusdeploy_tenant**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `ids` |  | ✓ |
| `is_clone` |  | ✓ |
| `partial_name` |  | ✓ |
| `project_environments` | ✓ |  |
| `project_id` |  | ✓ |
| `skip` |  | ✓ |
| `tags` |  | ✓ |
| `take` |  | ✓ |
| `tenants` |  | ✓ |

### UserResource ↔ **octopusdeploy_user**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `claim` |  | ✓ |
| `external_id` |  | ✓ |
| `filter` |  | ✓ |
| `identities` | ✓ |  |
| `identity` |  | ✓ |
| `ids` |  | ✓ |
| `is_identifying_claim` |  | ✓ |
| `name` |  | ✓ |
| `provider` |  | ✓ |
| `skip` |  | ✓ |
| `take` |  | ✓ |
| `users` |  | ✓ |
| `value` |  | ✓ |

### UsernamePasswordAccountResource ↔ **octopusdeploy_username_password_account**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `details` | ✓ |  |

### WorkerResource ↔ **octopusdeploy_listening_tentacle_worker**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `endpoint` | ✓ |  |
| `proxy_id` |  | ✓ |

### WorkerResource ↔ **octopusdeploy_ssh_connection_worker**

| Property | Missing from Terraform | Missing from Server/API |
|----------|------------------------|-------------------------|
| `account_id` |  | ✓ |
| `dotnet_platform` |  | ✓ |
| `endpoint` | ✓ |  |
| `fingerprint` |  | ✓ |
| `host` |  | ✓ |
| `port` |  | ✓ |
| `proxy_id` |  | ✓ |

---

## 4. Server endpoints and TF coverage (reference)

| Method | Route | TF coverage |
|--------|-------|-------------|
| GET | `api/.well-known/acme-challenge/{token}` | — |
| GET | `api/.well-known/jwks` | — |
| GET | `api/.well-known/openid-configuration` | — |
| GET | `api/accounts` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| POST | `api/accounts` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| GET | `api/accounts/all` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| PUT | `api/accounts/{accountId}` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| DELETE | `api/accounts/{id}` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| GET | `api/accounts/{id}` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| GET | `api/accounts/{id}/pk` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| GET | `api/accounts/{id}/usages` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| DELETE | `api/accounts/{id}/v1` | octopusdeploy_username_password_account, octopusdeploy_certificate ... |
| DELETE | `api/actionTemplates/{id}` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| DELETE | `api/actionTemplates/{id}/v1` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| POST | `api/actiontemplates` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/all` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/categories` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/search` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/search/v2` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{idOrActionType}/logo/v2` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{id}` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| PUT | `api/actiontemplates/{id}` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| POST | `api/actiontemplates/{id}/actionsUpdate` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| POST | `api/actiontemplates/{id}/actionsUpdate/bulk` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{id}/logo` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| POST | `api/actiontemplates/{id}/logo` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| PUT | `api/actiontemplates/{id}/logo` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| POST | `api/actiontemplates/{id}/logo/v2` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| PUT | `api/actiontemplates/{id}/logo/v2` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{id}/usage` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{id}/v1` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{id}/versions` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{id}/versions/{version}` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/actiontemplates/{typeOrId}/versions/{version}/logo` | octopusdeploy_step_template, octopusdeploy_community_step_template |
| GET | `api/api/capabilities` | — |
| GET | `api/api/capabilities/{capability}` | — |
| GET | `api/api/diagnostics/throw` | — |
| GET | `api/artifacts` | — |
| POST | `api/artifacts` | — |
| DELETE | `api/artifacts/{id}` | — |
| GET | `api/artifacts/{id}` | — |
| PUT | `api/artifacts/{id}` | — |
| GET | `api/artifacts/{id}/content` | — |
| PUT | `api/artifacts/{id}/content` | — |
| GET | `api/audit-stream` | — |
| PUT | `api/audit-stream` | — |
| GET | `api/authentication` | — |
| POST | `api/authentication/checklogininitiated` | — |
| GET | `api/build-information` | — |
| POST | `api/build-information` | — |
| DELETE | `api/build-information/bulk` | — |
| DELETE | `api/build-information/{id}` | — |
| GET | `api/build-information/{id}` | — |
| GET | `api/certificates` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| GET | `api/certificates/all` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| GET | `api/certificates/certificate-global` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates/generate` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| DELETE | `api/certificates/{id}` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| GET | `api/certificates/{id}` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| PUT | `api/certificates/{id}` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates/{id}/archive` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates/{id}/archive/v1` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| GET | `api/certificates/{id}/export` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates/{id}/replace` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates/{id}/unarchive` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| POST | `api/certificates/{id}/unarchive/v1` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| GET | `api/certificates/{id}/usages` | octopusdeploy_certificate, octopusdeploy_tentacle_certificate |
| GET | `api/channels` | octopusdeploy_channel |
| POST | `api/channels` | octopusdeploy_channel |
| GET | `api/channels/all` | octopusdeploy_channel |
| GET | `api/channels/rule-test` | octopusdeploy_channel |
| POST | `api/channels/rule-test` | octopusdeploy_channel |
| GET | `api/channels/rule-test/v1` | octopusdeploy_channel |
| POST | `api/channels/rule-test/v1` | octopusdeploy_channel |
| DELETE | `api/channels/{id}` | octopusdeploy_channel |
| GET | `api/channels/{id}` | octopusdeploy_channel |
| PUT | `api/channels/{id}` | octopusdeploy_channel |
| GET | `api/channels/{id}/releases` | octopusdeploy_channel |
| POST | `api/cloudtemplate/{id}/metadata` | — |
| GET | `api/communityactiontemplates` | — |
| GET | `api/communityactiontemplates/{id}` | — |
| GET | `api/communityactiontemplates/{id}/actiontemplate` | — |
| GET | `api/communityactiontemplates/{id}/actiontemplate/{actiontemplatespaceId}` | — |
| POST | `api/communityactiontemplates/{id}/installation` | — |
| PUT | `api/communityactiontemplates/{id}/installation` | — |
| POST | `api/communityactiontemplates/{id}/installation/{actiontemplatespaceId}` | — |
| PUT | `api/communityactiontemplates/{id}/installation/{actiontemplatespaceId}` | — |
| GET | `api/communityactiontemplates/{id}/logo` | — |
| GET | `api/configuration` | — |
| GET | `api/configuration/certificates` | — |
| GET | `api/configuration/certificates/{id}` | — |
| GET | `api/configuration/certificates/{id}/public-cer` | — |
| GET | `api/configuration/expired-feature-toggles` | — |
| GET | `api/configuration/feature-toggles` | — |
| GET | `api/configuration/open-telemetry-trace-file-export` | — |
| PUT | `api/configuration/open-telemetry-trace-file-export` | — |
| GET | `api/configuration/retention-default` | — |
| PUT | `api/configuration/retention-default` | — |
| POST | `api/configuration/versioncontrol/clear-cache` | — |
| POST | `api/configuration/versioncontrol/clear-cache/v1` | — |
| GET | `api/configuration/{id}` | — |
| GET | `api/configuration/{id}/metadata` | — |
| GET | `api/configuration/{id}/values` | — |
| PUT | `api/configuration/{id}/values` | — |
| GET | `api/connectionagent/registrations` | — |
| POST | `api/connectionagent/registrations` | — |
| DELETE | `api/connectionagent/registrations/{agentIdValue}` | — |
| GET | `api/dashboard` | — |
| GET | `api/dashboard/dynamic` | — |
| GET | `api/dashboardconfiguration` | — |
| PUT | `api/dashboardconfiguration` | — |
| GET | `api/debug/nsbendpoints` | — |
| GET | `api/debug/nsbendpoints/{endpointName}` | — |
| POST | `api/debug/nsbendpoints/{endpointName}/disable` | — |
| POST | `api/debug/nsbendpoints/{endpointName}/enable` | — |
| GET | `api/deploymentTargetTags/{tag}` | — |
| GET | `api/deploymentfreezes` | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project ... |
| POST | `api/deploymentfreezes` | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project ... |
| DELETE | `api/deploymentfreezes/{id}` | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project ... |
| GET | `api/deploymentfreezes/{id}` | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project ... |
| PUT | `api/deploymentfreezes/{id}` | octopusdeploy_deployment_freeze, octopusdeploy_deployment_freeze_project ... |
| GET | `api/deploymentprocesses` | octopusdeploy_process |
| GET | `api/deploymentprocesses/{deploymentProcessId}/template` | octopusdeploy_process |
| GET | `api/deploymentprocesses/{id}` | octopusdeploy_process |
| PUT | `api/deploymentprocesses/{id}` | octopusdeploy_process |
| GET | `api/deployments` | — |
| POST | `api/deployments` | — |
| POST | `api/deployments/override` | — |
| POST | `api/deployments/v1` | — |
| DELETE | `api/deployments/{id}` | — |
| GET | `api/deployments/{id}` | — |
| GET | `api/deploymentsettings/{id}` | — |
| PUT | `api/deploymentsettings/{projectId}` | — |
| POST | `api/deploymenttargets/upgrade` | — |
| GET | `api/deploymenttargettags` | — |
| POST | `api/deploymenttargettags` | — |
| DELETE | `api/deploymenttargettags/{tag}` | — |
| POST | `api/deprecations/example` | — |
| GET | `api/deprecations/schedules` | — |
| POST | `api/deprecations/toggle` | — |
| POST | `api/deprecations/toggle/v1` | — |
| GET | `api/deprecationsconfiguration` | — |
| POST | `api/deprecationsconfiguration` | — |
| GET | `api/dynamic-extensions/features/metadata` | — |
| GET | `api/dynamic-extensions/features/values` | — |
| PUT | `api/dynamic-extensions/features/values` | — |
| GET | `api/dynamic-extensions/scripts` | — |
| GET | `api/environments` | octopusdeploy_environment, octopusdeploy_parent_environment |
| POST | `api/environments` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/all` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/all/v1` | octopusdeploy_environment, octopusdeploy_parent_environment |
| PUT | `api/environments/sortorder` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/summary` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/v1` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/v2` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{environmentId:regex(Environments-\\d+)}/metadata` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{environmentId:regex(Environments-\\d+)}/singlyScopedVariableDetails` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{environmentId:regex(Environments-d+)}/metadata` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{environmentId:regex(Environments-d+)}/singlyScopedVariableDetails` | octopusdeploy_environment, octopusdeploy_parent_environment |
| PUT | `api/environments/{environmentId}` | octopusdeploy_environment, octopusdeploy_parent_environment |
| DELETE | `api/environments/{id}` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{id}` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{id}/machines` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/environments/{id}/v2` | octopusdeploy_environment, octopusdeploy_parent_environment |
| GET | `api/errors/400` | — |
| GET | `api/errors/404-exception` | — |
| GET | `api/errors/divide-by-zero` | — |
| GET | `api/errors/no-permission` | — |
| GET | `api/errors/server-returns-html-instead-of-json-for-some-reason` | — |
| GET | `api/errors/unauthorized` | — |
| GET | `api/events` | — |
| GET | `api/events/agents` | — |
| GET | `api/events/archives` | — |
| GET | `api/events/archives/v1` | — |
| DELETE | `api/events/archives/{fileName}` | — |
| GET | `api/events/archives/{fileName}` | — |
| DELETE | `api/events/archives/{fileName}/v1` | — |
| GET | `api/events/categories` | — |
| GET | `api/events/documenttypes` | — |
| GET | `api/events/groups` | — |
| GET | `api/events/{id}` | — |
| GET | `api/externalsecuritygroupproviders` | — |
| GET | `api/featuresconfiguration` | — |
| PUT | `api/featuresconfiguration` | — |
| GET | `api/feeds` | — |
| POST | `api/feeds` | — |
| GET | `api/feeds/all` | — |
| GET | `api/feeds/stats` | — |
| GET | `api/feeds/{feedId}/packages` | — |
| GET | `api/feeds/{feedId}/packages/notes` | — |
| DELETE | `api/feeds/{id}` | — |
| GET | `api/feeds/{id}` | — |
| PUT | `api/feeds/{id}` | — |
| GET | `api/feeds/{id}/packages/search` | — |
| GET | `api/feeds/{id}/packages/versions` | — |
| GET | `api/github/accounts/install-url` | — |
| GET | `api/github/app/settings` | — |
| GET | `api/github/installations/updated` | — |
| POST | `api/github/reset-registration` | — |
| GET | `api/github/user/app/authorization_status` | — |
| GET | `api/github/user/app/authorize` | — |
| POST | `api/github/user/app/exchange-access-code` | — |
| DELETE | `api/github/user/app/token` | — |
| GET | `api/github/user/app/token` | — |
| POST | `api/github/user/app/token/refresh` | — |
| POST | `api/githubissuetracker/connectivitycheck` | — |
| GET | `api/iamreplacedwitharealroute` | — |
| GET | `api/icons/all` | — |
| GET | `api/icons/categories` | — |
| GET | `api/insights/reports/{reportId}/logo` | — |
| GET | `api/interruptions` | — |
| GET | `api/interruptions/{id}` | — |
| GET | `api/interruptions/{id}/responsible` | — |
| PUT | `api/interruptions/{id}/responsible` | — |
| POST | `api/interruptions/{id}/submit` | — |
| GET | `api/letsencryptconfiguration` | — |
| PUT | `api/letsencryptconfiguration` | — |
| GET | `api/libraryvariablesets` | octopusdeploy_library_variable_set |
| POST | `api/libraryvariablesets` | octopusdeploy_library_variable_set |
| GET | `api/libraryvariablesets/all` | octopusdeploy_library_variable_set |
| GET | `api/libraryvariablesets/all/v1` | octopusdeploy_library_variable_set |
| POST | `api/libraryvariablesets/all/v1` | octopusdeploy_library_variable_set |
| DELETE | `api/libraryvariablesets/{id}` | octopusdeploy_library_variable_set |
| GET | `api/libraryvariablesets/{id}` | octopusdeploy_library_variable_set |
| PUT | `api/libraryvariablesets/{id}` | octopusdeploy_library_variable_set |
| GET | `api/libraryvariablesets/{id}/usages` | octopusdeploy_library_variable_set |
| GET | `api/licenses/licenses-current` | — |
| PUT | `api/licenses/licenses-current` | — |
| GET | `api/licenses/licenses-current-features` | — |
| GET | `api/licenses/licenses-current-status` | — |
| GET | `api/licenses/licenses-current-usage` | — |
| GET | `api/lifecycles` | octopusdeploy_lifecycle |
| POST | `api/lifecycles` | octopusdeploy_lifecycle |
| GET | `api/lifecycles/all` | octopusdeploy_lifecycle |
| GET | `api/lifecycles/previews` | octopusdeploy_lifecycle |
| DELETE | `api/lifecycles/{id}` | octopusdeploy_lifecycle |
| GET | `api/lifecycles/{id}` | octopusdeploy_lifecycle |
| PUT | `api/lifecycles/{id}` | octopusdeploy_lifecycle |
| GET | `api/lifecycles/{id}/preview` | octopusdeploy_lifecycle |
| GET | `api/lifecycles/{id}/projects` | octopusdeploy_lifecycle |
| GET | `api/machinepolicies` | — |
| POST | `api/machinepolicies` | — |
| GET | `api/machinepolicies/all` | — |
| GET | `api/machinepolicies/template` | — |
| DELETE | `api/machinepolicies/{id}` | — |
| GET | `api/machinepolicies/{id}` | — |
| PUT | `api/machinepolicies/{id}` | — |
| GET | `api/machinepolicies/{id}/machines` | — |
| DELETE | `api/machinepolicies/{id}/v1` | — |
| GET | `api/machinepolicies/{id}/workers` | — |
| GET | `api/machineroles/all` | — |
| GET | `api/machineroles/all/v1` | — |
| GET | `api/machines` | — |
| POST | `api/machines` | — |
| GET | `api/machines/all` | — |
| GET | `api/machines/all/v1` | — |
| GET | `api/machines/discover` | — |
| GET | `api/machines/operatingsystem/names/all` | — |
| GET | `api/machines/operatingsystem/shells/all` | — |
| POST | `api/machines/upgrade` | — |
| DELETE | `api/machines/{id}` | — |
| GET | `api/machines/{id}` | — |
| GET | `api/machines/{id}/connection` | — |
| GET | `api/machines/{id}/latestdeployments` | — |
| GET | `api/machines/{id}/tasks` | — |
| GET | `api/machines/{id}/tasks/v1` | — |
| GET | `api/machines/{machineId:regex(Machines-\\d+)}/singlyScopedVariableDetails` | — |
| GET | `api/machines/{machineId:regex(Machines-d+)}/singlyScopedVariableDetails` | — |
| PUT | `api/machines/{machineid}` | — |
| GET | `api/maintenanceconfiguration` | — |
| PUT | `api/maintenanceconfiguration` | — |
| PUT | `api/maintenanceconfiguration/v1` | — |
| GET | `api/maintenanceoperations/databasemaintenance` | — |
| GET | `api/maintenanceoperations/databasemaintenance/{id}` | — |
| POST | `api/maintenanceoperations/databasemaintenance/{id}/cancel` | — |
| GET | `api/maintenanceoperations/isdatabasemaintenancerequired` | — |
| POST | `api/maintenanceoperations/startdatabasemaintenance` | — |
| POST | `api/migrations/import` | — |
| POST | `api/migrations/partialexport` | — |
| GET | `api/octopusservernodes` | — |
| GET | `api/octopusservernodes/all` | — |
| GET | `api/octopusservernodes/ping` | — |
| GET | `api/octopusservernodes/summary` | — |
| DELETE | `api/octopusservernodes/{id}` | — |
| GET | `api/octopusservernodes/{id}` | — |
| PUT | `api/octopusservernodes/{id}` | — |
| GET | `api/octopusservernodes/{id}/details` | — |
| GET | `api/packages` | — |
| DELETE | `api/packages/bulk` | — |
| DELETE | `api/packages/bulk/v1` | — |
| GET | `api/packages/notes` | — |
| POST | `api/packages/raw` | — |
| DELETE | `api/packages/{id}` | — |
| GET | `api/packages/{id}` | — |
| GET | `api/packages/{id}/raw` | — |
| DELETE | `api/packages/{id}/v1` | — |
| POST | `api/packages/{packageId}/{baseVersion}/delta` | — |
| GET | `api/packages/{packageId}/{version}/delta-signature` | — |
| GET | `api/performanceconfiguration` | — |
| PUT | `api/performanceconfiguration` | — |
| GET | `api/permissions/all` | — |
| GET | `api/platformhub/accounts` | — |
| POST | `api/platformhub/accounts` | — |
| DELETE | `api/platformhub/accounts/{id}` | — |
| GET | `api/platformhub/accounts/{id}` | — |
| PUT | `api/platformhub/accounts/{id}` | — |
| GET | `api/platformhub/git-credentials` | — |
| POST | `api/platformhub/git-credentials` | — |
| DELETE | `api/platformhub/git-credentials/{id}` | — |
| GET | `api/platformhub/git-credentials/{id}` | — |
| PUT | `api/platformhub/git-credentials/{id}` | — |
| GET | `api/platformhub/git/branches` | — |
| POST | `api/platformhub/git/branches` | — |
| GET | `api/platformhub/git/tags` | — |
| POST | `api/platformhub/github/connections` | — |
| GET | `api/platformhub/github/installations/{installationId}/repositories` | — |
| GET | `api/platformhub/policies/{slug}/versions` | — |
| POST | `api/platformhub/policies/{slug}/versions/{version}/modify-status` | — |
| GET | `api/platformhub/processtemplates/{slug}/share` | — |
| GET | `api/platformhub/processtemplates/{slug}/versions` | — |
| GET | `api/platformhub/processtemplates/{slug}/{versionMask}` | — |
| PUT | `api/platformhub/projecttemplates/{slug}` | — |
| GET | `api/platformhub/projecttemplates/{slug}/share` | — |
| GET | `api/platformhub/projecttemplates/{slug}/versions` | — |
| GET | `api/platformhub/projecttemplates/{templateSlug}/variables/sensitive` | — |
| PUT | `api/platformhub/projecttemplates/{templateSlug}/variables/sensitive` | — |
| GET | `api/platformhub/versioncontrol` | — |
| PUT | `api/platformhub/versioncontrol` | — |
| GET | `api/platformhub/{gitRef}/policies` | — |
| POST | `api/platformhub/{gitRef}/policies` | — |
| GET | `api/platformhub/{gitRef}/policies/{slug}` | — |
| PUT | `api/platformhub/{gitRef}/policies/{slug}` | — |
| POST | `api/platformhub/{gitRef}/policies/{slug}/publish` | — |
| GET | `api/platformhub/{gitRef}/processtemplates` | — |
| POST | `api/platformhub/{gitRef}/processtemplates` | — |
| GET | `api/platformhub/{gitRef}/processtemplates/summaries` | — |
| GET | `api/platformhub/{gitRef}/processtemplates/{slug}` | — |
| POST | `api/platformhub/{gitRef}/processtemplates/{slug}/share` | — |
| GET | `api/platformhub/{gitRef}/processtemplates/{slug}/variables/names` | — |
| POST | `api/platformhub/{gitRef}/processtemplates/{slug}/versions` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates` | — |
| POST | `api/platformhub/{gitRef}/projecttemplates` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates/{slug}` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates/{slug}/deploymentprocess` | — |
| PUT | `api/platformhub/{gitRef}/projecttemplates/{slug}/deploymentprocess` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates/{slug}/link` | — |
| POST | `api/platformhub/{gitRef}/projecttemplates/{slug}/publish` | — |
| POST | `api/platformhub/{gitRef}/projecttemplates/{slug}/share` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates/{slug}/status` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates/{templateSlug}/parameters` | — |
| PUT | `api/platformhub/{gitRef}/projecttemplates/{templateSlug}/parameters` | — |
| GET | `api/platformhub/{gitRef}/projecttemplates/{templateSlug}/variables` | — |
| PUT | `api/platformhub/{gitRef}/projecttemplates/{templateSlug}/variables` | — |
| DELETE | `api/platformhub/{gitref}/processtemplates/{slug}` | — |
| PUT | `api/platformhub/{gitref}/processtemplates/{slug}` | — |
| DELETE | `api/platformhub/{gitref}/projecttemplates/{slug}` | — |
| GET | `api/progression/runbooks/taskRuns` | — |
| GET | `api/progression/runbooks/{runbookId}` | — |
| GET | `api/progression/runbooks/{runbookId}/v1` | — |
| GET | `api/progression/{projectId}` | — |
| GET | `api/projectgroups` | octopusdeploy_project_group |
| POST | `api/projectgroups` | octopusdeploy_project_group |
| GET | `api/projectgroups/all` | octopusdeploy_project_group |
| DELETE | `api/projectgroups/{id}` | octopusdeploy_project_group |
| GET | `api/projectgroups/{id}` | octopusdeploy_project_group |
| PUT | `api/projectgroups/{id}` | octopusdeploy_project_group |
| GET | `api/projectgroups/{id}/projects` | octopusdeploy_project_group |
| GET | `api/projects` | octopusdeploy_project |
| POST | `api/projects` | octopusdeploy_project |
| GET | `api/projects/all` | octopusdeploy_project |
| GET | `api/projects/experimental/summaries` | octopusdeploy_project |
| GET | `api/projects/featuretoggles/clientidentifier` | octopusdeploy_project |
| POST | `api/projects/import-export/export` | octopusdeploy_project |
| POST | `api/projects/import-export/import` | octopusdeploy_project |
| GET | `api/projects/import-export/import-files` | octopusdeploy_project |
| POST | `api/projects/import-export/import-files` | octopusdeploy_project |
| POST | `api/projects/import-export/import/preview` | octopusdeploy_project |
| GET | `api/projects/{id}/logo` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}` | octopusdeploy_project |
| GET | `api/projects/{projectId}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/channels` | octopusdeploy_project |
| POST | `api/projects/{projectId}/channels` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/channels/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/channels/{id}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/channels/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/channels/{id}/releases` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/channels/{id}/v2` | octopusdeploy_project |
| GET | `api/projects/{projectId}/deploymentprocesses` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/deploymentprocesses` | octopusdeploy_project |
| GET | `api/projects/{projectId}/deploymentprocesses/template` | octopusdeploy_project |
| GET | `api/projects/{projectId}/deploymentprocesses/v2` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/deploymentprocesses/v2` | octopusdeploy_project |
| POST | `api/projects/{projectId}/deploymentprocesses/validate` | octopusdeploy_project |
| GET | `api/projects/{projectId}/deploymentsettings` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/deploymentsettings` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/branches` | octopusdeploy_project |
| POST | `api/projects/{projectId}/git/branches/` | octopusdeploy_project |
| POST | `api/projects/{projectId}/git/branches/v2/` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/branches/{branchName}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/commits/{hash}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/compatibility-report` | octopusdeploy_project |
| POST | `api/projects/{projectId}/git/connectivity-test` | octopusdeploy_project |
| POST | `api/projects/{projectId}/git/convert` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/refs/{refName}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/tags` | octopusdeploy_project |
| GET | `api/projects/{projectId}/git/tags/{tagName}` | octopusdeploy_project |
| POST | `api/projects/{projectId}/git/validate` | octopusdeploy_project |
| POST | `api/projects/{projectId}/logo` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/logo` | octopusdeploy_project |
| GET | `api/projects/{projectId}/metadata` | octopusdeploy_project |
| GET | `api/projects/{projectId}/progression` | octopusdeploy_project |
| GET | `api/projects/{projectId}/progression/v1` | octopusdeploy_project |
| GET | `api/projects/{projectId}/releases` | octopusdeploy_project |
| GET | `api/projects/{projectId}/releases/{id}/variables` | octopusdeploy_project |
| GET | `api/projects/{projectId}/releases/{version}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookProcesses` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookProcesses/{id}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/runbookProcesses/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookProcesses/{id}/runbookSnapshotTemplate` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookProcesses/{id}/v2` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/runbookProcesses/{id}/v2` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookRuns` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbookRuns` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookRuns/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbookSnapshots` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots/{idOrName}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/runbookSnapshots/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots/{id}/runbookRuns` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/preview/{environmentId}/{tenant}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots/{id}/runbookRuns/template` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbookSnapshots/{id}/snapshot-variables/v1` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbookSnapshots/{id}/variables` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/runbookruns/{id}` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbookruns/{runbookRunId}/retry/v1` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbooks` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/all` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/all/v2` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/runbooks/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{id}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/runbooks/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{id}/environments` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{id}/runbookRunTemplate` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{id}/runbookSnapshots` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbooks/{runbookId}/run` | octopusdeploy_project |
| POST | `api/projects/{projectId}/runbooks/{runbookId}/runbookRuns/previews` | octopusdeploy_project |
| GET | `api/projects/{projectId}/runbooks/{runbookId}/runbookSnapshotTemplate` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/runbooksnapshots/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/summary` | octopusdeploy_project |
| GET | `api/projects/{projectId}/summary/v1` | octopusdeploy_project |
| GET | `api/projects/{projectId}/triggers` | octopusdeploy_project |
| POST | `api/projects/{projectId}/triggers` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/triggers/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/triggers/{id}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/triggers/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/deploymentprocesses` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/{gitRef}/deploymentprocesses` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/deploymentprocesses/template` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/deploymentprocesses/v2` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/{gitRef}/deploymentprocesses/v2` | octopusdeploy_project |
| POST | `api/projects/{projectId}/{gitRef}/deploymentprocesses/validate` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/deploymentsettings` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/{gitRef}/deploymentsettings` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbookProcesses/{id}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/{gitRef}/runbookProcesses/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbookProcesses/{id}/v2` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/{gitRef}/runbookProcesses/{id}/v2` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks` | octopusdeploy_project |
| DELETE | `api/projects/{projectId}/{gitRef}/runbooks/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks/{id}` | octopusdeploy_project |
| PUT | `api/projects/{projectId}/{gitRef}/runbooks/{id}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks/{id}/environments` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks/{id}/preview/{environment}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRunTemplate` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/runbooks/{id}/runbookRuns/preview/{environment}/{tenant}` | octopusdeploy_project |
| POST | `api/projects/{projectId}/{gitRef}/runbooks/{runbookId}/run/v1` | octopusdeploy_project |
| POST | `api/projects/{projectId}/{gitRef}/runbooks/{runbookId}/runbookRuns/previews` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/summary` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitRef}/summary/v1` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{gitref}/runbooks/{runbookId}/runbookSnapshotTemplate` | octopusdeploy_project |
| GET | `api/projects/{projectId}/{unusedGitRef}` | octopusdeploy_project |
| GET | `api/projecttriggers` | octopusdeploy_git_trigger, octopusdeploy_built_in_trigger |
| POST | `api/projecttriggers` | octopusdeploy_git_trigger, octopusdeploy_built_in_trigger |
| DELETE | `api/projecttriggers/{id}` | octopusdeploy_git_trigger, octopusdeploy_built_in_trigger |
| GET | `api/projecttriggers/{id}` | octopusdeploy_git_trigger, octopusdeploy_built_in_trigger |
| PUT | `api/projecttriggers/{id}` | octopusdeploy_git_trigger, octopusdeploy_built_in_trigger |
| GET | `api/proxies` | octopusdeploy_machine_proxy |
| POST | `api/proxies` | octopusdeploy_machine_proxy |
| GET | `api/proxies/all` | octopusdeploy_machine_proxy |
| DELETE | `api/proxies/{id}` | octopusdeploy_machine_proxy |
| GET | `api/proxies/{id}` | octopusdeploy_machine_proxy |
| ... | _(and 1654 more)_ | |
