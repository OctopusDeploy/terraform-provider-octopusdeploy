---
page_title: "octopusdeploy_process Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages Runbook and Deployment Processes in Octopus Deploy. It's used in collaboration with octopusdeploy_process_step and octopusdeploy_process_step_order.
---

# octopusdeploy_process (Resource)

This resource manages Runbook and Deployment Processes in Octopus Deploy. It's used in collaboration with `octopusdeploy_process_step` and `octopusdeploy_process_step_order`.

~> This resource is the successor to the original `octopusdeploy_deployment_process` resource, which suffered from numerous problems including: state drift, data inconsistency when reordering or inserting steps, and lack of awareness of Version-Controlled projects.

### Remarks

The `octopusdeploy_process` resource is used in conjunction with a series of other building-block resources to form a full process. They are deliberately designed with dependencies between them so that the deployment process will be incrementally "built up" in Octopus with a series of incremental updates. You can use only the building blocks you need for your process (i.e. if your process doesn't involve Child Steps, you don't need to deal with the `octopusdeploy_process_child_step` resource.)

At a minimum, to get a functional deployment process, you will need:

1. A Project (`octopusdeploy_project`)
1. A Deployment Process referencing the Project (`octopusdeploy_process`)
1. One or more Steps referencing the Deployment Process (`octopusdeploy_process_step`)

The `octopusdeploy_process_step_order` resource isn't strictly required, but it's highly recommended. If you need to change the order of steps in your process, or insert a new step within an existing process, you'll need the Step Order defined first.

Without a defined Step Order, the Steps will be added to the process in the order they're applied by Terraform. This is usually the order they appear in your HCL, but may not always be deterministic. 

## Example Usage
~> See the docs for `octopusdeploy_process_step`, `octopusdeploy_process_steps_order`, `octopusdeploy_process_child_step` and `octopusdeploy_process_child_steps_order` for more detailed examples.

### Deployment Process
```terraform
# Example of a Deployment Process with three steps and an explicit Step Order
resource "octopusdeploy_process" "example" {
  space_id = "Spaces-1"
  project_id  = "Projects-21"
}

resource "octopusdeploy_process_step" "run_script" {
  # Run script step
  process_id  = octopusdeploy_process.example.id  
  name = "Run My Script"
  properties = {
    "Octopus.Action.MaxParallelism" = "2"
    "Octopus.Action.TargetRoles" = "role-1,role-2"
  }
  type = "Octopus.Script"
  environments = [octopusdeploy_environment.development.id]
  excluded_environments = [octopusdeploy_environment.production.id]
  channels = [octopusdeploy_channel.example.id]
  notes = "Script example"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody" = <<-EOT
      Write-Host "Executing step..."
    EOT
  }
}

resource "octopusdeploy_process_step" "approval" {
  # Manual intervention
  process_id  = octopusdeploy_process.example.id
  name = "Approve deployment"
  type = "Octopus.Manual"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Manual.Instructions" = "Example of manual blocking step"
    "Octopus.Action.Manual.BlockConcurrentDeployments" = "True"
    "Octopus.Action.Manual.ResponsibleTeamIds" = "teams-managers"
  }
}

resource "octopusdeploy_process_step" "deploy_package" {
  # Package deployment with primary package
  process_id  = octopusdeploy_process.example.id
  name = "Package deployment"
  properties = {
    "Octopus.Action.TargetRoles" = "role-one"
  }
  type = "Octopus.TentaclePackage"
  packages = {
    "": {
      package_id: "my.package"
      feed_id: "Feeds-1"
    }
  }
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"    
    # Reference primary package in execution properties for legacy purposes
    "Octopus.Action.Package.DownloadOnTentacle" = "False"
    "Octopus.Action.Package.FeedId" = "Feeds-1"
    "Octopus.Action.Package.PackageId" = "my.package"
  }
}

resource "octopusdeploy_process_steps_order" "example" {
  process_id  = octopusdeploy_process.example.id
  steps = [
    octopusdeploy_process_step.run_script.id,
    octopusdeploy_process_step.approval.id,
    octopusdeploy_process_step.deploy_package.id,
  ]
}
```

