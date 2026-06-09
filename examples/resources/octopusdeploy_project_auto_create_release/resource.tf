# Example configuration for octopusdeploy_project_auto_create_release resource

# Basic project setup
resource "octopusdeploy_project_group" "example" {
  description = "Example project group for auto create release demo"
  name        = "Example Project Group"
}

resource "octopusdeploy_lifecycle" "example" {
  description = "Example lifecycle"
  name        = "Example Lifecycle"

  release_retention_with_strategy {
    strategy         = "Count"
    quantity_to_keep = 30
    unit             = "Days"
  }
}

resource "octopusdeploy_environment" "development" {
  description = "Development environment"
  name        = "Development"
}

resource "octopusdeploy_project" "example" {
  description      = "Example project with auto create release"
  lifecycle_id     = octopusdeploy_lifecycle.example.id
  name             = "Example Project with Auto Create Release"
  project_group_id = octopusdeploy_project_group.example.id

  # Note: auto_create_release is NOT set here - it will be managed by the separate resource
}

# Get the built-in feed
data "octopusdeploy_feeds" "built_in" {
  feed_type = "BuiltIn"
  take      = 1
}

# Channel for the project
resource "octopusdeploy_channel" "default" {
  description  = "Auto release channel"
  lifecycle_id = octopusdeploy_lifecycle.example.id
  name         = "Auto Release Channel"
  project_id   = octopusdeploy_project.example.id
}

# Deployment process with a package step that uses built-in feed
resource "octopusdeploy_process" "example" {
  project_id = octopusdeploy_project.example.id
}

resource "octopusdeploy_process_step" "example" {
  process_id   = octopusdeploy_process.example.id
  name         = "Deploy Package Action"
  type         = "Octopus.Script"
  condition    = "Success"
  environments = [octopusdeploy_environment.development.id]

  execution_properties = {
    "Octopus.Action.RunOnServer"         = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody"   = "Write-Host \"Deploying MyApp\""
  }

  # Auto create release points at a package by its reference name, so use a
  # named package reference (the map key is the reference name).
  packages = {
    "MyApp" = {
      package_id           = "MyApp"
      acquisition_location = "Server"
      feed_id              = data.octopusdeploy_feeds.built_in.feeds[0].id
    }
  }
}

# Auto create release configuration
resource "octopusdeploy_project_auto_create_release" "example" {
  deployment_process_id = octopusdeploy_process.example.id
  channel_id            = octopusdeploy_channel.default.id

  release_creation_package {
    deployment_action = octopusdeploy_process_step.example.name
    package_reference = "MyApp"
  }

  # release_creation_package_step_id is computed automatically if not provided
}