resource "octopusdeploy_space" "example" {
  name                 = "Dev Team Space"
  space_managers_teams = ["teams-everyone"]
}

resource "octopusdeploy_environment" "development" {
  name     = "Development"
  space_id = octopusdeploy_space.example.id
}

resource "octopusdeploy_environment" "production" {
  name     = "Production"
  space_id = octopusdeploy_space.example.id
}

# Configures the "Default Lifecycle" that Octopus creates together with the space
resource "octopusdeploy_space_default_lifecycle" "example" {
  space_id    = octopusdeploy_space.example.id
  name        = "Default Lifecycle"
  description = "Managed by Terraform."

  phases = [
    {
      name                         = "Development"
      automatic_deployment_targets = [octopusdeploy_environment.development.id]
    },
    {
      name                        = "Production"
      optional_deployment_targets = [octopusdeploy_environment.production.id]
      release_retention_with_strategy = {
        strategy = "Forever"
      }
    },
  ]

  release_retention_with_strategy = {
    strategy         = "Count"
    quantity_to_keep = 30
    unit             = "Days"
  }

  tentacle_retention_with_strategy = {
    strategy = "Default"
  }
}
