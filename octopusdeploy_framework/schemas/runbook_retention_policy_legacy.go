package schemas

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var legacyRunbookRetentionPolicySchemeAttributeNames = struct {
	QuantityToKeep    string
	ShouldKeepForever string
}{
	QuantityToKeep:    "quantity_to_keep",
	ShouldKeepForever: "should_keep_forever",
}

func GetLegacyRunbookRetentionPolicyObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		legacyRunbookRetentionPolicySchemeAttributeNames.QuantityToKeep:    types.Int64Type,
		legacyRunbookRetentionPolicySchemeAttributeNames.ShouldKeepForever: types.BoolType,
	}
}

func getLegacyRunbookRetentionPolicySchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		legacyRunbookRetentionPolicySchemeAttributeNames.QuantityToKeep: resourceSchema.Int64Attribute{
			Description: "How many runs to keep per environment.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		legacyRunbookRetentionPolicySchemeAttributeNames.ShouldKeepForever: resourceSchema.BoolAttribute{
			Description: "Indicates if items should never be deleted. The default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(legacyRunbookRetentionPolicySchemeAttributeNames.QuantityToKeep)),
			},
		},
	}
}

func MapFromLegacyRunbookRetentionPolicy(retentionPolicy *runbooks.RunbookRetentionPolicy) attr.Value {
	if retentionPolicy == nil {
		return MapFromLegacyRunbookRetentionPolicy(runbooks.NewDefaultRunbookRetentionPolicy())
	}

	attrs := map[string]attr.Value{
		legacyRunbookRetentionPolicySchemeAttributeNames.QuantityToKeep:    types.Int64Value(int64(retentionPolicy.QuantityToKeep)),
		legacyRunbookRetentionPolicySchemeAttributeNames.ShouldKeepForever: types.BoolValue(retentionPolicy.Strategy == runbooks.RunbookRetentionStrategyForever),
	}

	return types.ObjectValueMust(GetLegacyRunbookRetentionPolicyObjectType(), attrs)
}

func MapToLegacyRunbookRetentionPolicy(flattenedRunbookRetentionPolicy types.List) *runbooks.RunbookRetentionPolicy {
	if flattenedRunbookRetentionPolicy.IsNull() || len(flattenedRunbookRetentionPolicy.Elements()) == 0 {
		return runbooks.NewDefaultRunbookRetentionPolicy()
	}
	obj := flattenedRunbookRetentionPolicy.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	if shouldKeepForever, ok := attrs[legacyRunbookRetentionPolicySchemeAttributeNames.ShouldKeepForever].(types.Bool); ok && !shouldKeepForever.IsNull() {
		if shouldKeepForever.ValueBool() {
			return runbooks.NewKeepForeverRunbookRetentionPolicy()
		}
	}
	if quantityToKeep, ok := attrs[legacyRunbookRetentionPolicySchemeAttributeNames.QuantityToKeep].(types.Int64); ok && !quantityToKeep.IsNull() {
		if quantityToKeep.ValueInt64() == 0 {
			return runbooks.NewKeepForeverRunbookRetentionPolicy()
		}
		if quantityToKeep.ValueInt64() == 100 {
			return runbooks.NewDefaultRunbookRetentionPolicy()
		}
		policy, err := runbooks.NewCountBasedRunbookRetentionPolicy(int32(quantityToKeep.ValueInt64()), runbooks.RunbookRetentionUnitItems)
		if err == nil {
			return policy
		}
		fmt.Printf("Error creating count based runbook retention policy: %v\n", err)
	}
	return runbooks.NewDefaultRunbookRetentionPolicy()
}
