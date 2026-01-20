resource "octopusdeploy_platform_hub_generic_oidc_account" "example" {
  name                   = "Generic OpenID Connect Account (OK to Delete)"
  description            = "My Generic OIDC account managed by terraform"
  execution_subject_keys = ["space", "project"]
  audience               = "api://default"
}
