resource "octopusdeploy_channel" "example" {
  name       = "Development Channel (OK to Delete)"
  project_id = "Projects-123"
}

# A channel whose package version rule orders by publish date instead of
# SemVer ("Most Recently Published"). Use this for packages versioned with
# non-SemVer schemes such as date stamps or feature-branch tags. The
# version_tag_regex matches the full version string in place of the SemVer
# version_range / tag filtering.
#
# Requires the `non-semver-ordering` feature toggle on the Octopus instance;
# without it the server silently ignores the MostRecentlyPublished strategy.
resource "octopusdeploy_channel" "most_recently_published" {
  name       = "Most Recently Published (OK to Delete)"
  project_id = "Projects-123"

  rule {
    versioning_strategy = "MostRecentlyPublished"
    version_tag_regex   = ".*"

    action_package {
      deployment_action = "Reference NuGet package"
      package_reference = "nuget-pkg"
    }
  }
}
