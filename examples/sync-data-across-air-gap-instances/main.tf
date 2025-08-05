provider "octopusdeploy" {
  address  = "https://test.prod-octopus.app"
  api_key  = var.prod_api_key
  alias = "Production"
  space_id = "Spaces-1"
}

provider "octopusdeploy" {
  address  = "https://test.staging-octopus.app"
  api_key  = var.staging_api_key
  alias = "Staging"
  space_id = "Spaces-1"
}
# variable for AWS Account
variable "aws_account_secret" {
  type        = string
  sensitive   = true
  description = "The AWS secret key associated with the account"
  default     = "Change Me!"
}

variable "feed_uri"{
  type = string
  sensitive = true
  description = "The Feed URI" 
  default = "default_URL"
}

variable "certificate_data"{
  type = string
  sensitive = true
  description = "The base 64 encoded string for certificate"
  default = "default_data"
}

# Read from Production instance 
data "octopusdeploy_projects" "Production"{
  space_id = "Spaces-1"
  provider = octopusdeploy.Production
}

#AWS Account
data "octopusdeploy_accounts" "Production" {
  space_id = "Spaces-1"
  provider = octopusdeploy.Production
}

# Environments
data "octopusdeploy_environments" "Production"{
  space_id = "Spaces-1"
  provider = octopusdeploy.Production
}

# Certificates
data "octopusdeploy_certificates" "Production"{
  space_id = "Spaces-1"
  partial_name = "Terraform-test"
  provider = octopusdeploy.Production
}

# Variables
data "octopusdeploy_library_variable_sets" "Production"{
  space_id = "Spaces-1"
  provider = octopusdeploy.Production
}

# External Feeds
data "octopusdeploy_feeds" "Production" {
  space_id = "Spaces-1"
  partial_name = "test"
  provider = octopusdeploy.Production
}

# Create project from datasource
resource "octopusdeploy_project" "Staging" {
  lifecycle_id = "Lifecycles-1"
  name = data.octopusdeploy_projects.Production.projects[0].name
  provider = octopusdeploy.Staging
  project_group_id = "ProjectGroups-1"
}

# Create AWS Account from other project
resource "octopusdeploy_aws_account" "Staging"{
  provider = octopusdeploy.Staging
  access_key = data.octopusdeploy_accounts.Production.accounts[0].access_key
  name = data.octopusdeploy_accounts.Production.accounts[0].name
  secret_key = var.aws_account_secret
}

# Create Environments from other project
resource "octopusdeploy_environment" "Staging"{
  provider = octopusdeploy.Staging
  for_each = { for env in data.octopusdeploy_environments.Production.environments : env.name => env }
  name        = each.value.name
}

# Create Certificate from other project
resource "octopusdeploy_certificate" "Staging" {
  provider = octopusdeploy.Staging
  certificate_data = "certificate_data_goes_here"
  name = data.octopusdeploy_certificates.Production.certificates[0].name
  self_signed = "true"
  subject_common_name = data.octopusdeploy_certificates.Production.certificates[0].subject_common_name
}