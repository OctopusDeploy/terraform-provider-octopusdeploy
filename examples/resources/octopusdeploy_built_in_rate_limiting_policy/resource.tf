resource "octopusdeploy_built_in_rate_limiting_policy" "authenticated_user" {
  slug              = "user"
  is_enabled        = true
  requests_per_hour = 10000
  burst_limit       = 100
  audit_mode        = false
}
