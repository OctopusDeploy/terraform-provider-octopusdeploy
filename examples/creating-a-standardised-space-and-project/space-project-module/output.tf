output "space" {
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