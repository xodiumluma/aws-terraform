// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ExpandFrameworkStringSet(ctx context.Context, v basetypes.SetValuable) []*string {
	var output []*string

	panicOnError(Expand(ctx, v, &output))

	return output
}

func ExpandFrameworkStringValueSet(ctx context.Context, v basetypes.SetValuable) Set[string] {
	var output []string

	panicOnError(Expand(ctx, v, &output))

	return output
}

// FlattenFrameworkStringSet converts a slice of string pointers to a framework Set value.
//
// A nil slice is converted to a null Set.
// An empty slice is converted to a null Set.
func FlattenFrameworkStringSet(ctx context.Context, v []*string) types.Set {
	if len(v) == 0 {
		return types.SetNull(types.StringType)
	}

	var output types.Set

	panicOnError(Flatten(ctx, v, &output))

	return output
}

// FlattenFrameworkStringSetLegacy converts a slice of string pointers to a framework Set value.
//
// A nil slice is converted to an empty (non-null) Set.
func FlattenFrameworkStringSetLegacy(_ context.Context, vs []*string) types.Set {
	elems := make([]attr.Value, len(vs))

	for i, v := range vs {
		elems[i] = types.StringValue(aws.ToString(v))
	}

	return types.SetValueMust(types.StringType, elems)
}

// FlattenFrameworkStringValueSet converts a slice of string values to a framework Set value.
//
// A nil slice is converted to a null Set.
// An empty slice is converted to a null Set.
func FlattenFrameworkStringValueSet[T ~string](ctx context.Context, v []T) types.Set {
	if len(v) == 0 {
		return types.SetNull(types.StringType)
	}

	var output types.Set

	panicOnError(Flatten(ctx, v, &output))

	return output
}

// FlattenFrameworkStringValueSetLegacy is the Plugin Framework variant of FlattenStringValueSet.
// A nil slice is converted to an empty (non-null) Set.
func FlattenFrameworkStringValueSetLegacy[T ~string](_ context.Context, vs []T) types.Set {
	elems := make([]attr.Value, len(vs))

	for i, v := range vs {
		elems[i] = types.StringValue(string(v))
	}

	return types.SetValueMust(types.StringType, elems)
}
