resource "octopusdeploy_platform_hub_git_credential" "example" {
  name        = "GitHub Credentials"
  description = "Git credentials for production repositories"
  username    = "my-github-username"
  password    = "ghp_yourPersonalAccessTokenHere"

  repository_restrictions = {
    enabled = true
    allowed_repositories = [
      "https://github.com/myorg/myrepo"
    ]
  }
}
