resource "octopusdeploy_tenant" "three"{
    name = "New Tenant Three"
    tenant_tags = [octopusdeploy_tag.alpha.canonical_tag_name, octopusdeploy_tag.beta.canonical_tag_name]
}

resource "octopusdeploy_azure_service_principal" "three" {
  application_id  = "00000000-0000-0000-0000-000000000001"
  name            = "Azure Service Principal Account (OK to Delete)"
  password        = "###########" # required; get from secure environment/store
  subscription_id = "00000000-0000-0000-0000-000000000002"
  tenant_id       = "00000000-0000-0000-0000-000000000003"
  tenants = [octopusdeploy_tenant.three.id]
}

resource "octopusdeploy_azure_web_app_deployment_target" "web_app" {
    name = "Azure Web App Deployment Target"
    roles = [data.octopusdeploy_teams.everyone.id]
    account_id = octopusdeploy_azure_service_principal.three.id
    environments = [ octopusdeploy_environment.dev.id, octopusdeploy_environment.prod.id  ]
    web_app_name = "Test Web App"
    resource_group_name = "resource-group-name"
}

resource "octopusdeploy_azure_web_app_deployment_target" "static_app" {
    name = "Azure Static Web App Deployment Target"
    roles = [data.octopusdeploy_teams.everyone.id]
    account_id = octopusdeploy_azure_service_principal.three.id
    environments = [ octopusdeploy_environment.dev.id, octopusdeploy_environment.prod.id  ]
    web_app_name = "Test Static Web App"
    resource_group_name = "resource-group-name"
}