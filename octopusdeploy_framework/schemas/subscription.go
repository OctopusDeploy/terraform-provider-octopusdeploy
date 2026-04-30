package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

	notificationAttributes := map[string]resourceSchema.Attribute{
		"filter": resourceSchema.SingleNestedAttribute{
			Optional:    true,
			Description: "Filter criteria to limit which events trigger this subscription. When omitted, all events trigger the subscription.",
			Attributes:  filterAttributes,
		},
		"email_teams": util.ResourceList(types.StringType).
			Optional().
			Computed().
			Description("Team IDs to notify via email.").
			PlanModifiers(listplanmodifier.UseStateForUnknown()).
			Build(),
		"email_frequency_period": resourceSchema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "How often to send email digests (e.g. '01:00:00' for hourly).",
			Default:     stringdefault.StaticString("01:00:00"),
		},
		"email_priority": resourceSchema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Priority of notification emails. Valid values: Normal, High.",
			Default:     stringdefault.StaticString("Normal"),
			Validators: []validator.String{
				stringvalidator.OneOf("Normal", "High"),
			},
		},
		"email_show_dates_in_timezone_id": resourceSchema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Timezone ID for dates shown in emails (e.g. 'UTC').",
			Default:     stringdefault.StaticString("UTC"),
		},
		"webhook_uri": resourceSchema.StringAttribute{
			Optional:    true,
			Description: "URI to send webhook notifications to.",
		},
		"webhook_teams": util.ResourceList(types.StringType).
			Optional().
			Computed().
			Description("Team IDs to notify via webhook.").
			PlanModifiers(listplanmodifier.UseStateForUnknown()).
			Build(),
		"webhook_timeout": resourceSchema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Timeout for webhook calls (e.g. '00:00:10' for 10 seconds).",
			Default:     stringdefault.StaticString("00:00:10"),
		},
		"webhook_header_key": resourceSchema.StringAttribute{
			Optional:    true,
			Description: "Custom header key to include in webhook requests.",
		},
		"webhook_header_value": resourceSchema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "Custom header value to include in webhook requests.",
		},
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
				Attributes:  notificationAttributes,
			},
		},
	}
}
