resource "octopusdeploy_machine_policy" "example" {
  name        = "Example Machine Policy (OK to Delete)"
  description = "A machine policy with custom Tentacle communication timeouts."

  # All timeout values are specified in nanoseconds.
  connection_connect_timeout      = 60000000000 # 1 minute
  connection_retry_count_limit    = 5
  connection_retry_sleep_interval = 1000000000 # 1 second

  # Octopus Server caps these two timeouts at a maximum of 30 minutes (1800000000000 nanoseconds).
  # A larger value is rejected when the policy is created or updated.
  connection_retry_time_limit   = 300000000000 # 5 minutes
  polling_request_queue_timeout = 120000000000 # 2 minutes
}
