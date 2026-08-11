package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebhookTriggerSchema struct{}

var _ EntitySchema = WebhookTriggerSchema{}

func (w WebhookTriggerSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages webhook triggers in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema("webhook trigger"),
			"space_id":    GetSpaceIdResourceSchema("webhook trigger"),
			"project_id": util.ResourceString().
				Description("The ID of the project to attach the trigger.").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Required().
				Build(),
			"is_disabled": GetOptionalBooleanResourceAttribute("Disables the trigger from being run when set.", false),
			"secret": GetSensitiveResourceSchema(
				"The secret used to authenticate incoming webhook requests, sent in the 'X-Octopus-Webhook-Secret' header. Required unless 'require_api_key' is set.", false),
			"require_api_key": GetOptionalBooleanResourceAttribute(
				"Authenticate using an Octopus API key, sent in the 'X-Octopus-ApiKey' header.", false),
			"webhook_id": resourceSchema.StringAttribute{
				Description: "The server generated identifier used in the webhook URL that invokes this trigger.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_ids": resourceSchema.ListAttribute{
				Description: "The IDs of the tenants trigger should apply to.",
				ElementType: types.StringType,
				Optional:    true,
			},
			// Required - this trigger only allows RunRunbook action at this stage.
			"run_runbook_action": GetRunRunbookActionResourceSchema(true),
		},
	}
}

func (w WebhookTriggerSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type WebhookTriggerResourceModel struct {
	Name             types.String           `tfsdk:"name"`
	Description      types.String           `tfsdk:"description"`
	SpaceId          types.String           `tfsdk:"space_id"`
	ProjectId        types.String           `tfsdk:"project_id"`
	IsDisabled       types.Bool             `tfsdk:"is_disabled"`
	Secret           types.String           `tfsdk:"secret"`
	RequireApiKey    types.Bool             `tfsdk:"require_api_key"`
	WebhookId        types.String           `tfsdk:"webhook_id"`
	TenantIds        types.List             `tfsdk:"tenant_ids"`
	RunRunbookAction *RunRunbookActionModel `tfsdk:"run_runbook_action"`

	ResourceModel
}
