resource "octopusdeploy_azure_service_fabric_cluster_deployment_target" "example" {
  account_id                        = "Accounts-123"
  connection_endpoint               = "[connection-endpoint]"
  environments                      = ["Environments-123", "Environment-321"]
  name                              = "Azure Service Fabric Cluster Deployment Target (OK to Delete)"
  storage_account_name              = "[storage_account_name]"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"
}
