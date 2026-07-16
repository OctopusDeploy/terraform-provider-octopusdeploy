resource "octopusdeploy_project_group" "example" {
  name        = "Example"
  description = "Example Group"
}

resource "octopusdeploy_project" "example" {
  name                                 = "Example"
  lifecycle_id                         = "Lifecycles-101"
  project_group_id                     = octopusdeploy_project_group.example.id
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Project with Built-In Trigger"
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  tenanted_deployment_participation    = "Untenanted"
  included_library_variable_sets       = []

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_channel" "example" {
    name = "Example Channel"
    project_id = octopusdeploy_project.example.id
    lifecycle_id = "Lifecycles-101"
}

data "octopusdeploy_feeds" "built_in" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
  skip         = 0
  take         = 1
}

resource "octopusdeploy_process" "example" {
  project_id = octopusdeploy_project.example.id
}

resource "octopusdeploy_process_step" "example" {
  process_id = octopusdeploy_process.example.id
  name       = "Action One"
  type       = "Octopus.Script"
  execution_properties = {
    "Octopus.Action.RunOnServer"         = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody"   = <<-EOT
        $ExtractedPath = $OctopusParameters["Octopus.Action.Package[my.package].ExtractedPath"]
        Write-Host $ExtractedPath
      EOT
  }
  packages = {
    "my.package" = {
      package_id           = "my.package"
      feed_id              = data.octopusdeploy_feeds.built_in.feeds[0].id
      acquisition_location = "Server"
      properties = {
        "Extract" = "True"
      }
    }
  }
}

resource "octopusdeploy_built_in_trigger" "example" {
  project_id = octopusdeploy_project.example.id
  channel_id = octopusdeploy_channel.example.id

  release_creation_package = {
    deployment_action = "Action One"
    package_reference = "my.package"
  }

  depends_on = [
    octopusdeploy_project.example,
    octopusdeploy_channel.example,
    octopusdeploy_process_step.example
  ]
}
