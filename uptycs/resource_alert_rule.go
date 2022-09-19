package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"throttled": {
				Type:     types.BoolType,
				Required: true,
			},
			"is_internal": {
				Type:     types.BoolType,
				Required: true,
			},
			"alert_tags": {
				Type:     types.ListType{ElemType: types.StringType},
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
			"lock": {
				Type:     types.BoolType,
				Required: true,
			},
			"notify_interval": {
				Type:     types.NumberType,
				Optional: true,
			},
			"notify_count": {
				Type:     types.NumberType,
				Optional: true,
			},
			"rule_exceptions": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"destinations": {
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"severity": {
							Type:     types.StringType,
							Optional: true,
						},
						"destination_id": {
							Type:     types.StringType,
							Optional: true,
						},
						"notify_every_alert": {
							Type:     types.BoolType,
							Optional: true,
						},
						"close_after_delivery": {
							Type:     types.BoolType,
							Optional: true,
						},
					},
				),
			},
			"sql_config": {
				Optional: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"interval_seconds": {
						Type:     types.NumberType,
						Optional: true,
					},
				}),
			},
		},
	}, nil
}

// New resource instance
func (r resourceAlertRuleType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceAlertRule{
		p: *(p.(*Provider)),
	}, nil
}

type resourceAlertRule struct {
	p Provider
}

