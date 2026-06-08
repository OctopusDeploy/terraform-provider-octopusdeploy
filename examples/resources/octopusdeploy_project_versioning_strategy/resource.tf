resource "octopusdeploy_project_group" "tp" {
  name        = "DevOps Projects"
  description = "My DevOps projects group"
}

resource "octopusdeploy_project" "tp" {
  name             = "My DevOps Project"
  description      = "test project"
  lifecycle_id     = "Lifecycles-1"
  project_group_id = octopusdeploy_project_group.tp.id

  depends_on  = [octopusdeploy_project_group.tp]
}

resource "octopusdeploy_process" "process" {
  project_id = octopusdeploy_project.tp.id

  depends_on = [octopusdeploy_project.tp]
}

resource "octopusdeploy_process_step" "hello_world" {
  process_id = octopusdeploy_process.process.id
  name       = "Hello World"
  type       = "Octopus.Script"

  properties = {
    "Octopus.Action.TargetRoles" = "hello-world"
  }

  execution_properties = {
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody"   = "Write-Host 'hello world'"
  }

  packages = {
    "Package" = {
      feed_id              = "feeds-builtin"
      package_id           = "myExpressApp"
      acquisition_location = "Server"
    }
  }
}

######
# NOTE: You can use either template or donor_package, not both
######
resource "octopusdeploy_project_versioning_strategy" "using_template" {
  project_id = octopusdeploy_project.tp.id
  space_id = octopusdeploy_project.tp.space_id
  # See https://octopus.com/docs/releases/release-versioning for variable syntax
  template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.NextMinor}-alpha"
  depends_on = [
    octopusdeploy_project_group.tp,
    octopusdeploy_process_step.hello_world
  ]
}

resource "octopusdeploy_project_versioning_strategy" "using_donor_package" {
  project_id = octopusdeploy_project.tp.id
  space_id = octopusdeploy_project.tp.space_id
  donor_package_step_id = octopusdeploy_process_step.hello_world.id
  donor_package = {
    deployment_action = "Hello World"
    package_reference = "Package"
  }
  depends_on = [
    octopusdeploy_project_group.tp,
    octopusdeploy_process_step.hello_world
  ]
}
