resource "octopusdeploy_parent_environment" "example" {
  name                          = "Terraform Environment"
  space_id                      = "Spaces-1"
  description                   = "A parent environment."
  use_guided_failure            = false
  automatic_deprovisioning_rule = {
    days = 7
    hours = 12
  }
}