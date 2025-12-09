resource "octopusdeploy_platform_hub_aws_account" "example" {
    name        = "My AWS Account"
    description = "My AWS account managed by terraform"
    access_key  = "MY-ACCESS-KEY"
    secret_key  = "my-secret-key"
 }