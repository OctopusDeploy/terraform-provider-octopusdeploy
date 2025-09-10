resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle (OK to Delete)"

  release_retention_policy {
    strategy = "Forever"
  }

  tentacle_retention_policy {
    strategy         = "Count"
    quantity_to_keep = 30
    unit             = "Items"
  }

  phase {
    automatic_deployment_targets = ["Environments-321"]
    name = "Test Phase 1"

    release_retention_policy {
      quantity_to_keep    = 1
      should_keep_forever = false
      unit                = "Days"
    }

    tentacle_retention_policy {
      strategy = "Default"
    }
  }

  phase {
    is_optional_phase = true
    name              = "Test Phase 1"
    optional_deployment_targets = ["Environments-321"]
  }
}
