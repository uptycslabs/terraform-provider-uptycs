package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &alertRuleResource{}
	_ resource.ResourceWithConfigure   = &alertRuleResource{}
	_ resource.ResourceWithImportState = &alertRuleResource{}
)

func AlertRuleResource() resource.Resource {
	return &alertRuleResource{}
}

type alertRuleResource struct {
	client *uptycs.Client
}

func (r *alertRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_rule"
}

// Configure resource instance
func (r *alertRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *alertRuleResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r *alertRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var alertRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &alertRuleID)...)
	alertRuleResp, err := r.client.GetAlertRule(uptycs.AlertRule{
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
		ID:          types.StringValue(alertRuleResp.ID),
		Name:        types.StringValue(alertRuleResp.Name),
		Description: types.StringValue(alertRuleResp.Description),
		Code:        types.StringValue(alertRuleResp.Code),
		Type:        types.StringValue(alertRuleResp.Type),
		Rule:        types.StringValue(alertRuleResp.Rule),
		Grouping:    types.StringValue(alertRuleResp.Grouping),
		Enabled:     types.BoolValue(alertRuleResp.Enabled),
		Throttled:   types.BoolValue(alertRuleResp.Throttled),
		IsInternal:  types.BoolValue(alertRuleResp.IsInternal),
		AlertTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		GroupingL2: types.StringValue(alertRuleResp.GroupingL2),
		GroupingL3: types.StringValue(alertRuleResp.GroupingL3),
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
			Severity:           types.StringValue(d.Severity),
			DestinationID:      types.StringValue(d.DestinationID),
			NotifyEveryAlert:   types.BoolValue(d.NotifyEveryAlert),
			CloseAfterDelivery: types.BoolValue(d.CloseAfterDelivery),
		})
	}
	result.Destinations = destinations

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *alertRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	alertRule := uptycs.AlertRule{
		Name:                plan.Name.Value,
		Description:         plan.Description.Value,
		Code:                plan.Code.Value,
		Type:                plan.Type.Value,
		Rule:                plan.Rule.Value,
		Grouping:            plan.Grouping.Value,
		Enabled:             plan.Enabled.Value,
		Throttled:           plan.Throttled.Value,
		IsInternal:          plan.IsInternal.Value,
		AlertTags:           tags,
		GroupingL2:          plan.GroupingL2.Value,
		GroupingL3:          plan.GroupingL3.Value,
		AlertRuleExceptions: _ruleExceptions,
		Destinations:        _destinations,
	}
	if plan.AlertNotifyInterval != nil {
		alertRule.AlertNotifyInterval = *plan.AlertNotifyInterval
	}

	if plan.AlertNotifyCount != nil {
		alertRule.AlertNotifyCount = *plan.AlertNotifyCount
	}

	if plan.SQLConfig != nil {
		alertRule.SQLConfig = &uptycs.SQLConfig{
			IntervalSeconds: plan.SQLConfig.IntervalSeconds,
		}
	}
	alertRuleResp, err := r.client.CreateAlertRule(alertRule)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create alertRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.StringValue(alertRuleResp.ID),
		Name:        types.StringValue(alertRuleResp.Name),
		Description: types.StringValue(alertRuleResp.Description),
		Code:        types.StringValue(alertRuleResp.Code),
		Type:        types.StringValue(alertRuleResp.Type),
		Rule:        types.StringValue(alertRuleResp.Rule),
		Grouping:    types.StringValue(alertRuleResp.Grouping),
		Enabled:     types.BoolValue(alertRuleResp.Enabled),
		Throttled:   types.BoolValue(alertRuleResp.Throttled),
		IsInternal:  types.BoolValue(alertRuleResp.IsInternal),
		AlertTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		GroupingL2: types.StringValue(alertRuleResp.GroupingL2),
		GroupingL3: types.StringValue(alertRuleResp.GroupingL3),
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
			Severity:           types.StringValue(d.Severity),
			DestinationID:      types.StringValue(d.DestinationID),
			NotifyEveryAlert:   types.BoolValue(d.NotifyEveryAlert),
			CloseAfterDelivery: types.BoolValue(d.CloseAfterDelivery),
		})
	}
	result.Destinations = destinations

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *alertRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	alertRule := uptycs.AlertRule{
		ID:                  alertRuleID,
		Name:                plan.Name.Value,
		Description:         plan.Description.Value,
		Code:                plan.Code.Value,
		Type:                plan.Type.Value,
		Rule:                plan.Rule.Value,
		Grouping:            plan.Grouping.Value,
		Enabled:             plan.Enabled.Value,
		Throttled:           plan.Throttled.Value,
		IsInternal:          plan.IsInternal.Value,
		AlertTags:           tags,
		GroupingL2:          plan.GroupingL2.Value,
		GroupingL3:          plan.GroupingL3.Value,
		AlertRuleExceptions: _ruleExceptions,
		Destinations:        _destinations,
	}
	if plan.AlertNotifyInterval != nil {
		alertRule.AlertNotifyInterval = *plan.AlertNotifyInterval
	}

	if plan.AlertNotifyCount != nil {
		alertRule.AlertNotifyCount = *plan.AlertNotifyCount
	}

	if plan.SQLConfig != nil {
		alertRule.SQLConfig = &uptycs.SQLConfig{
			IntervalSeconds: plan.SQLConfig.IntervalSeconds,
		}
	}

	alertRuleResp, err := r.client.UpdateAlertRule(alertRule)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create alertRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = AlertRule{
		ID:          types.StringValue(alertRuleResp.ID),
		Name:        types.StringValue(alertRuleResp.Name),
		Description: types.StringValue(alertRuleResp.Description),
		Code:        types.StringValue(alertRuleResp.Code),
		Type:        types.StringValue(alertRuleResp.Type),
		Rule:        types.StringValue(alertRuleResp.Rule),
		Grouping:    types.StringValue(alertRuleResp.Grouping),
		Enabled:     types.BoolValue(alertRuleResp.Enabled),
		Throttled:   types.BoolValue(alertRuleResp.Throttled),
		IsInternal:  types.BoolValue(alertRuleResp.IsInternal),
		AlertTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		GroupingL2: types.StringValue(alertRuleResp.GroupingL2),
		GroupingL3: types.StringValue(alertRuleResp.GroupingL3),
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
			Severity:           types.StringValue(d.Severity),
			DestinationID:      types.StringValue(d.DestinationID),
			NotifyEveryAlert:   types.BoolValue(d.NotifyEveryAlert),
			CloseAfterDelivery: types.BoolValue(d.CloseAfterDelivery),
		})
	}
	result.Destinations = destinations

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *alertRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertRuleID := state.ID.Value

	_, err := r.client.DeleteAlertRule(uptycs.AlertRule{
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

func (r *alertRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
