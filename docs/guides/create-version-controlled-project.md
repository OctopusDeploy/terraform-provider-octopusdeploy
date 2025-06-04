---
page_title: "Create a new Version-Controlled Project"
subcategory: "Examples"
---

# Create a new Version-Controlled Project

This example show how to configure new version-controlled Project and let application team manage it with OCL

## Existing ecosystem
These resources are the "pre-requisites" and assumed already exists
```terraform
﻿resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle"
}

resource "octopusdeploy_project_group" "example" {
  description  = "The project group."
  name         = "Version Controlled Project Group"
}
```

## Version-Controlled Project
```terraform
﻿resource "octopusdeploy_git_credential" "app"{
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
```
