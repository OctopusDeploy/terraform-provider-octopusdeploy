resource "octopusdeploy_platform_hub_aws_openid_connect_account" "example" {
  name                       = "AWS OIDC Account"
  description                = "AWS OIDC Connect Account"
  role_arn                   = "arn:aws:iam::123456789012:role/MyRole"
  account_test_subject_keys  = ["space", ]
  execution_subject_keys     = ["space", "environment"]
  health_subject_keys        = ["space"]
  session_duration           = "3600"
}
