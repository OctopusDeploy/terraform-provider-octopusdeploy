resource "octopusdeploy_user" "github_actions" {
  display_name  = "GitHub Actions"
  email_address = "github-actions@example.com"
  is_active     = true
  is_service    = true
  username      = "github-actions"
}

resource "octopusdeploy_service_account_oidc_identity" "github_actions" {
  name               = "GitHub Actions main branch"
  service_account_id = octopusdeploy_user.github_actions.id
  issuer             = "https://token.actions.githubusercontent.com"
  subject            = "repo:my-org/my-repo:*"

  # The token must pass every condition. Each condition passes when the claim matches any one of its values.
  # Leave this out to keep the existing conditions. Set it to [] to remove them all.
  claim_conditions = [
    {
      claim  = "repository_owner"
      values = ["my-org"]
    },
    {
      claim  = "ref"
      values = ["refs/heads/main", "refs/tags/v*"]
    },
  ]
}
