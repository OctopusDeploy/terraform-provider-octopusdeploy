# every runbook in the space
data "octopusdeploy_runbooks" "all" {}

# every runbook belonging to one project
data "octopusdeploy_runbooks" "project" {
  project_id = "Projects-123"
}

# runbooks whose name contains "Restart", within one project
data "octopusdeploy_runbooks" "restart" {
  project_id   = "Projects-123"
  partial_name = "Restart"
}

# specific runbooks by ID
data "octopusdeploy_runbooks" "by_id" {
  ids = ["Runbooks-1", "Runbooks-2"]
}

output "runbook_names" {
  value = [for runbook in data.octopusdeploy_runbooks.project.runbooks : runbook.name]
}
