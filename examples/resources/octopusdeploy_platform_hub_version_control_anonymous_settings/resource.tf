resource "octopusdeploy_platform_hub_version_control_anonymous_settings" "example" {
  url            = "https://github.com/acme/hello-world.git"
  default_branch = "main"
  base_path      = ".octopus"
}
