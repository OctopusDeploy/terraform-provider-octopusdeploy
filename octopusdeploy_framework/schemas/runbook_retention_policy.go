package schemas

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var runbookRetentionPolicySchemeAttributeNames = struct {
	QuantityToKeep string
	Strategy       string
	Unit           string
}{
	QuantityToKeep: "quantity_to_keep",
	Strategy:       "strategy",
	Unit:           "unit",
}

func GetRunbookRetentionPolicyObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		runbookRetentionPolicySchemeAttributeNames.QuantityToKeep: types.Int64Type,
		runbookRetentionPolicySchemeAttributeNames.Strategy:       types.StringType,
		runbookRetentionPolicySchemeAttributeNames.Unit:           types.StringType,
	}
}

func getRunbookRetentionPolicySchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		runbookRetentionPolicySchemeAttributeNames.Strategy: resourceSchema.StringAttribute{
			Description: "How retention will be set. Valid strategies are `Default`, `Forever` and `Count`.",
			Required:    true,
		},
		runbookRetentionPolicySchemeAttributeNames.QuantityToKeep: resourceSchema.Int64Attribute{
			Description: "The number of runs or days of runs to keep, depending on the unit selected. Required when strategy is `Count`.",
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				NewStrategyAttributeValidator("Count"),
			},
		},
		runbookRetentionPolicySchemeAttributeNames.Unit: resourceSchema.StringAttribute{
			Description: "The unit of the quantity to keep. Valid units are `Items` and `Days`. Required when strategy is `Count`.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("Days", "Items"),
				NewStrategyAttributeValidator("Count"),
			},
		},
	}
}

func MapFromRunbookRetentionPolicy(retentionPolicy *runbooks.RunbookRetentionPolicy) attr.Value {
	if retentionPolicy == nil {
		return MapFromRunbookRetentionPolicy(runbooks.NewDefaultRunbookRetentionPolicy())
	}

	attrs := map[string]attr.Value{
		runbookRetentionPolicySchemeAttributeNames.Strategy: types.StringValue(string(retentionPolicy.Strategy)),
		runbookRetentionPolicySchemeAttributeNames.QuantityToKeep: util.Ternary(retentionPolicy.Strategy == runbooks.RunbookRetentionStrategyCount,
			types.Int64Value(int64(retentionPolicy.QuantityToKeep)),
			types.Int64Null()),
		runbookRetentionPolicySchemeAttributeNames.Unit: util.Ternary(retentionPolicy.Strategy == runbooks.RunbookRetentionStrategyCount,
			types.StringValue(string(retentionPolicy.Unit)),
			types.StringNull()),
	}

	return types.ObjectValueMust(GetRunbookRetentionPolicyObjectType(), attrs)
}

func MapToRunbookRetentionPolicy(policy types.List) (*runbooks.RunbookRetentionPolicy, error) {
	if policy.IsNull() || len(policy.Elements()) == 0 {
		return runbooks.NewDefaultRunbookRetentionPolicy(), nil
	}
	obj := policy.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	strategy, ok := attrs[runbookRetentionPolicySchemeAttributeNames.Strategy].(types.String)
	if !ok {
		return nil, fmt.Errorf("Runbook retention policy strategy is required")
	}

	switch strategy.ValueString() {
	case string(runbooks.RunbookRetentionStrategyDefault):
		return runbooks.NewDefaultRunbookRetentionPolicy(), nil
	case string(runbooks.RunbookRetentionStrategyForever):
		return runbooks.NewKeepForeverRunbookRetentionPolicy(), nil
	case string(runbooks.RunbookRetentionStrategyCount):
		quantityAttr, ok := attrs[runbookRetentionPolicySchemeAttributeNames.QuantityToKeep].(types.Int64)
		if !ok || quantityAttr.IsNull() {
			return nil, fmt.Errorf("Runbook retention policy quantity to keep is required when strategy is 'Count'")
		}
		unitAttr, ok := attrs[runbookRetentionPolicySchemeAttributeNames.Unit].(types.String)
		if !ok || unitAttr.IsNull() {
			return nil, fmt.Errorf("Runbook retention policy unit is required when strategy is 'Count'")
		}
		switch unitAttr.ValueString() {
		case string(runbooks.RunbookRetentionUnitItems):
			return runbooks.NewCountBasedRunbookRetentionPolicy(int32(quantityAttr.ValueInt64()), runbooks.RunbookRetentionUnitItems)
		case string(runbooks.RunbookRetentionUnitDays):
			return runbooks.NewCountBasedRunbookRetentionPolicy(int32(quantityAttr.ValueInt64()), runbooks.RunbookRetentionUnitDays)
		default:
			return nil, fmt.Errorf("Invalid runbook retention policy unit: %s", unitAttr.ValueString())
		}
	default:
		return nil, fmt.Errorf("Invalid runbook retention policy strategy: %s", strategy.ValueString())
	}
}
