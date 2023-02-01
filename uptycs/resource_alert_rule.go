package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *alertRuleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"code": schema.StringAttribute{
				Required: true,
				Computed: false,
			},
			"type": schema.StringAttribute{
				Required: true,
				Computed: false,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive([]string{"sql"}...),
				},
			},
			"rule":     schema.StringAttribute{Required: true},
			"grouping": schema.StringAttribute{Required: true},
			"enabled": schema.BoolAttribute{Required: false,
				Computed: true,
			},
			"throttled":   schema.BoolAttribute{Required: true},
			"is_internal": schema.BoolAttribute{Required: true},
			"alert_tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"grouping_l2":     schema.StringAttribute{Required: true},
			"grouping_l3":     schema.StringAttribute{Required: true},
			"notify_interval": schema.Int64Attribute{Required: true},
			"notify_count":    schema.Int64Attribute{Required: true},
			"rule_exceptions": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"destinations": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"severity":             schema.StringAttribute{Optional: true},
						"destination_id":       schema.StringAttribute{Optional: true},
						"notify_every_alert":   schema.BoolAttribute{Optional: true},
						"close_after_delivery": schema.BoolAttribute{Optional: true},
					},
				},
			},
			"sql_config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"interval_seconds": schema.Int64Attribute{Optional: true},
				},
			},
		},
	}
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
		ID:                  types.StringValue(alertRuleResp.ID),
		Name:                types.StringValue(alertRuleResp.Name),
		Description:         types.StringValue(alertRuleResp.Description),
		Code:                types.StringValue(alertRuleResp.Code),
		Type:                types.StringValue(alertRuleResp.Type),
		Rule:                types.StringValue(alertRuleResp.Rule),
		Grouping:            types.StringValue(alertRuleResp.Grouping),
		Enabled:             types.BoolValue(alertRuleResp.Enabled),
		Throttled:           types.BoolValue(alertRuleResp.Throttled),
		IsInternal:          types.BoolValue(alertRuleResp.IsInternal),
		AlertNotifyCount:    types.Int64Value(int64(alertRuleResp.AlertNotifyCount)),
		AlertNotifyInterval: types.Int64Value(int64(alertRuleResp.AlertNotifyInterval)),
		AlertTags:           makeListStringAttribute(alertRuleResp.AlertTags),
		GroupingL2:          types.StringValue(alertRuleResp.GroupingL2),
		GroupingL3:          types.StringValue(alertRuleResp.GroupingL3),
		AlertRuleExceptions: makeListStringAttributeFn(alertRuleResp.AlertRuleExceptions, func(v uptycs.RuleException) (string, bool) { return v.ExceptionID, true }),
	}

	if alertRuleResp.SQLConfig != nil {
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: types.Int64Value(int64(alertRuleResp.SQLConfig.IntervalSeconds)),
		}
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
			Severity:           d.Severity.ValueString(),
			DestinationID:      d.DestinationID.ValueString(),
			NotifyEveryAlert:   d.NotifyEveryAlert.ValueBool(),
			CloseAfterDelivery: d.CloseAfterDelivery.ValueBool(),
		})
	}

	alertRule := uptycs.AlertRule{
		Name:                plan.Name.ValueString(),
		Description:         plan.Description.ValueString(),
		Code:                plan.Code.ValueString(),
		Type:                plan.Type.ValueString(),
		Rule:                plan.Rule.ValueString(),
		Grouping:            plan.Grouping.ValueString(),
		Enabled:             plan.Enabled.ValueBool(),
		Throttled:           plan.Throttled.ValueBool(),
		IsInternal:          plan.IsInternal.ValueBool(),
		AlertTags:           tags,
		GroupingL2:          plan.GroupingL2.ValueString(),
		GroupingL3:          plan.GroupingL3.ValueString(),
		AlertRuleExceptions: _ruleExceptions,
		Destinations:        _destinations,
		AlertNotifyInterval: int(plan.AlertNotifyInterval.ValueInt64()),
		AlertNotifyCount:    int(plan.AlertNotifyCount.ValueInt64()),
	}

	if plan.SQLConfig != nil {
		alertRule.SQLConfig = &uptycs.SQLConfig{
			IntervalSeconds: int(plan.SQLConfig.IntervalSeconds.ValueInt64()),
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
		ID:                  types.StringValue(alertRuleResp.ID),
		Name:                types.StringValue(alertRuleResp.Name),
		Description:         types.StringValue(alertRuleResp.Description),
		Code:                types.StringValue(alertRuleResp.Code),
		Type:                types.StringValue(alertRuleResp.Type),
		Rule:                types.StringValue(alertRuleResp.Rule),
		Grouping:            types.StringValue(alertRuleResp.Grouping),
		Enabled:             types.BoolValue(alertRuleResp.Enabled),
		Throttled:           types.BoolValue(alertRuleResp.Throttled),
		IsInternal:          types.BoolValue(alertRuleResp.IsInternal),
		AlertNotifyCount:    types.Int64Value(int64(alertRuleResp.AlertNotifyCount)),
		AlertNotifyInterval: types.Int64Value(int64(alertRuleResp.AlertNotifyInterval)),
		AlertTags:           makeListStringAttribute(alertRuleResp.AlertTags),
		GroupingL2:          types.StringValue(alertRuleResp.GroupingL2),
		GroupingL3:          types.StringValue(alertRuleResp.GroupingL3),
		AlertRuleExceptions: makeListStringAttributeFn(alertRuleResp.AlertRuleExceptions, func(v uptycs.RuleException) (string, bool) { return v.ExceptionID, true }),
	}

	if alertRuleResp.SQLConfig != nil {
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: types.Int64Value(int64(alertRuleResp.SQLConfig.IntervalSeconds)),
		}
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

	alertRuleID := state.ID.ValueString()

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
			Severity:           d.Severity.ValueString(),
			DestinationID:      d.DestinationID.ValueString(),
			NotifyEveryAlert:   d.NotifyEveryAlert.ValueBool(),
			CloseAfterDelivery: d.CloseAfterDelivery.ValueBool(),
		})
	}

	alertRule := uptycs.AlertRule{
		ID:                  alertRuleID,
		Name:                plan.Name.ValueString(),
		Description:         plan.Description.ValueString(),
		Code:                plan.Code.ValueString(),
		Type:                plan.Type.ValueString(),
		Rule:                plan.Rule.ValueString(),
		Grouping:            plan.Grouping.ValueString(),
		Enabled:             plan.Enabled.ValueBool(),
		Throttled:           plan.Throttled.ValueBool(),
		IsInternal:          plan.IsInternal.ValueBool(),
		AlertNotifyInterval: int(plan.AlertNotifyInterval.ValueInt64()),
		AlertNotifyCount:    int(plan.AlertNotifyCount.ValueInt64()),
		AlertTags:           tags,
		GroupingL2:          plan.GroupingL2.ValueString(),
		GroupingL3:          plan.GroupingL3.ValueString(),
		AlertRuleExceptions: _ruleExceptions,
		Destinations:        _destinations,
	}

	if plan.SQLConfig != nil {
		alertRule.SQLConfig = &uptycs.SQLConfig{
			IntervalSeconds: int(plan.SQLConfig.IntervalSeconds.ValueInt64()),
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
		ID:                  types.StringValue(alertRuleResp.ID),
		Name:                types.StringValue(alertRuleResp.Name),
		Description:         types.StringValue(alertRuleResp.Description),
		Code:                types.StringValue(alertRuleResp.Code),
		Type:                types.StringValue(alertRuleResp.Type),
		Rule:                types.StringValue(alertRuleResp.Rule),
		Grouping:            types.StringValue(alertRuleResp.Grouping),
		Enabled:             types.BoolValue(alertRuleResp.Enabled),
		Throttled:           types.BoolValue(alertRuleResp.Throttled),
		IsInternal:          types.BoolValue(alertRuleResp.IsInternal),
		AlertNotifyCount:    types.Int64Value(int64(alertRuleResp.AlertNotifyCount)),
		AlertNotifyInterval: types.Int64Value(int64(alertRuleResp.AlertNotifyInterval)),
		AlertTags:           makeListStringAttribute(alertRuleResp.AlertTags),
		GroupingL2:          types.StringValue(alertRuleResp.GroupingL2),
		GroupingL3:          types.StringValue(alertRuleResp.GroupingL3),
		AlertRuleExceptions: makeListStringAttributeFn(alertRuleResp.AlertRuleExceptions, func(v uptycs.RuleException) (string, bool) { return v.ExceptionID, true }),
	}

	if alertRuleResp.SQLConfig != nil {
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: types.Int64Value(int64(alertRuleResp.SQLConfig.IntervalSeconds)),
		}
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

	alertRuleID := state.ID.ValueString()

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
