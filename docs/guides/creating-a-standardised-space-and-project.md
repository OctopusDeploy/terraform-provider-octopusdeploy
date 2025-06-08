---
page_title: "Creating a Standardised Space and Project"
subcategory: "Examples"
---

# Creating a Standardised Space and Project

This example show how to create a standardised space and project that can be used to quickly setup applications used as
microservices.

## Create Space, Project with Deployment Process
Consuming terraform module with standard configuration for Space and Project

```terraform
﻿module "standard" {
  source       = "../modules/standardised-space-and-project"

  space = {
    name        = "Testing"
    description = "Testing infrastructure"
  }
  
  project = {
    name        = "Test Project"
    description = "Testing resources" 
  }
}

resource "octopusdeploy_process" "test" {
  space_id = module.standard.space
  project_id = module.standard.project
}

# Manual Intervention Step (Production only)
resource "octopusdeploy_process_step" "approve" {
  space_id = module.standard.space
  process_id = octopusdeploy_process.test.id
  name = "Approve"
  type = "Octopus.Manual"
  is_required = true
  environments = [
    module.standard.environments.production
  ]
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Manual.Instructions" = "Please review and approve the deployment to production."
    "Octopus.Action.Manual.ResponsibleTeamIds" = octopusdeploy_team.approvers.id
  }
}

# Run Azure Script Step
resource "octopusdeploy_process_step" "azure" {
  space_id = module.standard.space
  process_id = octopusdeploy_process.test.id
  name = "Run Azure Script"
  type = "Octopus.AzurePowerShell"
  is_required = true
  environments = [
    module.standard.environments.development,
    module.standard.environments.staging,
    module.standard.environments.production
  ]
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Azure.AccountId" = octopusdeploy_azure_service_principal.azure_account.id
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax" = "PowerShell"
    "Octopus.Action.Script.ScriptBody" = <<-EOT
      # Your Azure PowerShell script here
      Write-Output "Running Azure Script"
    EOT
    "OctopusUseBundledTooling" = "False"
  }
  container = {
    feed_id = module.standard.docker_registry
    image   = "your-private-registry/your-execution-container:latest"
  }
}

# Deploy Kubernetes YAML Step
resource "octopusdeploy_process_step" "kubernetes" {
  space_id = module.standard.space
  process_id = octopusdeploy_process.test.id
  name = "Deploy Kubernetes YAML"
  properties = {
    "Octopus.Action.TargetRoles" = "k8s_cluster"
  }
  type = "Octopus.KubernetesDeployRawYaml"
  is_required = true
  environments = [
    module.standard.environments.development,
    module.standard.environments.staging,
    module.standard.environments.production
  ]
  worker_pool_id = octopusdeploy_static_worker_pool.default.id
  execution_properties = {
    "Octopus.Action.RunOnServer" = "true"
    "Octopus.Action.KubernetesContainers.CustomResourceYamlFileName" = "deployment.yaml"
    "Octopus.Action.KubernetesContainers.Namespace" = "#{Octopus.Environment.Name | ToLower}"
    "Octopus.Action.Script.ScriptSource" = "GitRepository"
    "Octopus.Action.Git.Url" = "https://github.com/your-org/your-repo.git"
    "Octopus.Action.Git.Branch" = "main"
    "OctopusUseBundledTooling" = "False"
  }
  container = {
    feed_id = module.standard.docker_registry
    image   = "octopusdeploy/worker-tools:latest"
  }
  git_dependencies = {
    "": {
      git_credential_id = octopusdeploy_git_credential.git_credential.id
      git_credential_type = "Library"
      repository_uri = "https://github.com/your-org/your-repo.git"
      default_branch = "main"
    }
  }
}

resource "octopusdeploy_process_steps_order" "test" {
  space_id = module.standard.space
  process_id  = octopusdeploy_process.test.id
  steps = [
    octopusdeploy_process_step.approve.id,
    octopusdeploy_process_step.azure.id,
    octopusdeploy_process_step.kubernetes.id,
  ]
}

resource "octopusdeploy_team" "approvers" {
  space_id = module.standard.space
  name        = "Deployment Approvers"
  description = "Team responsible for approving production deployments"
}

resource "octopusdeploy_azure_service_principal" "azure_account" {
  space_id = module.standard.space
  name             = "Example Azure Service Principal"
  description      = "Azure service principal account for deployments"
  subscription_id  = "00000000-0000-0000-0000-000000000001"
  tenant_id        = "00000000-0000-0000-0000-000000000002"
  application_id   = "00000000-0000-0000-0000-000000000003"
  password         = "TopSecretPassword01!"
  environments     = [
    module.standard.environments.development,
    module.standard.environments.staging,
    module.standard.environments.production,
  ]
}

resource "octopusdeploy_static_worker_pool" "default" {
    space_id = module.standard.space
    name = "Worker Pool 1"
    is_default = true
}

resource "octopusdeploy_git_credential" "git_credential" {
  space_id = module.standard.space
  name = "TerraformGitCred"
  username = "TFP"
  password = "password01!"
}

resource "octopusdeploy_kubernetes_cluster_deployment_target" "k8s_cluster_1" {
  space_id = module.standard.space
  name                              = "Example Kubernetes Cluster"
  cluster_url                       = "https://your-cluster-url.com"
  skip_tls_verification             = false
  default_worker_pool_id            = octopusdeploy_static_worker_pool.default.id
  namespace                         = "default"
  environments = [
    module.standard.environments.development,
    module.standard.environments.staging,
    module.standard.environments.production
  ]
  roles = ["k8s-cluster"]
  azure_service_principal_authentication {
    account_id = octopusdeploy_azure_service_principal.azure_account.id
    cluster_name = "Cluster"
    cluster_resource_group = "cluster-rg"
  }
}
```

