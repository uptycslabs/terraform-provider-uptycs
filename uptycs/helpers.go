package uptycs

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/constraints"
)

type JSONUnpackError struct{}

func (m *JSONUnpackError) Error() string {
	return "key not found"
}

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

func getKeyValueFromRawJSON(in string, key string) (string, string, error) {

	_raw := make(map[string]json.RawMessage)
	err := json.Unmarshal([]byte(in), &_raw)
	if err != nil {
		panic(err)
	}

	for k, v := range _raw {
		if k == key {
			return k, string(v)[1 : len(v)-1], nil
		}
	}
	return in, "", &JSONUnpackError{}
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func interSection[T constraints.Ordered](pS ...[]T) []T {
	hash := make(map[T]*int) // value, counter
	result := make([]T, 0)
	for _, slice := range pS {
		duplicationHash := make(map[T]bool) // duplication checking for individual slice
		for _, value := range slice {
			if _, isDup := duplicationHash[value]; !isDup { // is not duplicated in slice
				if counter := hash[value]; counter != nil { // is found in hash counter map
					if *counter++; *counter >= len(pS) { // is found in every slice
						result = append(result, value)
					}
				} else { // not found in hash counter map
					i := 1
					hash[value] = &i
				}
				duplicationHash[value] = true
			}
		}
	}
	return result
}
