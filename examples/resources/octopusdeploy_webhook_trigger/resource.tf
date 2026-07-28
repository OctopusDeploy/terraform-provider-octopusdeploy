resource "octopusdeploy_webhook_trigger" "my_trigger" {
  name        = "My webhook trigger"
  space_id    = "Spaces-1"
  description = "Runs a runbook when the webhook is called"
  project_id  = "Projects-1"
  secret      = "secret"

  run_runbook_action = {
    runbook_id             = "Runbooks-1"
    target_environment_ids = ["Environments-1"]
  }
}

# Authenticate with an Octopus API key instead of a shared secret
resource "octopusdeploy_webhook_trigger" "my_api_key_trigger" {
  name            = "My API key webhook trigger"
  space_id        = "Spaces-1"
  description     = "Runs a runbook when the webhook is called"
  project_id      = "Projects-1"
  require_api_key = true

  run_runbook_action = {
    runbook_id             = "Runbooks-1"
    target_environment_ids = ["Environments-1"]
  }
}
