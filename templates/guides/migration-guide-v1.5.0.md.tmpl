---
page_title: "Migrating to v1.6.0"
subcategory: "Upgrades & Migrations"
---

# v1.6.0 Migration Guide

In this release, we've announced a deprecation that will require action from some customers, depending on their configuration

## Deprecated - runbook `retention_policy` blocks

In this release, we announced the deprecation the following runbook resource block:

- `octopusdeploy_runbook.retention_policy`

in favour of the following runbook resource block:

- `octopusdeploy_runbook.retention_policy_with_strategy`

**In addition to this**, the default retention policy for runbooks without explicit retention was previously set to "100" "Items". After switching to `retention_policy_with_strategy`, the default retention will be "Default", for Octopus server version that support default retention strategies.

### Rationale

In octopus 2026.1 and higher, a customizable default retention setting can be added at the space level and applied to desired runbooks within that space. The new `retention_policy_with_strategy` block support this feature. This feature enables retention of many runbooks to be changed on mass or more easily initiated at the user’s desired retention.

Runbooks will now be initially set to the Space Default retention setting (the space default retention setting is initially set to "Keep 100 runs").

### Impact

This change affects all customers using runbooks.

Runbook resources using `retention_policy` retention settings:

- users will need to change their HCL to use the new `retention_policy_with_strategy` block

Runbook resources **without** explicit retention settings:

- if users don't want these runbooks to have the space default retention, they will need to update their HCL to have explicit retention settings

### Timeline

Migration will be required no earlier than 21 Oct 2026

| Date        | What we'll do                                                        | What you need to do                                                      |
| ----------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| 05 JUN 2026 | **Enactment**: Soft-delete the deprecated block (Major release)      | Migrate your Terraform config, or use the escape-hatch, before upgrading |
| 05 DEC 2026 | **Completion**: Remove the deprecated block entirely (Patch release) | Migrate your Terraform config before upgrading                           |

### How to migrate

Please ensure you are working from a clean slate and have no pending changes to your Terraform config, by running a `terraform plan`. If you have outstanding changes, please resolve them before proceeding with this guide.

-> This migration substitutes equivalent resources. This is non-destructive as long as you complete the migration in one go.

1.  Within your runbook resources, replace all `retention_policy` blocks with `retention_policy_with_strategy` blocks and equivalent retention settings

    - e.g. when setting keep forever for the runbook retention, use:

          retention_policy_with_strategy  = {
              strategy = “Forever”
          }

2.  Review all runbooks without explicit retention blocks and set explicit retention where needed.
    All runbooks without explicit retention blocks are currently being initialised with a retention of “100 runs”. If the same behaviour is required after migration, please explicitly set the retention strategy as follows:

        retention_policy_with_strategy = {
            strategy = “Count”
            quantity_to_keep = 100
            unit = “Items”
        }

3.  Run `terraform plan`
4.  When satisfied, run `terraform apply` to complete the migration

### Escape hatch

We expect customers to migrate their configs in the 6 months between Announcement and Enactment of a deprecation. However, we know that this isn't always possible, so we have a further 6 months grace period.

If you're caught out during this period and need a bit more time to migrate, you can use this escape hatch to revert the soft-deletion from the Enactment stage.

| Environment Variable               | Required Value                           |
| ---------------------------------- | ---------------------------------------- |
| `TF_OCTOPUS_DEPRECATION_REVERSALS` | `octopusdeploy_runbook.retention_policy` |

This escape hatch will be removed and migration will be required during the [Completion phase](#Timeline)
