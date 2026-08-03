data "octopusdeploy_step_templates" "example" {
  ids          = ["ActionTemplates-123", "ActionTemplates-321"]
  partial_name = "Deploy"
  skip         = 5
  take         = 100
}

# Build a name to template lookup, so processes can reference templates by name instead of
# hard-coding an ID and a version that changes every time the template is edited.
data "octopusdeploy_step_templates" "shared" {
  partial_name = "Shared - "
}

locals {
  shared_templates = {
    for template in data.octopusdeploy_step_templates.shared.step_templates :
    template.name => template
  }
}

resource "octopusdeploy_process_templated_step" "deploy" {
  process_id       = octopusdeploy_process.example.id
  name             = "Deploy the application"
  template_id      = local.shared_templates["Shared - Deploy Application"].id
  template_version = local.shared_templates["Shared - Deploy Application"].version
}
