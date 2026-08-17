package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetRunRunbookActionResourceSchema returns the nested attribute for a trigger action that runs a runbook.
func GetRunRunbookActionResourceSchema(isRequired bool) resourceSchema.SingleNestedAttribute {
	return resourceSchema.SingleNestedAttribute{
		Description: "Runs a runbook when the trigger fires.",
		Required:    isRequired,
		Optional:    !isRequired,
		Attributes: map[string]resourceSchema.Attribute{
			"runbook_id": GetRequiredStringResourceSchema("The ID of the runbook to run."),
			"target_environment_ids": resourceSchema.ListAttribute{
				Description: "The IDs of the environments to run the runbook in. At least one is required.",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

type RunRunbookActionModel struct {
	RunbookId            types.String `tfsdk:"runbook_id"`
	TargetEnvironmentIds types.List   `tfsdk:"target_environment_ids"`
}
