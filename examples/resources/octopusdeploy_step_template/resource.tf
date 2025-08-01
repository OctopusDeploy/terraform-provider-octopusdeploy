resource "octopusdeploy_step_template" "example" {
  action_type     = "Octopus.Script"
  name            = "Example"
  description     = "Example script step template"
  step_package_id = "Octopus.Script"
  packages        = [
    {
      package_id = "my.scripts"
      acquisition_location = "Server"
      feed_id = module.predefined.feeds.built_in
      name = "My Scripts"
      properties = {
        extract = "True"
        purpose = ""
        selection_mode = "immediate"
      }
    }
  ]

  parameters = [
    {
      name      = "Text.Plain"
      id = "10001000-0000-0000-0000-100010001001"
      label     = "Plain text"
      help_text = "Random text value"
      display_settings = {
        "Octopus.ControlType" : "SingleLineText"
      }
      default_value = "initial text"
    },
    {
      name      = "My.Secret"
      id = "10001000-0000-0000-0000-100010001002"
      label     = "Secret value"
      help_text = "Some secret value"
      display_settings = {
        "Octopus.ControlType" : "Sensitive"
      }
      default_sensitive_value = var.secrets.example
    },
  ]

  properties = {
    "Octopus.Action.Script.ScriptBody" : "echo '1.#{Text.Plan}'"
    "Octopus.Action.Script.ScriptSource" : "Inline"
    "Octopus.Action.Script.Syntax" : "Bash"
  }
}
