resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle"
}

resource "octopusdeploy_project_group" "example" {
  description  = "The project group."
  name         = "Version Controlled Project Group"
}