resource "octopusdeploy_aws_openid_connect_account" "example" {
  name = "AWS OIDC Account"
  description = "AWS OIDC Connect Account"
  role_arn = "Amazon Resource Name"
  account_test_subject_keys = ["space"]
  environments = ["environment-123"]
  tenants = ["tenants-123"]
  execution_subject_keys = ["space"]
  custom_claims = {
    "claim1" = "value1"
    "claim2" = "{\"nestedClaim1\":\"value2\",\"nestedClaim2\":\"value3\"}"
    "claim3" = "[\"value4\",\"value5\",\"value6\"]"
  }
}