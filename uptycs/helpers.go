package uptycs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func stringDefault(defaultValue string) stringDefaultModifier {
	return stringDefaultModifier{
		Default: defaultValue,
	}
}

// stringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type stringDefaultModifier struct {
	Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// PlanModifyString runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m stringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic plan modifiers, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.PlanValue, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if !str.IsNull() {
		return
	}

	resp.PlanValue = types.StringValue(m.Default)
}

type boolDefaultModifier struct {
	Default bool
}

func boolDefault(defaultValue bool) boolDefaultModifier {
	return boolDefaultModifier{
		Default: defaultValue,
	}
}

func (m boolDefaultModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// types.Bool must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic plan modifiers, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.Bool
	diags := tfsdk.ValueAs(ctx, req.PlanValue, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if !str.IsNull() {
		return
	}

	resp.PlanValue = types.BoolValue(m.Default)
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %t", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%t`", m.Default)
}
