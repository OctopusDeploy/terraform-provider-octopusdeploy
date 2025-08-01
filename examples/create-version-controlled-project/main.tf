resource "octopusdeploy_git_credential" "app"{
    name = "Git credential"
    password = "Password"
    username = "Username"
}

resource "octopusdeploy_project" "app" {
  name                                 = "Version-Controlled Project "
  project_group_id                     = octopusdeploy_project_group.example.id
  lifecycle_id                         = octopusdeploy_lifecycle.example.id
  is_version_controlled = true
  git_library_persistence_settings {
    git_credential_id = octopusdeploy_git_credential.app.id
    url =               "<git-url>"
    default_branch      = "main"
    base_path           = ".octopus"
    protected_branches  = []
  }
}