### Runbook Process
```terraform
# Example of a Runbook Process with two steps and an explicit Step Order
# To manage a Runbook process, specify both the Project and Runbook IDs (usually via Terraform resource references)
resource "octopusdeploy_process" "example" {
  space_id = "Spaces-1"
  project_id  = "Projects-21"
  runbook_id  = "Runbooks-42"
}

resource "octopusdeploy_process_step" "run_script" {
  # Run script step
  process_id  = octopusdeploy_process.example.id  
  name = "Run My Script"
  properties = {
    "Octopus.Action.MaxParallelism" = "2"
    "Octopus.Action.TargetRoles" = "role-1,role-2"
  }
  type = "Octopus.Script"
  environments = [octopusdeploy_environment.development.id]
  excluded_environments = [octopusdeploy_environment.production.id]
  channels = [octopusdeploy_channel.example.id]
  notes = "Script example"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody" = <<-EOT
      Write-Host "Executing step..."
    EOT
  }
}

resource "octopusdeploy_process_step" "deploy_package" {
  # Package deployment with primary package
  process_id  = octopusdeploy_process.example.id
  name = "Package deployment"
  properties = {
    "Octopus.Action.TargetRoles" = "role-one"
  }
  type = "Octopus.TentaclePackage"
  packages = {
    "": {
      package_id: "my.package"
      feed_id: "Feeds-1"
    }
  }
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"    
    # Reference primary package in execution properties for legacy purposes
    "Octopus.Action.Package.DownloadOnTentacle" = "False"
    "Octopus.Action.Package.FeedId" = "Feeds-1"
    "Octopus.Action.Package.PackageId" = "my.package"
  }
}

resource "octopusdeploy_process_steps_order" "example" {
  process_id  = octopusdeploy_process.example.id
  steps = [
    octopusdeploy_process_step.run_script.id,
    octopusdeploy_process_step.deploy_package.id,
  ]
}
```

### Using Process Templates
Process Templates can be consumed in your process using `octopusdeploy_process_step` with `type = "Octopus.ProcessTemplate"`. Process template parameters are configured through `execution_properties`.

Process template specific properties of `"Octopus.Action.ProcessTemplate.Reference.Slug"` and `"Octopus.Action.ProcessTemplate.Reference.VersionMask"` are required to be set within `execution_properties` to specify which process template and version to use.

```terraform
# Example of how to configure a Process Template Usage Step

resource "octopusdeploy_project" "example" {
  project_group_id = "ProjectGroups-1"
  lifecycle_id     = "Lifecycles-1"
  name             = "Example Project"
}

resource "octopusdeploy_process" "example" {
  project_id = octopusdeploy_project.example.id
}

# Process template usage step
resource "octopusdeploy_process_step" "deploy_microservice" {
  process_id = octopusdeploy_process.example.id
  name       = "Run process template - microservice template"
  type       = "Octopus.ProcessTemplate"

  # Package references required by the process template
  packages = {
    "db-package" = {
      package_id           = "myservice-db"
      feed_id              = "Feeds-1"
      acquisition_location = ""
      properties = {
        "PackageParameterName" = "db-package"
        "SelectionMode"        = "deferred"
      }
    }
  }

  execution_properties = {
    # Process template reference - required
    "Octopus.Action.ProcessTemplate.Reference.Slug"        = "deploy-process-microservice-template"
    "Octopus.Action.ProcessTemplate.Reference.VersionMask" = "1.X"

    # Process template parameters
    "k8s-cluster"    = "k8s-deployment-target"
    "email-subject"  = "Deployment succeeded"
    "email-body"     = "The deployment of the microservice was successful."
    "worker-pool"     = "WorkerPools-1"
    "db-connection-string" = "#{Project.DatabaseConnectionString}"

    # Package parameters require JSON format
    "db-package" = "{\"PackageId\":\"myservice-db\",\"FeedId\":\"Feeds-1\"}"
  }
}
```

#### Required Parameter Validation
~> When a new required parameter is added to a process template, any modifications to the process will fail validation until that parameter is configured - even if you're only modifying unrelated steps.

This typically occurs when your process template usage is set to auto-accept minor or patch version updates and a new required parameter gets added. While this goes against recommended practices (adding required parameters should be a major version change), you may encounter validation errors like:

```
╷
│ Error: unable to update process step
│
│   with octopusdeploy_process_step.deploy_microservice,
│   on example.tf line 10, in resource "octopusdeploy_process_step" "deploy_microservice":
│   10: resource "octopusdeploy_process_step" "deploy_microservice" {
│
│ Octopus API error: There was a problem with your request. [Parameter 'NewRequiredParam' in step 'Run process template - microservice template' is required]
╵
```

To resolve this error, you'll need to add the missing required parameter to your process template step's `execution_properties`, even if you weren't intending to modify that step.

For more information on process template versioning best practices, see the [Octopus Deploy Process Templates documentation](https://octopus.com/docs/platform-hub/process-templates/best-practices).

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Id of the project this process belongs to.

### Optional

- `runbook_id` (String) Id of the runbook this process belongs to. When not set this resource represents deployment process of the project
- `space_id` (String) The space ID associated with this process.

### Read-Only

- `id` (String) The unique ID for this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import [options] octopusdeploy_process.<name> <process-id>
```
