# Octopus Deploy API vs Terraform Provider Coverage

Comparison of the [Octopus Deploy REST API](https://deploy.octopus.app/api/swagger.json) (Swagger 2026.1) with this Terraform provider.

**For full coverage:**  
- **Swagger-based:** [API_COVERAGE_FULL.md](API_COVERAGE_FULL.md) — all 1,314 paths, 102 path resources, 233 API *Resource definitions.  
- **Server source (canonical):** [API_COVERAGE_SERVER_SOURCE.md](API_COVERAGE_SERVER_SOURCE.md) — endpoints and contract properties from the Octopus Server C# source (controllers + Resource/Command/Request). Use this to catch properties missing from Swagger (e.g. Generic OIDC `CustomClaims`).  
  That document includes **relevant files** (scripts, inputs, TF schemas) and **how the data is generated** (extract Server source → JSON; compare → report).

The document below is a curated summary. The provider uses a **muxed architecture**: Framework resources in `octopusdeploy_framework/` (preferred for new work) and legacy SDK resources in `octopusdeploy/`.

---

## 1. Missing resources (API has CRUD; provider has no equivalent)

These API areas have list/create/update/delete operations but no corresponding Terraform resource.

| API area | Path pattern | Notes |
|----------|--------------|--------|
| **User API Keys** | `POST/GET/DELETE /users/{userId}/apikeys` | Per-user API keys; no `octopusdeploy_user_apikey` resource. |
| **Access Tokens** | `POST /users/access-token` | Generate token (operation, not a resource). |
| **Artifacts** | `GET/POST /spaces/.../artifacts` | Deployment artifacts; usually not managed as Terraform resources. |
| **Build information** | `GET/POST /spaces/.../build-information` | Build metadata for deployments. |
| **Dashboard configuration** | `GET/PUT /spaces/.../dashboardconfiguration` | Space dashboard layout/settings. |
| **Deployments** | `GET/POST /spaces/.../deployments` | Create/read deployments; often desired for automation but not yet a resource. |
| **Deployment settings** | `GET/PUT .../projects/.../deploymentsettings` | Project-level deployment settings (e.g. config). |
| **Deployment target tags** | `GET/PUT /spaces/.../deploymenttargettags` | Tag set values used for deployment targeting. |
| **Invitations** | Invitations tag | User invitations. |
| **Releases** | `GET/POST .../projects/.../releases` | Create/read releases; no `octopusdeploy_release` resource. |
| **Runbook runs** | RunbookRuns tag | Execute runbooks / run history; operational, not typically a resource. |
| **Runbook snapshots** | RunbookSnapshots tag | Snapshots of runbook versions. |
| **Subscriptions** | Subscriptions tag | Notifications (e.g. email, webhook). |
| **Interruptions** | Interruptions tag | Manual intervention; action-oriented. |
| **Tasks** | Tasks tag | Task queue / logs; read-only or operational. |
| **Events** | `GET /spaces/.../events` | Audit/events; read-only. |
| **Ephemeral environments** | `.../projects/.../environments/ephemeral` | Provision/deprovision ephemeral environments; project-scoped operations. |
| **Platform Hub policies** | `GET/PUT .../platformhub/.../policies` | Policy CRUD in Platform Hub. |
| **Platform Hub process templates** | `.../platformhub/.../processtemplates` | Process template versions/sharing. |
| **Platform Hub project templates** | `.../platformhub/.../projecttemplates` | Project template sharing. |
| **Reporting** | `GET .../reporting/...` | Reporting endpoints; read-only. |
| **Retention policies** | `GET/PUT .../retentionpolicies/{id}` | Generic retention policy; provider has space default retention policies only. |
| **Configuration** | `GET/PUT /configuration`, `.../configuration/{id}/values` | Server/extension configuration. |
| **Features configuration** | `GET/PUT /featuresconfiguration` | Feature flags. |
| **External security group providers** | `GET /externalsecuritygroupproviders` | Directory/external group providers. |
| **Performance configuration** | `GET/PUT /performanceconfiguration` | Performance tuning. |
| **Maintenance configuration** | MaintenanceConfiguration tag | Maintenance windows. |
| **Licenses** | Licenses tag | License info. |
| **Migrations** | Migrations tag | Migration operations. |

---

## 2. Resources only in SDK (legacy); not in Framework

Per `CLAUDE.md`, new resources must be in the Framework; the SDK is legacy. These exist only as SDK resources and have no Framework counterpart yet.

| Terraform resource (SDK) | API area | Note |
|--------------------------|----------|------|
| `octopusdeploy_machine_policy` | MachinePolicies | Machine policy for deployment targets. |
| `octopusdeploy_static_worker_pool` | WorkerPools | Static worker pool. |
| `octopusdeploy_dynamic_worker_pool` | WorkerPools | Dynamic worker pool. |
| `octopusdeploy_listening_tentacle_deployment_target` | Machines | Listening tentacle deployment target. |
| `octopusdeploy_polling_tentacle_deployment_target` | Machines | Polling tentacle deployment target. |
| `octopusdeploy_ssh_connection_deployment_target` | Machines | SSH deployment target. |
| `octopusdeploy_cloud_region_deployment_target` | Machines | Cloud region target. |
| `octopusdeploy_kubernetes_agent_deployment_target` | Machines | Kubernetes agent target. |
| `octopusdeploy_kubernetes_agent_worker` | Workers | Kubernetes agent worker. |
| `octopusdeploy_kubernetes_cluster_deployment_target` | Machines | Kubernetes cluster target. |
| `octopusdeploy_azure_cloud_service_deployment_target` | Machines | Azure cloud service target. |
| `octopusdeploy_azure_service_fabric_cluster_deployment_target` | Machines | Azure Service Fabric target. |
| `octopusdeploy_azure_web_app_deployment_target` | Machines | Azure Web App target. |
| `octopusdeploy_offline_package_drop_deployment_target` | Machines | Offline package drop target. |
| `octopusdeploy_deployment_process` | DeploymentProcesses | Deployment process (SDK); Framework has Process + steps. |
| `octopusdeploy_runbook_process` | RunbookProcesses | Runbook process (SDK). |
| `octopusdeploy_project_scheduled_trigger` | ProjectTriggers | Scheduled trigger. |
| `octopusdeploy_project_deployment_target_trigger` | ProjectTriggers | Deployment target trigger. |
| `octopusdeploy_external_feed_create_release_trigger` | ProjectTriggers | External feed create release trigger. |
| `octopusdeploy_user_role` | UserRoles | User role (SDK); Framework has Team + ScopedUserRole. |
| `octopusdeploy_token_account` | Accounts | Token account. |
| `octopusdeploy_ssh_key_account` | Accounts | SSH key account. |
| `octopusdeploy_gcp_account` | Accounts | GCP account (SDK); Framework has Platform Hub GCP. |
| `octopusdeploy_aws_openid_connect_account` | Accounts | AWS OIDC (SDK); Framework has Platform Hub AWS OIDC. |
| `octopusdeploy_azure_openid_connect` | Accounts | Azure OIDC (SDK); Framework has Platform Hub Azure OIDC. |
| `octopusdeploy_polling_subscription_id` | Subscriptions | Polling subscription. |

Framework already has: workers (listening tentacle, SSH, Kubernetes monitor), process/process steps, project triggers (Git, built-in), accounts (username/password, certificate, AWS, Azure, generic OIDC, Platform Hub variants), and others. The table above is the set that is **only** in SDK.

---

## 3. Missing or partial properties on existing resources

### AccountResource (API) vs account resources (provider)

| API property | In provider? | Note |
|-------------|-------------|------|
| `EnvironmentIds` | Partial | Scoping; not all account resources expose it. |
| `TenantIds` | Partial | Scoping; not all account resources expose it. |
| `TenantTags` | Yes | Exposed on multiple account types. |
| `TenantedDeploymentParticipation` | Yes | Exposed on multiple account types. |

### ProjectResource (API) vs `octopusdeploy_project`

| API property | In provider? | Note |
|-------------|-------------|------|
| `CombineHealthAndSyncStatusInDashboardLiveStatus` | No | Dashboard live status behavior. |
| `ExecuteDeploymentsOnResilientPipeline` | No | Resilient pipeline option. |
| `DefaultPowerShellEdition` | No | Default PowerShell edition for the project. |
| `ForcePackageDownload` | No (on project) | Exists on runbook; not on project. |
| `AllowIgnoreChannelRules` | No | Allow ignoring channel rules. |
| `ReleaseCreationStrategy` | Partial | Handled via `octopusdeploy_project_auto_create_release`. |
| `VersioningStrategy` | Partial | Handled via `octopusdeploy_project_versioning_strategy`. |
| `ProjectConnectivityPolicy` | Yes | As `connectivity_policy` (including `exclude_unhealthy_targets`). |
| `ProjectTemplateDetails` | No | Resolved template details (read-only). |

### EnvironmentResource (API) vs `octopusdeploy_environment`

| API property | In provider? | Note |
|-------------|-------------|------|
| `EnvironmentTags` | Yes | As `environment_tags`. |
| `AllowDynamicInfrastructure` | Yes | As `allow_dynamic_infrastructure`. |
| `UseGuidedFailure` | Yes | As `use_guided_failure`. |
| `SortOrder` | Yes | As `sort_order`. |
| `ExtensionSettings` | No | Read-only in API; extension-specific. |

### LifecycleResource (API) vs `octopusdeploy_lifecycle`

| API property | In provider? | Note |
|-------------|-------------|------|
| `Phases` | Read-only in API | Provider manages lifecycle with phases. |
| `ReleaseRetentionPolicy` | Yes | Via lifecycle + space default release retention. |
| `TentacleRetentionPolicy` | Yes | Via lifecycle + space default tentacle retention. |

### MachineResource / deployment targets (API)

- Full **Machine** CRUD exists in API; in provider, “machines” are represented as:
  - **Framework:** Workers (listening tentacle, SSH, Kubernetes monitor) and `octopusdeploy_kubernetes_monitor`.
  - **SDK:** Various `octopusdeploy_*_deployment_target` and `octopusdeploy_*_worker` resources.
- API fields such as `HasLatestCalamari`, `HealthStatus`, `SkipInitialHealthCheck`, `ShellName`, `ShellVersion` are not necessarily exposed on every provider resource.

### WorkerPoolResource (API) vs worker pool resources

- **SDK:** `octopusdeploy_static_worker_pool`, `octopusdeploy_dynamic_worker_pool` exist.
- **Framework:** No dedicated worker pool resource; worker pool *members* (workers) are in Framework.
- API properties: `CanAddWorkers`, `IsDefault`, `SortOrder`, `WorkerPoolType` — check SDK schema for exact coverage.

### Certificate (API) vs `octopusdeploy_certificate`

- Certificate resources in provider; API also has `TenantIds`, `TenantTags`, `EnvironmentIds` for scoping — verify all are exposed where applicable.

### ChannelResource (API) vs `octopusdeploy_channel`

- Provider has `tenant_tags` and channel rules; confirm any newer API fields (e.g. Git-related rules) are fully mapped.

---

## 4. Summary

- **Missing resources:** User API Keys, Dashboard configuration, Deployments, Deployment settings, Deployment target tags, Releases, Runbook runs/snapshots, Subscriptions, Invitations, and several Platform Hub, reporting, and configuration areas. Some (e.g. Artifacts, Events, Tasks) are read-only or operational and may be low priority for resources.
- **SDK-only (no Framework):** Machine policies, static/dynamic worker pools, all deployment target types, deployment/runbook process (SDK), project triggers (scheduled, deployment target, external feed), user role, token account, SSH key account, GCP/AWS OIDC/Azure OIDC accounts (SDK), polling subscription. Migration path: implement or re-expose these in Framework as needed.
- **Missing properties:** On **Project:** `CombineHealthAndSyncStatusInDashboardLiveStatus`, `ExecuteDeploymentsOnResilientPipeline`, `DefaultPowerShellEdition`, `ForcePackageDownload`, `AllowIgnoreChannelRules`. On **Accounts:** ensure `EnvironmentIds`/`TenantIds` are exposed where the API supports them. **Environment/Lifecycle/Certificate/Channel** are largely covered; confirm new API fields as the API evolves.

Use this document to prioritize new resources (e.g. `octopusdeploy_release`, `octopusdeploy_deployment`, dashboard configuration) and new attributes (e.g. on project and accounts) when aligning the provider with the Octopus Deploy API.
