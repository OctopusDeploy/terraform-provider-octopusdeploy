data "octopusdeploy_step_template" "example" {
    id = "<Source-ActionTemplates-Id>"
    space_id = "<Source-Spaces-Id>"

}

resource "octopusdeploy_step_template" "hello_blank_step_template" {
  action_type     = "Octopus.Script"
  name            = data.octopusdeploy_step_template.example.step_template.name
  description     = data.octopusdeploy_step_template.example.step_template.description
  space_id        = "<Target-Spaces-Id>"
  step_package_id = "Octopus.Script"
  packages        = []
  parameters = [
    {
      default_value = data.octopusdeploy_step_template.example.step_template.parameters[0].default_value
      display_settings = {
        "Octopus.ControlType" : "SingleLineText"
      }
      help_text = data.octopusdeploy_step_template.example.step_template.parameters[0].help_text
      label     = data.octopusdeploy_step_template.example.step_template.parameters[0].label
      name      = data.octopusdeploy_step_template.example.step_template.parameters[0].name
      id = data.octopusdeploy_step_template.example.step_template.parameters[0].id
    },
  ]
  properties = {
   "Octopus.Action.Script.ScriptBody" : data.octopusdeploy_step_template.example.step_template.properties["Octopus.Action.Script.ScriptBody"]
   "Octopus.Action.Script.ScriptSource" : data.octopusdeploy_step_template.example.step_template.properties["Octopus.Action.Script.ScriptSource"]
   "Octopus.Action.Script.Syntax" : data.octopusdeploy_step_template.example.step_template.properties["Octopus.Action.Script.Syntax"]
  }
}