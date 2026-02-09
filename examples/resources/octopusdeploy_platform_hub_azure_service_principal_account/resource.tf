resource "octopusdeploy_platform_hub_azure_service_principal_account" "example" {
  name                        = "Azure Service Principal Account (OK to Delete)"
  description                 = "Azure service principal account managed by terraform"
  subscription_id             = "00000000-0000-0000-0000-000000000000"
  tenant_id                   = "11111111-1111-1111-1111-111111111111"
  application_id              = "22222222-2222-2222-2222-222222222222"
  password                    = "my-secret-password"
  azure_environment           = "AzureCloud"
  authentication_endpoint     = ""
  resource_management_endpoint = ""
}