// Read resource information
func (r resourceAlertRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var alertRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &alertRuleID)...)
	alertRuleResp, err := r.p.client.GetAlertRule(uptycs.AlertRule{
		ID: alertRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get alertRule with ID  "+alertRuleID+": "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.String{Value: alertRuleResp.ID},
		Name:        types.String{Value: alertRuleResp.Name},
		Description: types.String{Value: alertRuleResp.Description},
		Code:        types.String{Value: alertRuleResp.Code},
		Type:        types.String{Value: alertRuleResp.Type},
		Rule:        types.String{Value: alertRuleResp.Rule},
		Grouping:    types.String{Value: alertRuleResp.Grouping},
		Enabled:     types.Bool{Value: alertRuleResp.Enabled},
		Custom:      types.Bool{Value: alertRuleResp.Custom},
		Throttled:   types.Bool{Value: alertRuleResp.Throttled},
		IsInternal:  types.Bool{Value: alertRuleResp.IsInternal},
		AlertTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		GroupingL2: types.String{Value: alertRuleResp.GroupingL2},
		GroupingL3: types.String{Value: alertRuleResp.GroupingL3},
		Lock:       types.Bool{Value: alertRuleResp.Lock},
		AlertRuleExceptions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	if alertRuleResp.AlertNotifyInterval != 0 {
		result.AlertNotifyInterval = &alertRuleResp.AlertNotifyInterval
	}
	if alertRuleResp.AlertNotifyCount != 0 {
		result.AlertNotifyCount = &alertRuleResp.AlertNotifyCount
	}
	if alertRuleResp.SQLConfig != nil {
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: alertRuleResp.SQLConfig.IntervalSeconds,
		}
	}

	for _, at := range alertRuleResp.AlertTags {
		result.AlertTags.Elems = append(result.AlertTags.Elems, types.String{Value: at})
	}

	for _, are := range alertRuleResp.AlertRuleExceptions {
		result.AlertRuleExceptions.Elems = append(result.AlertRuleExceptions.Elems, types.String{Value: are.ExceptionID})
	}

	destinations := make([]AlertRuleDestination, 0)
	for _, d := range alertRuleResp.Destinations {
		destinations = append(destinations, AlertRuleDestination{
			Severity:           types.String{Value: d.Severity},
			DestinationID:      types.String{Value: d.DestinationID},
			NotifyEveryAlert:   types.Bool{Value: d.NotifyEveryAlert},
			CloseAfterDelivery: types.Bool{Value: d.CloseAfterDelivery},
		})
	}
	result.Destinations = destinations

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Create a new resource
func (r resourceAlertRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	var tags []string
	plan.AlertTags.ElementsAs(ctx, &tags, false)

	var ruleExceptions []string
	plan.AlertRuleExceptions.ElementsAs(ctx, &ruleExceptions, false)
	_ruleExceptions := make([]uptycs.RuleException, 0)
	for _, _re := range ruleExceptions {
		_ruleExceptions = append(_ruleExceptions, uptycs.RuleException{
			ExceptionID: _re,
		})
	}

	_destinations := make([]uptycs.AlertRuleDestination, 0)
	for _, d := range plan.Destinations {
		_destinations = append(_destinations, uptycs.AlertRuleDestination{
			Severity:           d.Severity.Value,
			DestinationID:      d.DestinationID.Value,
			NotifyEveryAlert:   d.NotifyEveryAlert.Value,
			CloseAfterDelivery: d.CloseAfterDelivery.Value,
		})
	}

	alertRuleResp, err := r.p.client.CreateAlertRule(uptycs.AlertRule{
		Name:                plan.Name.Value,
		Description:         plan.Description.Value,
		Code:                plan.Code.Value,
		Type:                plan.Type.Value,
		Rule:                plan.Rule.Value,
		Grouping:            plan.Grouping.Value,
		Enabled:             plan.Enabled.Value,
		Custom:              plan.Custom.Value,
		Throttled:           plan.Throttled.Value,
		IsInternal:          plan.IsInternal.Value,
		AlertTags:           tags,
		GroupingL2:          plan.GroupingL2.Value,
		GroupingL3:          plan.GroupingL3.Value,
		Lock:                plan.Lock.Value,
		AlertNotifyInterval: *plan.AlertNotifyInterval,
		AlertNotifyCount:    *plan.AlertNotifyCount,
		AlertRuleExceptions: _ruleExceptions,
		Destinations:        _destinations,
		SQLConfig: &uptycs.SQLConfig{
			IntervalSeconds: plan.SQLConfig.IntervalSeconds,
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create alertRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.String{Value: alertRuleResp.ID},
		Name:        types.String{Value: alertRuleResp.Name},
		Description: types.String{Value: alertRuleResp.Description},
		Code:        types.String{Value: alertRuleResp.Code},
		Type:        types.String{Value: alertRuleResp.Type},
		Rule:        types.String{Value: alertRuleResp.Rule},
		Grouping:    types.String{Value: alertRuleResp.Grouping},
		Enabled:     types.Bool{Value: alertRuleResp.Enabled},
		Custom:      types.Bool{Value: alertRuleResp.Custom},
		Throttled:   types.Bool{Value: alertRuleResp.Throttled},
		IsInternal:  types.Bool{Value: alertRuleResp.IsInternal},
		AlertTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		GroupingL2: types.String{Value: alertRuleResp.GroupingL2},
		GroupingL3: types.String{Value: alertRuleResp.GroupingL3},
		Lock:       types.Bool{Value: alertRuleResp.Lock},
		AlertRuleExceptions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	if alertRuleResp.AlertNotifyInterval != 0 {
		result.AlertNotifyInterval = &alertRuleResp.AlertNotifyInterval
	}
	if alertRuleResp.AlertNotifyCount != 0 {
		result.AlertNotifyCount = &alertRuleResp.AlertNotifyCount
	}

	if alertRuleResp.SQLConfig != nil {
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: alertRuleResp.SQLConfig.IntervalSeconds,
		}
	}

	for _, at := range alertRuleResp.AlertTags {
		result.AlertTags.Elems = append(result.AlertTags.Elems, types.String{Value: at})
	}

	for _, are := range alertRuleResp.AlertRuleExceptions {
		result.AlertRuleExceptions.Elems = append(result.AlertRuleExceptions.Elems, types.String{Value: are.ExceptionID})
	}

	destinations := make([]AlertRuleDestination, 0)
	for _, d := range alertRuleResp.Destinations {
		destinations = append(destinations, AlertRuleDestination{
			Severity:           types.String{Value: d.Severity},
			DestinationID:      types.String{Value: d.DestinationID},
			NotifyEveryAlert:   types.Bool{Value: d.NotifyEveryAlert},
			CloseAfterDelivery: types.Bool{Value: d.CloseAfterDelivery},
		})
	}
	result.Destinations = destinations

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceAlertRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	var tags []string
	plan.AlertTags.ElementsAs(ctx, &tags, false)

	var ruleExceptions []string
	plan.AlertRuleExceptions.ElementsAs(ctx, &ruleExceptions, false)
	_ruleExceptions := make([]uptycs.RuleException, 0)
	for _, _re := range ruleExceptions {
		_ruleExceptions = append(_ruleExceptions, uptycs.RuleException{
			ExceptionID: _re,
		})
	}

	_destinations := make([]uptycs.AlertRuleDestination, 0)
	for _, d := range plan.Destinations {
		_destinations = append(_destinations, uptycs.AlertRuleDestination{
			Severity:           d.Severity.Value,
			DestinationID:      d.DestinationID.Value,
			NotifyEveryAlert:   d.NotifyEveryAlert.Value,
			CloseAfterDelivery: d.CloseAfterDelivery.Value,
		})
	}

	alertRuleResp, err := r.p.client.UpdateAlertRule(uptycs.AlertRule{
		ID:                  alertRuleID,
		Name:                plan.Name.Value,
		Description:         plan.Description.Value,
		Code:                plan.Code.Value,
		Type:                plan.Type.Value,
		Rule:                plan.Rule.Value,
		Grouping:            plan.Grouping.Value,
		Enabled:             plan.Enabled.Value,
		Custom:              plan.Custom.Value,
		Throttled:           plan.Throttled.Value,
		IsInternal:          plan.IsInternal.Value,
		AlertTags:           tags,
		GroupingL2:          plan.GroupingL2.Value,
		GroupingL3:          plan.GroupingL3.Value,
		Lock:                plan.Lock.Value,
		AlertNotifyInterval: *plan.AlertNotifyInterval,
		AlertNotifyCount:    *plan.AlertNotifyCount,
		AlertRuleExceptions: _ruleExceptions,
		Destinations:        _destinations,
		SQLConfig: &uptycs.SQLConfig{
			IntervalSeconds: plan.SQLConfig.IntervalSeconds,
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create alertRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.String{Value: alertRuleResp.ID},
		Name:        types.String{Value: alertRuleResp.Name},
		Description: types.String{Value: alertRuleResp.Description},
		Code:        types.String{Value: alertRuleResp.Code},
		Type:        types.String{Value: alertRuleResp.Type},
		Rule:        types.String{Value: alertRuleResp.Rule},
		Grouping:    types.String{Value: alertRuleResp.Grouping},
		Enabled:     types.Bool{Value: alertRuleResp.Enabled},
		Custom:      types.Bool{Value: alertRuleResp.Custom},
		Throttled:   types.Bool{Value: alertRuleResp.Throttled},
		IsInternal:  types.Bool{Value: alertRuleResp.IsInternal},
		AlertTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		GroupingL2: types.String{Value: alertRuleResp.GroupingL2},
		GroupingL3: types.String{Value: alertRuleResp.GroupingL3},
		Lock:       types.Bool{Value: alertRuleResp.Lock},
		AlertRuleExceptions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	if alertRuleResp.AlertNotifyInterval != 0 {
		result.AlertNotifyInterval = &alertRuleResp.AlertNotifyInterval
	}
	if alertRuleResp.AlertNotifyCount != 0 {
		result.AlertNotifyCount = &alertRuleResp.AlertNotifyCount
	}

	if alertRuleResp.SQLConfig != nil {
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: alertRuleResp.SQLConfig.IntervalSeconds,
		}
	}

	for _, at := range alertRuleResp.AlertTags {
		result.AlertTags.Elems = append(result.AlertTags.Elems, types.String{Value: at})
	}

	for _, are := range alertRuleResp.AlertRuleExceptions {
		result.AlertRuleExceptions.Elems = append(result.AlertRuleExceptions.Elems, types.String{Value: are.ExceptionID})
	}

	destinations := make([]AlertRuleDestination, 0)
	for _, d := range alertRuleResp.Destinations {
		destinations = append(destinations, AlertRuleDestination{
			Severity:           types.String{Value: d.Severity},
			DestinationID:      types.String{Value: d.DestinationID},
			NotifyEveryAlert:   types.Bool{Value: d.NotifyEveryAlert},
			CloseAfterDelivery: types.Bool{Value: d.CloseAfterDelivery},
		})
	}
	result.Destinations = destinations

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceAlertRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
			"Error Deleting",
			"Could not delete alertRule with ID  "+alertRuleID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceAlertRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