## Space and Project configuration module
Module with standard configuration to be reused for multiple applications

### variables.tf
```terraform
﻿variable "space" {
  description = "Space details"
  type = object({
    name        = string
    description = optional(string, "")
  })
}

variable "project" {
  description = "Application details"
  type = object({
    name        = string
    description = optional(string, "")
  })
}
```

### main.tf
```terraform
﻿data "octopusdeploy_teams" "everyone" {
  partial_name = "Everyone"
  skip         = 0
  take         = 1
}

# Create a Space
resource "octopusdeploy_space" "example" {
  name        = var.space.name
  description = var.space.description
  space_managers_teams = [ data.octopusdeploy_teams.everyone.teams[0].id ]
}

# Create 3 Environments
resource "octopusdeploy_environment" "development" {
  name        = "Development"
  description = "Development environment"
  space_id    = octopusdeploy_space.example.id
}

resource "octopusdeploy_environment" "staging" {
  name        = "Staging"
  description = "Staging environment"
  space_id    = octopusdeploy_space.example.id
}

resource "octopusdeploy_environment" "production" {
  name        = "Production"
  description = "Production environment"
  space_id    = octopusdeploy_space.example.id
}

# Create a Container Registry
resource "octopusdeploy_docker_container_registry" "example" {
  name        = "DockerHub"
  space_id    = octopusdeploy_space.example.id
  feed_uri    = "https://registry.hub.docker.com"
}

# Additional required resources
resource "octopusdeploy_lifecycle" "example" {
  name     = "Example Lifecycle"
  space_id = octopusdeploy_space.example.id
}

resource "octopusdeploy_project_group" "example" {
  name     = "Example Project Group"
  space_id = octopusdeploy_space.example.id
}

# Create a Project
resource "octopusdeploy_project" "example" {
  name        = var.project.name
  description = var.project.description
  space_id    = octopusdeploy_space.example.id
  lifecycle_id = octopusdeploy_lifecycle.example.id
  project_group_id = octopusdeploy_project_group.example.id
}
```

### output.tf
```terraform
﻿output "space" {
  value = octopusdeploy_space.example.id
}

output "project" {
  value = octopusdeploy_project.example.id
}

output "environments" {
  value = {
    development = octopusdeploy_environment.development.id
    staging = octopusdeploy_environment.staging.id
    production = octopusdeploy_environment.production.id
  }
}

output "docker_registry" {
  value = octopusdeploy_docker_container_registry.example.id
}
```
