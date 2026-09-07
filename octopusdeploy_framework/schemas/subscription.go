package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const SubscriptionResourceName = "subscription"

type SubscriptionSchema struct{}

func (s SubscriptionSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = SubscriptionSchema{}

func (s SubscriptionSchema) GetResourceSchema() resourceSchema.Schema {
	filterList := func(description string) resourceSchema.Attribute {
		return util.ResourceList(types.StringType).
			Optional().
			Computed().
			Description(description).
			PlanModifiers(listplanmodifier.UseStateForUnknown()).
			Build()
	}

	filterAttributes := map[string]resourceSchema.Attribute{
		"users":            filterList("Filter by user IDs."),
		"projects":         filterList("Filter by project IDs."),
		"project_groups":   filterList("Filter by project group IDs."),
		"environments":     filterList("Filter by environment IDs."),
		"event_groups":     filterList("Filter by event groups."),
		"event_categories": filterList("Filter by event categories (e.g. Created, Modified, Deleted)."),
		"event_agents":     filterList("Filter by event agents."),
		"tenants":          filterList("Filter by tenant IDs."),
		"tags":             filterList("Filter by tenant tags."),
		"document_types":   filterList("Filter by document types (e.g. Machines, Projects, Deployments)."),
	}

	return resourceSchema.Schema{
		Description: "This resource manages event notification subscriptions in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema("subscription"),
			"name":     GetNameResourceSchema(true),
			"is_disabled": resourceSchema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this subscription is disabled.",
				Default:     booldefault.StaticBool(false),
			},
			"event_notification_subscription": resourceSchema.SingleNestedAttribute{
				Required:    true,
				Description: "Event notification configuration for this subscription.",
				Attributes: map[string]resourceSchema.Attribute{
					"filter": resourceSchema.SingleNestedAttribute{
						Optional:    true,
						Description: "Filter criteria to limit which events trigger this subscription. When omitted, all events trigger the subscription.",
						Attributes:  filterAttributes,
					},
					"email_teams": util.ResourceSet(types.StringType).
						Optional().
						Computed().
						Description("Team IDs to notify via email.").
						PlanModifiers(setplanmodifier.UseStateForUnknown()).
						Build(),
					"email_frequency_period": util.ResourceString().
						Optional().
						Computed().
						Description("How often to send email digests (e.g. '01:00:00' for hourly).").
						Default("01:00:00").
						Build(),
					"email_priority": util.ResourceString().
						Optional().
						Computed().
						Description("Priority of notification emails. Valid values: Normal, High, Low.").
						Default("Normal").
						Validators(stringvalidator.OneOf("Normal", "High", "Low")).
						Build(),
					"email_show_dates_in_timezone_id": util.ResourceString().
						Optional().
						Computed().
						Description("Timezone ID for dates shown in emails (e.g. 'UTC').").
						Default("UTC").
						Build(),
					"slack_channel_ids": util.ResourceList(types.StringType).
						Optional().
						Computed().
						Description("Slack channel IDs to post to.").
						PlanModifiers(listplanmodifier.UseStateForUnknown()).
						Build(),
					"slack_channel_names": util.ResourceList(types.StringType).
						Optional().
						Computed().
						Description("Display names for the channels in slack_channel_ids, in the same order. If a name is omitted, the channel ID is shown instead.").
						PlanModifiers(listplanmodifier.UseStateForUnknown()).
						Build(),
					"slack_frequency_period": util.ResourceString().
						Optional().
						Computed().
						Description("How often to send Slack digests (e.g. '01:00:00' for hourly).").
						Default("01:00:00").
						Build(),
					"slack_digest_format": util.ResourceString().
						Optional().
						Description("Deprecated and ignored. Slack digests always send a summary.").
						Deprecated("slack_digest_format is no longer used and will be removed in a future release.").
						Build(),
					"teams_webhooks": resourceSchema.ListNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Microsoft Teams channels to post to via incoming webhooks.",
						NestedObject: resourceSchema.NestedAttributeObject{
							Attributes: map[string]resourceSchema.Attribute{
								"id": util.ResourceString().
									Required().
									Description("A unique identifier for this Teams channel.").
									Validators(stringvalidator.LengthAtLeast(1)).
									Build(),
								"name": util.ResourceString().
									Required().
									Description("Display name for this Teams channel.").
									Build(),
								"url": util.ResourceString().
									Required().
									Sensitive().
									Description("Incoming webhook URL for this Teams channel.").
									Build(),
							},
						},
					},
					"teams_frequency_period": util.ResourceString().
						Optional().
						Computed().
						Description("How often to send Teams digests (e.g. '01:00:00' for hourly).").
						Default("01:00:00").
						Build(),
					"webhook_uri": util.ResourceString().
						Optional().
						Description("URI to send webhook notifications to.").
						Build(),
					"webhook_teams": util.ResourceSet(types.StringType).
						Optional().
						Computed().
						Description("Team IDs to notify via webhook.").
						PlanModifiers(setplanmodifier.UseStateForUnknown()).
						Build(),
					"webhook_timeout": util.ResourceString().
						Optional().
						Computed().
						Description("Timeout for webhook calls (e.g. '00:00:10' for 10 seconds).").
						Default("00:00:10").
						Build(),
					"webhook_header_key": util.ResourceString().
						Optional().
						Description("Custom header key to include in webhook requests.").
						Build(),
					"webhook_header_value": util.ResourceString().
						Optional().
						Sensitive().
						Description("Custom header value to include in webhook requests.").
						Build(),
				},
			},
		},
	}
}
