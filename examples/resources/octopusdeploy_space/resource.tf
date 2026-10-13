resource "octopusdeploy_space" "example" {
  description                 = "A space for the development team."
  name                        = "Dev Team Space"
  is_default                  = false
  is_task_queue_stopped       = false
  space_managers_team_members = ["Users-123", "Users-321"]
  space_managers_teams        = ["teams-everyone"]
}

resource "octopusdeploy_space" "with_default_lifecycle" {
  name                 = "Ops Team Space"
  space_managers_teams = ["teams-everyone"]

  # Configures the "Default Lifecycle" that Octopus creates together with the space
  default_lifecycle = {
    name        = "Default Lifecycle"
    description = "Managed by Terraform."

    phases = [
      {
        name                         = "Development"
        automatic_deployment_targets = ["Environments-1"]
      },
      {
        name                        = "Production"
        optional_deployment_targets = ["Environments-2"]
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
}
