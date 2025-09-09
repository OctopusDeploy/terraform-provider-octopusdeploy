---
page_title: "Onboard a new Tenant to existing Project"
subcategory: "Examples"
---

# Onboard a new Tenant to existing Project

This example show how to configure new Tenant and connect it to existing projects. This can be used to onboard new customers to existing applications


## Existing ecosystem
These resources are the "pre-requisites" and assumed already exists in Terraform state

```terraform
﻿resource "octopusdeploy_environment" "dev" {    
  name        = "Development"
  description = "Development environment"
}

resource "octopusdeploy_environment" "prod" {
  name        = "Production"
  description = "Production environment"
}

resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle"
}

resource "octopusdeploy_project_group" "example" {
  description  = "The project group."
  name         = "Tenanted Project Group"
}

resource "octopusdeploy_project" "frontend" {
  name                                 = "Frontend"
  project_group_id                     = octopusdeploy_project_group.example.id
  lifecycle_id                         = octopusdeploy_lifecycle.example.id
  tenanted_deployment_participation    = "TenantedOrUntenanted"
}

resource "octopusdeploy_project" "backend" {
  name                                 = "Backend"
  project_group_id                     = octopusdeploy_project_group.example.id
  lifecycle_id                         = octopusdeploy_lifecycle.example.id
  tenanted_deployment_participation    = "TenantedOrUntenanted"
}

resource "octopusdeploy_tag_set" "eap" {
  description = "Provides tenants with access to certain early access programs."
  name        = "Early Access Program (EAP)"
}

resource "octopusdeploy_tag" "alpha" {
  color      = "#00FF00"
  name       = "Redundancy"
  tag_set_id = octopusdeploy_tag_set.eap.id
}

resource "octopusdeploy_tag" "beta" {
  color      = "#FF0000"
  name       = "Industry Category"
  tag_set_id = octopusdeploy_tag_set.eap.id
}

resource "octopusdeploy_library_variable_set" "example"{
    name = "Library Variable set"
}

resource "octopusdeploy_tenant" "one"{
    name = "One"
}

resource "octopusdeploy_tenant" "two"{
    name = "Two"
}
```

## New Tenant
Adds new tenant with corresponding tags and connects to deployment targets

```terraform
﻿resource "octopusdeploy_tenant" "three"{
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
```
