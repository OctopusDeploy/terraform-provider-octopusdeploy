package schemas

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func DeprecatedGetResourceRetentionBlockSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		DeprecationMessage: "This block will deprecate when octopus 2025.2 is no longer supported. After upgrading to octopus 2025.3 or higher, please use the `release_retention_with_strategy` and `tentacle_retention_with_strategy` blocks instead.",
		Description:        "Defines the retention policy for releases or tentacles.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"quantity_to_keep": util.ResourceInt64().
					Optional().Computed().
					Default(int64default.StaticInt64(30)).
					Validators(int64validator.AtLeast(0)).
					Description("The number of days/releases to keep. The default value is 30. If 0 then all are kept.").
					Build(),
				"should_keep_forever": util.ResourceBool().
					Optional().Computed().
					Default(booldefault.StaticBool(false)).
					Description("Indicates if items should never be deleted. The default value is false.").
					Build(),
				"unit": util.ResourceString().
					Optional().Computed().
					Default(stringdefault.StaticString("Days")).
					Validators(stringvalidator.OneOfCaseInsensitive("Days", "Items")).
					Description("The unit of quantity to keep. Valid units are Days or Items. The default value is Days.").
					Build(),
			},
			Validators: []validator.Object{
				deprecatedRetentionValidator{},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func DeprecatedGetRetentionAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"quantity_to_keep":    util.DataSourceInt64().Computed().Description("The quantity of releases to keep.").Build(),
				"should_keep_forever": util.DataSourceBool().Computed().Description("Whether releases should be kept forever.").Build(),
				"unit":                util.DataSourceString().Computed().Description("The unit of time for the retention policy.").Build(),
			},
		},
	}
}

type deprecatedRetentionValidator struct{}

func (v deprecatedRetentionValidator) Description(ctx context.Context) string {
	return "validates that should_keep_forever is true only if quantity_to_keep is 0"
}

func (v deprecatedRetentionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v deprecatedRetentionValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var retentionPolicy struct {
		QuantityToKeep    types.Int64  `tfsdk:"quantity_to_keep"`
		ShouldKeepForever types.Bool   `tfsdk:"should_keep_forever"`
		Unit              types.String `tfsdk:"unit"`
	}

	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &retentionPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !retentionPolicy.QuantityToKeep.IsNull() && !retentionPolicy.QuantityToKeep.IsUnknown() && !retentionPolicy.ShouldKeepForever.IsNull() && !retentionPolicy.ShouldKeepForever.IsUnknown() {
		quantityToKeep := retentionPolicy.QuantityToKeep.ValueInt64()
		shouldKeepForever := retentionPolicy.ShouldKeepForever.ValueBool()

		if quantityToKeep == 0 && !shouldKeepForever {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("should_keep_forever"),
				"Invalid retention policy configuration",
				"should_keep_forever must be true when quantity_to_keep is 0",
			)
		} else if quantityToKeep != 0 && shouldKeepForever {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("should_keep_forever"),
				"Invalid retention policy configuration",
				"should_keep_forever must be false when quantity_to_keep is not 0",
			)
		}
	}

	if !retentionPolicy.Unit.IsNull() && !retentionPolicy.Unit.IsUnknown() {
		unit := retentionPolicy.Unit.ValueString()
		if !strings.EqualFold(unit, "Days") && !strings.EqualFold(unit, "Items") {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("unit"),
				"Invalid retention policy unit",
				"Unit must be either 'Days' or 'Items' (case insensitive)",
			)
		}
	}
}
