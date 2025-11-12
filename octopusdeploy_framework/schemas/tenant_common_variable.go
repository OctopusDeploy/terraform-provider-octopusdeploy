package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const (
	TenantCommonVariableResourceDescription = "Tenant Common Variable"
	TenantCommonVariableResourceName        = "tenant_common_variable"
)

func GetTenantCommonVariableResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a tenant common variable in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":                      GetIdResourceSchema(),
			"space_id":                GetSpaceIdResourceSchema(TenantCommonVariableResourceDescription),
			"tenant_id":               GetRequiredStringResourceSchema("The ID of the tenant."),
			"library_variable_set_id": GetRequiredStringResourceSchema("The ID of the library variable set."),
			"template_id":             GetRequiredStringResourceSchema("The ID of the variable template."),
			"value": schema.StringAttribute{
				Optional:    true,
				Description: "The value of the variable.",
				Sensitive:   true,
			},
		},
		Blocks: map[string]schema.Block{
			"scope": schema.ListNestedBlock{
				Description: "Sets the scope of the variable.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"environment_ids": getEnvironmentsResourceSchema("A set of environment IDs to scope this variable to."),
					},
				},
			},
		},
	}
}
