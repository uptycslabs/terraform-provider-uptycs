package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceTagRuleType struct{}

// Alert Rule Resource schema
func (r resourceTagRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Type:     types.StringType,
				Required: true,
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

// New resource instance
func (r resourceTagRuleType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceTagRule{
		p: *(p.(*Provider)),
	}, nil
}

type resourceTagRule struct {
	p Provider
}

// Read resource information
func (r resourceTagRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tagRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &tagRuleID)...)
	tagRuleResp, err := r.p.client.GetTagRule(uptycs.TagRule{
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
		ID:             types.String{Value: tagRuleResp.ID},
		Name:           types.String{Value: tagRuleResp.Name},
		Description:    types.String{Value: tagRuleResp.Description},
		Query:          types.String{Value: tagRuleResp.Query},
		Source:         types.String{Value: tagRuleResp.Source},
		RunOnce:        types.Bool{Value: tagRuleResp.RunOnce},
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.String{Value: tagRuleResp.OSqueryVersion},
		Platform:       types.String{Value: tagRuleResp.Platform},
		Enabled:        types.Bool{Value: tagRuleResp.Enabled},
		System:         types.Bool{Value: tagRuleResp.System},
		ResourceType:   types.String{Value: tagRuleResp.ResourceType},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Create a new resource
func (r resourceTagRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan TagRule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tagRuleResp, err := r.p.client.CreateTagRule(uptycs.TagRule{
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
		System:         plan.System.Value,
		ResourceType:   plan.ResourceType.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tagRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = TagRule{
		ID:             types.String{Value: tagRuleResp.ID},
		Name:           types.String{Value: tagRuleResp.Name},
		Description:    types.String{Value: tagRuleResp.Description},
		Query:          types.String{Value: tagRuleResp.Query},
		Source:         types.String{Value: tagRuleResp.Source},
		RunOnce:        types.Bool{Value: tagRuleResp.RunOnce},
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.String{Value: tagRuleResp.OSqueryVersion},
		Platform:       types.String{Value: tagRuleResp.Platform},
		Enabled:        types.Bool{Value: tagRuleResp.Enabled},
		System:         types.Bool{Value: tagRuleResp.System},
		ResourceType:   types.String{Value: tagRuleResp.ResourceType},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceTagRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	tagRuleResp, err := r.p.client.UpdateTagRule(uptycs.TagRule{
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
		System:         plan.System.Value,
		ResourceType:   plan.ResourceType.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tagRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = TagRule{
		ID:             types.String{Value: tagRuleResp.ID},
		Name:           types.String{Value: tagRuleResp.Name},
		Description:    types.String{Value: tagRuleResp.Description},
		Query:          types.String{Value: tagRuleResp.Query},
		Source:         types.String{Value: tagRuleResp.Source},
		RunOnce:        types.Bool{Value: tagRuleResp.RunOnce},
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.String{Value: tagRuleResp.OSqueryVersion},
		Platform:       types.String{Value: tagRuleResp.Platform},
		Enabled:        types.Bool{Value: tagRuleResp.Enabled},
		System:         types.Bool{Value: tagRuleResp.System},
		ResourceType:   types.String{Value: tagRuleResp.ResourceType},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceTagRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TagRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRuleID := state.ID.Value

	_, err := r.p.client.DeleteTagRule(uptycs.TagRule{
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

// Import resource
func (r resourceTagRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
