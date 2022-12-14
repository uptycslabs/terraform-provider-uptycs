package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &tagRuleResource{}
	_ resource.ResourceWithConfigure   = &tagRuleResource{}
	_ resource.ResourceWithImportState = &tagRuleResource{}
)

func TagRuleResource() resource.Resource {
	return &tagRuleResource{}
}

type tagRuleResource struct {
	client *uptycs.Client
}

func (r *tagRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag_rule"
}

func (r *tagRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *tagRuleResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
			"query": {
				Type:     types.StringType,
				Required: true,
			},
			"source": {
				Type:     types.StringType,
				Required: true,
			},
			"run_once": {
				Type:     types.BoolType,
				Required: true,
			},
			"interval": {
				Type:     types.NumberType,
				Optional: true,
			},
			"osquery_version": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"platform": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"resource_type": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("asset")},
			},
			"enabled": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"system": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
		},
	}, nil
}

func (r *tagRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tagRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &tagRuleID)...)
	tagRuleResp, err := r.client.GetTagRule(uptycs.TagRule{
		ID: tagRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get tagRule with ID  "+tagRuleID+": "+err.Error(),
		)
		return
	}

	var result = TagRule{
		ID:             types.StringValue(tagRuleResp.ID),
		Name:           types.StringValue(tagRuleResp.Name),
		Description:    types.StringValue(tagRuleResp.Description),
		Query:          types.StringValue(tagRuleResp.Query),
		Source:         types.StringValue(tagRuleResp.Source),
		RunOnce:        types.BoolValue(tagRuleResp.RunOnce),
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.StringValue(tagRuleResp.OSqueryVersion),
		Platform:       types.StringValue(tagRuleResp.Platform),
		Enabled:        types.BoolValue(tagRuleResp.Enabled),
		System:         types.BoolValue(tagRuleResp.System),
		ResourceType:   types.StringValue(tagRuleResp.ResourceType),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tagRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan TagRule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tagRuleResp, err := r.client.CreateTagRule(uptycs.TagRule{
		ID:             plan.ID.Value,
		Name:           plan.Name.Value,
		Description:    plan.Description.Value,
		Query:          plan.Query.Value,
		Source:         plan.Source.Value,
		RunOnce:        plan.RunOnce.Value,
		Interval:       plan.Interval,
		OSqueryVersion: plan.OSqueryVersion.Value,
		Platform:       plan.Platform.Value,
		Enabled:        plan.Enabled.Value,
		ResourceType:   plan.ResourceType.Value,
		// System:         plan.System.Value, //"error":{"status":400,"code":"INVALID_OR_REQUIRED_FIELD","message":{"brief":"","detail":"\"system\"│ is not allowed","developer":""}}}
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tagRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = TagRule{
		ID:             types.StringValue(tagRuleResp.ID),
		Name:           types.StringValue(tagRuleResp.Name),
		Description:    types.StringValue(tagRuleResp.Description),
		Query:          types.StringValue(tagRuleResp.Query),
		Source:         types.StringValue(tagRuleResp.Source),
		RunOnce:        types.BoolValue(tagRuleResp.RunOnce),
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.StringValue(tagRuleResp.OSqueryVersion),
		Platform:       types.StringValue(tagRuleResp.Platform),
		Enabled:        types.BoolValue(tagRuleResp.Enabled),
		System:         types.BoolValue(tagRuleResp.System),
		ResourceType:   types.StringValue(tagRuleResp.ResourceType),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state TagRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRuleID := state.ID.Value

	// Retrieve values from plan
	var plan TagRule
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRuleResp, err := r.client.UpdateTagRule(uptycs.TagRule{
		ID:             tagRuleID,
		Name:           plan.Name.Value,
		Description:    plan.Description.Value,
		Query:          plan.Query.Value,
		Source:         plan.Source.Value,
		RunOnce:        plan.RunOnce.Value,
		Interval:       plan.Interval,
		OSqueryVersion: plan.OSqueryVersion.Value,
		Platform:       plan.Platform.Value,
		Enabled:        plan.Enabled.Value,
		ResourceType:   plan.ResourceType.Value,
		// System:         plan.System.Value, //"error":{"status":400,"code":"INVALID_OR_REQUIRED_FIELD","message":{"brief":"","detail":"\"system\"│ is not allowed","developer":""}}}
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tagRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = TagRule{
		ID:             types.StringValue(tagRuleResp.ID),
		Name:           types.StringValue(tagRuleResp.Name),
		Description:    types.StringValue(tagRuleResp.Description),
		Query:          types.StringValue(tagRuleResp.Query),
		Source:         types.StringValue(tagRuleResp.Source),
		RunOnce:        types.BoolValue(tagRuleResp.RunOnce),
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.StringValue(tagRuleResp.OSqueryVersion),
		Platform:       types.StringValue(tagRuleResp.Platform),
		Enabled:        types.BoolValue(tagRuleResp.Enabled),
		System:         types.BoolValue(tagRuleResp.System),
		ResourceType:   types.StringValue(tagRuleResp.ResourceType),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TagRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRuleID := state.ID.Value

	_, err := r.client.DeleteTagRule(uptycs.TagRule{
		ID: tagRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete tagRule with ID  "+tagRuleID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r *tagRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
