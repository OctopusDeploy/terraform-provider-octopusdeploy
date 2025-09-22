terraform {
  required_providers {
    octopusdeploy = {
      source  = "octopus.com/com/octopusdeploy"
      version = "1.2.100"
    }
  }
}
provider "octopusdeploy" {
  # configuration options
  address = "http://localhost:8066/" # (required; string) the service endpoint of the Octopus REST API
  # // latest
  # api_key = "API-ZT38SPCVL2O6O6TCVNVFTAPA4HWYCOE5"
  # space_id = "Spaces-383"

  // 2025.2
  space_id = "Spaces-2"
  api_key  = "API-ZCIVPVCB7B4FG3DDOWNJTNUUBYVUOYQ"
}
#
# resource "octopusdeploy_lifecycle" "lifecycle-1" {
#   name = "Lifecycle-1"
# }

resource "octopusdeploy_lifecycle" "lifecycle-2" {
  name = "Lifecycle-2"
  release_retention_policy {
    quantity_to_keep = 0
  }
}
# data "octopusdeploy_lifecycles" "example" {
#   ids = [octopusdeploy_lifecycle.lifecycle-9991.id]
# }
#
# output "trying_out_lifecycle_dot_notation" {
#   value = data.octopusdeploy_lifecycles.example.lifecycles[0].release_retention_policy[0].should_keep_forever
# }
# output "retention_policy" {
#   value = data.octopusdeploy_lifecycles.example.lifecycles[0].phase[0].release_retention_policy[0].should_keep_forever
# }


# resource "octopusdeploy_lifecycle" "lifecycle-2" {
#   name = "Lifecycle-2"
#   release_retention_policy {
#     quantity_to_keep = "1"
#     unit             = "Days"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-3" {
#   name = "Lifecycle-3"
#   release_retention_policy {
#     should_keep_forever = "true"
#     quantity_to_keep    = "0"
#     unit                = "Items"
#   }
# }

# resource "octopusdeploy_lifecycle" "lifecycle-4" {
#   name = "Lifecycle-4"
#   release_retention_policy {
#     quantity_to_keep = "0"
#     unit             = "Items"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-5" {
#   name = "Lifecycle-5"
#   release_retention_policy {
#     should_keep_forever = "true"
#     unit                = "Items"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-6" {
#   name = "Lifecycle-6"
#   release_retention_policy {
#     should_keep_forever = "true"
#     quantity_to_keep    = "0"
#   }
# }

#   resource "octopusdeploy_lifecycle" "lifecycle-7" {
#     name = "Lifecycle-7"
#     release_retention_policy {
#       quantity_to_keep = "0"
#     }
#   }
#   resource "octopusdeploy_lifecycle" "lifecycle-8" {
#     name = "Lifecycle-8"
#     release_retention_policy {
#       should_keep_forever = "true"
#     }
#   }

# resource "octopusdeploy_lifecycle" "lifecycle-9" {
#   name = "Lifecycle-9"
#   release_retention_policy {
#     should_keep_forever = "true"
#     quantity_to_keep    = "1"
#     unit                = "Items"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-10" {
#   name = "Lifecycle-10"
#   release_retention_policy {
#
#     should_keep_forever = "false"
#     quantity_to_keep    = "0"
#     unit                = "Items"
#
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-11" {
#   name = "Lifecycle-11"
#   release_retention_policy {
#
#     should_keep_forever = "false"
#     unit                = "Items"
#
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-12" {
#   name = "Lifecycle-12"
#   release_retention_policy {
#     unit = "Days"
#   }
# }

# resource "octopusdeploy_lifecycle" "lifecycle-13" {
#   name = "Lifecycle-13"
#   release_retention_policy {
#     should_keep_forever = "true"
#     quantity_to_keep    = "1"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-14" {
#   name = "Lifecycle-14"
#   release_retention_policy {
#     should_keep_forever = "false"
#     quantity_to_keep    = "1"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-15" {
#   name = "Lifecycle-15"
#   release_retention_policy {
#     should_keep_forever = "false"
#     quantity_to_keep    = "0"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-16" {
#   name = "Lifecycle-16"
#   release_retention_policy {
#     should_keep_forever = "false"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-17" {
#   name = "Lifecycle-17"
#   release_retention_policy {
#     quantity_to_keep = "1"
#   }
# }
# resource "octopusdeploy_lifecycle" "lifecycle-18" {
#   name = "Lifecycle-18"
#   release_retention_policy {
#
#   }
# }