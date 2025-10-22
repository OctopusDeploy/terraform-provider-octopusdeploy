resource "octopusdeploy_docker_container_registry" "example" {
  feed_uri                       = "https://index.docker.io"
  name                           = "Test Docker Container Registry (OK to Delete)"
  password                       = "test-password"
  registry_path                  = "testing/test-image"
  username                       = "test-username"
  download_attempts              = 3
  download_retry_backoff_seconds = 15
}
