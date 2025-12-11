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
