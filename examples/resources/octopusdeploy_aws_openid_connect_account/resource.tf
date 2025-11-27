resource "octopusdeploy_aws_openid_connect_account" "example" {
  name = "AWS OIDC Account"
  description = "AWS OIDC Connect Account"
  role_arn = "Amazon Resource Name"
  account_test_subject_keys = ["space"]
  environments = ["environment-123"]
  tenants = ["tenants-123"]
  execution_subject_keys = ["space"]
  custom_claims = {
    "custom_claim_1" = "value1"
    "custom_claim_2" = "value2"
  }
}