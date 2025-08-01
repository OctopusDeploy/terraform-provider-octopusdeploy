resource "octopusdeploy_environment" "dev" {    
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