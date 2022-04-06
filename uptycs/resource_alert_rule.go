package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/myoung34/uptycs-client-go/uptycs"
)

type resourceAlertRuleType struct{}

// Alert Rule Resource schema
func (r resourceAlertRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
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
			"code": {
				Type:     types.StringType,
				Required: true,
				Computed: false,
			},
			"type": {
				Type:     types.StringType,
				Required: true,
				Computed: false,
			},
			"rule": {
				Type:     types.StringType,
				Required: true,
			},
			"grouping": {
				Type:     types.StringType,
				Required: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Required: true,
			},
			"grouping_l2": {
				Type:     types.StringType,
				Required: true,
			},
			"grouping_l3": {
				Type:     types.StringType,
				Required: true,
			},
			"sql_config": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"interval_seconds": {
						Type:     types.NumberType,
						Required: true,
					},
				}),
			},
		},
	}, nil
}

// New resource instance
func (r resourceAlertRuleType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceAlertRule{
		p: *(p.(*provider)),
	}, nil
}

type resourceAlertRule struct {
	p provider
}

// Create a new resource
func (r resourceAlertRule) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan AlertRule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertRuleResp, err := r.p.client.CreateAlertRule(uptycs.AlertRule{
		Name:        plan.Name.Value,
		Code:        plan.Code.Value,
		Description: plan.Description.Value,
		Rule:        plan.Rule.Value,
		Type:        plan.Type.Value,
		Enabled:     plan.Enabled.Value,
		SQLConfig: uptycs.SQLConfig{
			IntervalSeconds: plan.SQLConfig.IntervalSeconds,
		},
		Grouping:   plan.Grouping.Value,
		GroupingL2: plan.GroupingL2.Value,
		GroupingL3: plan.GroupingL3.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating alertRule",
			"Could not create alertRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.String{Value: alertRuleResp.ID},
		Enabled:     types.Bool{Value: alertRuleResp.Enabled},
		Name:        types.String{Value: alertRuleResp.Name},
		Description: types.String{Value: alertRuleResp.Description},
		Code:        types.String{Value: alertRuleResp.Code},
		Type:        types.String{Value: alertRuleResp.Type},
		Rule:        types.String{Value: alertRuleResp.Rule},
		SQLConfig: SQLConfig{
			IntervalSeconds: alertRuleResp.SQLConfig.IntervalSeconds,
		},
		Grouping:   types.String{Value: alertRuleResp.Grouping},
		GroupingL2: types.String{Value: alertRuleResp.GroupingL2},
		GroupingL3: types.String{Value: alertRuleResp.GroupingL3},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceAlertRule) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var alertRuleId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &alertRuleId)...)
	alertRuleResp, err := r.p.client.GetAlertRule(uptycs.AlertRule{
		ID: alertRuleId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading order",
			"Could not get alertRule with ID  "+alertRuleId+": "+err.Error(),
		)
		return
	}
	var result = AlertRule{
		ID:          types.String{Value: alertRuleResp.ID},
		Enabled:     types.Bool{Value: alertRuleResp.Enabled},
		Name:        types.String{Value: alertRuleResp.Name},
		Description: types.String{Value: alertRuleResp.Description},
		Code:        types.String{Value: alertRuleResp.Code},
		Type:        types.String{Value: alertRuleResp.Type},
		Rule:        types.String{Value: alertRuleResp.Rule},
		SQLConfig: SQLConfig{
			IntervalSeconds: alertRuleResp.SQLConfig.IntervalSeconds,
		},
		Grouping:   types.String{Value: alertRuleResp.Grouping},
		GroupingL2: types.String{Value: alertRuleResp.GroupingL2},
		GroupingL3: types.String{Value: alertRuleResp.GroupingL3},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update resource
func (r resourceAlertRule) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state AlertRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertRuleID := state.ID.Value

	// Retrieve values from plan
	var plan AlertRule
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertRuleResp, err := r.p.client.UpdateAlertRule(uptycs.AlertRule{
		ID:          alertRuleID,
		Name:        plan.Name.Value,
		Code:        plan.Code.Value,
		Description: plan.Description.Value,
		Rule:        plan.Rule.Value,
		Type:        plan.Type.Value,
		Enabled:     plan.Enabled.Value,
		SQLConfig: uptycs.SQLConfig{
			IntervalSeconds: plan.SQLConfig.IntervalSeconds,
		},
		Grouping:   plan.Grouping.Value,
		GroupingL2: plan.GroupingL2.Value,
		GroupingL3: plan.GroupingL3.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating alertRule",
			"Could not create alertRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.String{Value: alertRuleResp.ID},
		Enabled:     types.Bool{Value: alertRuleResp.Enabled},
		Name:        types.String{Value: alertRuleResp.Name},
		Description: types.String{Value: alertRuleResp.Description},
		Code:        types.String{Value: alertRuleResp.Code},
		Type:        types.String{Value: alertRuleResp.Type},
		Rule:        types.String{Value: alertRuleResp.Rule},
		SQLConfig: SQLConfig{
			IntervalSeconds: alertRuleResp.SQLConfig.IntervalSeconds,
		},
		Grouping:   types.String{Value: alertRuleResp.Grouping},
		GroupingL2: types.String{Value: alertRuleResp.GroupingL2},
		GroupingL3: types.String{Value: alertRuleResp.GroupingL3},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceAlertRule) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state AlertRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertRuleID := state.ID.Value

	_, err := r.p.client.DeleteAlertRule(uptycs.AlertRule{
		ID: alertRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting order",
			"Could not delete alertRule with ID  "+alertRuleID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceAlertRule) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
