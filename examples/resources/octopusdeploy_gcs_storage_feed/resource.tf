resource "octopusdeploy_gcs_storage_feed" "example" {
  name                           = "GCS Storage Feed (OK to Delete)"
  project                        = "my-gcp-project-id"
  use_service_account_key        = true
  service_account_json_key       = jsonencode({
    "type": "service_account",
    "project_id": "my-project",
    "private_key_id": "key-id",
    "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
    "client_email": "service-account@my-project.iam.gserviceaccount.com",
    "client_id": "123456789",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs"
  })
  download_attempts              = 5
  download_retry_backoff_seconds = 10
}

resource "octopusdeploy_gcs_storage_feed" "example_with_oidc" {
  name                    = "GCS Storage Feed with OIDC (OK to Delete)"
  project                 = "my-gcp-project-id"
  use_service_account_key = false
  oidc_authentication = {
    audience     = "https://gcs.googleapis.com"
    subject_keys = ["feed", "space"]
  }
}
