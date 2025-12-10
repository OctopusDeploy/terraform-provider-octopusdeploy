resource "octopusdeploy_platform_hub_username_password_account" "example" {
  name        = "Username-Password Account (OK to Delete)"
  description = "My Username-Password account managed by terraform"
  username    = "myusername"
  password    = "mypassword123"
}
