data "octopusdeploy_teams" "everyone" {
  partial_name = "Everyone"
  skip         = 0
  take         = 1
}

# Create a Space
resource "octopusdeploy_space" "example" {
  name        = var.space.name
  description = var.space.description
  space_managers_teams = [ data.octopusdeploy_teams.everyone.teams[0].id ]
}

# Create 3 Environments
resource "octopusdeploy_environment" "development" {
  name        = "Development"
  description = "Development environment"
  space_id    = octopusdeploy_space.example.id
}

resource "octopusdeploy_environment" "staging" {
  name        = "Staging"
  description = "Staging environment"
  space_id    = octopusdeploy_space.example.id
}

resource "octopusdeploy_environment" "production" {
  name        = "Production"
  description = "Production environment"
  space_id    = octopusdeploy_space.example.id
}

# Create a Container Registry
resource "octopusdeploy_docker_container_registry" "example" {
  name        = "DockerHub"
  space_id    = octopusdeploy_space.example.id
  feed_uri    = "https://registry.hub.docker.com"
}

# Additional required resources
resource "octopusdeploy_lifecycle" "example" {
  name     = "Example Lifecycle"
  space_id = octopusdeploy_space.example.id
}

resource "octopusdeploy_project_group" "example" {
  name     = "Example Project Group"
  space_id = octopusdeploy_space.example.id
}

# Create a Project
resource "octopusdeploy_project" "example" {
  name        = var.project.name
  description = var.project.description
  space_id    = octopusdeploy_space.example.id
  lifecycle_id = octopusdeploy_lifecycle.example.id
  project_group_id = octopusdeploy_project_group.example.id
}