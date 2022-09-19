package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceAlertRuleType struct {
	p Provider
}

func (r dataSourceAlertRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"code": {
				Type:     types.StringType,
				Optional: true,
			},
			"type": {
				Type:     types.StringType,
				Optional: true,
			},
			"rule": {
				Type:     types.StringType,
				Optional: true,
			},
			"grouping": {
				Type:     types.StringType,
				Optional: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Optional: true,
			},
			"custom": {
				Type:     types.BoolType,
				Optional: true,
			},
			"throttled": {
				Type:     types.BoolType,
				Optional: true,
			},
			"is_internal": {
				Type:     types.BoolType,
				Optional: true,
			},
			"alert_tags": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"grouping_l2": {
				Type:     types.StringType,
				Optional: true,
			},
			"grouping_l3": {
				Type:     types.StringType,
				Optional: true,
			},
			"lock": {
				Type:     types.BoolType,
				Optional: true,
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
				Optional: true,
			},
			"destinations": {
				Optional: true,
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

func (r dataSourceAlertRuleType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceAlertRuleType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceAlertRuleType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var alertRuleID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &alertRuleID)...)

	alertRuleResp, err := r.p.client.GetAlertRule(uptycs.AlertRule{
		ID: alertRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
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
