package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ObjectValue builds an object from the supplied attribute types and values, filling
// any attribute that the caller did not supply with a null of its declared type.
//
// types.ObjectValueMust panics when an attribute declared in the object type is absent
// from the value map, and that panic reaches users as a provider crash during plan or
// apply rather than as a diagnostic. A mapper reaches this state when a value only
// applies to some variants of a resource and the branch that handles the other variants
// skips the key, which is what caused the crash in issue #166.
//
// Attributes present in attrs but absent from attrTypes are left in place so that
// ObjectValueMust still rejects them. Those indicate a genuine mismatch between a
// mapper and its object type, rather than an attribute that does not apply.
func ObjectValue(attrTypes map[string]attr.Type, attrs map[string]attr.Value) basetypes.ObjectValue {
	complete := make(map[string]attr.Value, len(attrTypes))
	for name, value := range attrs {
		complete[name] = value
	}

	for name, attrType := range attrTypes {
		if _, ok := complete[name]; !ok {
			complete[name] = NullValueOf(attrType)
		}
	}

	return types.ObjectValueMust(attrTypes, complete)
}

// NullValueOf returns a null value of the supplied type. It handles nested collection
// and object types, so callers do not have to name the concrete null constructor.
func NullValueOf(attrType attr.Type) attr.Value {
	ctx := context.Background()

	value, err := attrType.ValueFromTerraform(ctx, tftypes.NewValue(attrType.TerraformType(ctx), nil))
	if err != nil {
		// The framework's own types never fail to produce a null value. A custom type
		// that does is a programming error, and silently substituting something else
		// would hide it.
		panic(fmt.Sprintf("unable to build a null value for type %s: %s", attrType, err))
	}

	return value
}
