resource "octopusdeploy_channel" "example" {
  name       = "Development Channel (OK to Delete)"
  project_id = "Projects-123"

  # Constrain which Git refs can be deployed through this channel.
  git_reference_rules = ["refs/heads/main", "refs/tags/release-*"]

  # Apply Git ref rules to specific Git dependency actions.
  git_resource_rules = [{
    rules = ["refs/heads/main"]

    git_dependency_actions = [{
      deployment_action_slug = "deploy-application"
      git_dependency_name    = "app-config"
    }]
  }]
}
