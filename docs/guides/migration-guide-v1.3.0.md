---
page_title: "Migrating to v1.3.0"
subcategory: "Upgrades & Migrations"
---

# v1.3.0 Migration Guide
In this release, we've announced a deprecation that will require action from some customers, depending on their configuration

## Deprecated - `octopusdeploy_project_auto_create_release`
In this release, we announced the deprecation of the `octopusdeploy_project_auto_create_relesae` resource in favour of the existing `octopusdeploy_built_in_trigger` resource.

### Rationale
The `octopusdeploy_project_auto_create_release` resource duplicates existing functionality provided by `octopusdeploy_built_in_trigger`.

### Impact
This change requires some customers to update their HCL.
Only customers who have adopted the recently introduced `octopusdeploy_project_auto_create_release` are effected by these changes.

### Timeline
Migration will be required no earlier than 2026-08-08

| Date       | What we'll do                                                            | What you need to do                                                      |
|------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------|
| 2026-02-08 | **Enactment**: Soft-delete the deprecated attribute (Major release)      | Migrate your Terraform config, or use the escape-hatch, before upgrading |
| 2026-08-08 | **Completion**: Remove the deprecated resources entirely (Patch release) | Migrate your Terraform config before upgrading                           |

### How to migrate

Please ensure you are working from a clean slate and have no pending changes to your Terraform config, by running a `terraform plan`. If you have outstanding changes, please resolve them before proceeding with this guide.

-> This migration substitutes equivalent resources. This is non-destructive as long as you complete the migration in one go.

1. Declare a new resource of type `octopusdeploy_built_in_trigger`
1. Set `project_id` to the parent project of your `deployment_process_id` attribute on your existing `octopusdeploy_project_auto_create_release` resource
1. Set all other `octopusdeploy_built_in_trigger` attributes to match the remaining `octopusdeploy_project_auto_create_release` attributes
1. Run `terraform plan`
1. When satisfied, run `terraform apply` to complete the migration

### Escape hatch

We expect customers to migrate their configs in the 6 months between Announcement and Enactment of a deprecation. However, we know that this isn't always possible, so we have a further 6 months grace period.

If you're caught out during this period and need a bit more time to migrate, you can use this escape hatch to revert the soft-deletion from the Enactment stage.

| Environment Variable | Required Value                       |
|----------------------|--------------------------------------|
| `TF_OCTOPUS_DEPRECATION_REVERSALS` | `project-auto-create-release-v1.3.0` |

This escape hatch will be removed and migration will be required during the [Completion phase](#Timeline)
