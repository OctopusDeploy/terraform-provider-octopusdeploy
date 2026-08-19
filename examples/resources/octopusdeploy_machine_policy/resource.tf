resource "octopusdeploy_machine_policy" "example" {
  name        = "Example Machine Policy (OK to Delete)"
  description = "A machine policy that checks its targets once a week."

  machine_health_check_policy {
    # 7 days, expressed in nanoseconds.
    health_check_interval = 604800000000000
    health_check_type     = "RunScript"

    bash_health_check_policy {
      run_type = "InheritFromDefault"
    }

    powershell_health_check_policy {
      run_type = "InheritFromDefault"
    }
  }
}

# Setting health_check_interval to 0 turns automatic health checks off, matching the "Never"
# schedule in the Octopus UI. Targets using this policy are then only checked on demand.
resource "octopusdeploy_machine_policy" "no_automatic_health_checks" {
  name        = "Manually Checked Targets (OK to Delete)"
  description = "A machine policy that never runs automatic health checks."

  machine_health_check_policy {
    health_check_interval = 0

    bash_health_check_policy {
      run_type = "InheritFromDefault"
    }

    powershell_health_check_policy {
      run_type = "InheritFromDefault"
    }
  }
}
