resource "octopusdeploy_platform_hub_azure_oidc_account" "example" {
  name                        = "Azure OpenID Connect Account (OK to Delete)"
  description                 = "Azure OIDC account managed by terraform"
  subscription_id             = "00000000-0000-0000-0000-000000000000"
  application_id              = "11111111-1111-1111-1111-111111111111"
  tenant_id                   = "22222222-2222-2222-2222-222222222222"
  execution_subject_keys      = ["space", "project"]
  health_subject_keys         = ["space", "target", "type"]
  account_test_subject_keys   = ["space", "type"]
  audience                    = "api://AzureADTokenExchange"
  azure_environment           = "AzureCloud"
  authentication_endpoint     = ""
  resource_management_endpoint = ""
}
