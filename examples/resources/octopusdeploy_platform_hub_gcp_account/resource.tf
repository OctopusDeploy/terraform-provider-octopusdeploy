resource "octopusdeploy_platform_hub_gcp_account" "example" {
  name        = "My GCP Account"
  description = "My GCP account managed by terraform"
  json_key    = file("path/to/service-account-key.json")
}
