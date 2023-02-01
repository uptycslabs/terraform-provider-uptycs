package uptycs

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func makeListStringAttribute(in []string) types.List {
	values := make([]attr.Value, len(in))
	for i, v := range in {
		values[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, values)
}

func makeListStringAttributeFn[T any](in []T, fn func(T) (string, bool)) types.List {
	values := make([]attr.Value, 0)
	for _, v := range in {
		if s, ok := fn(v); ok {
			values = append(values, types.StringValue(s))
		}
	}
	return types.ListValueMust(types.StringType, values)
}
