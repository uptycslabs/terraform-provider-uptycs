package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &alertRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &alertRuleDataSource{}
)

func AlertRuleDataSource() datasource.DataSource {
	return &alertRuleDataSource{}
}

type alertRuleDataSource struct {
	client *uptycs.Client
}

func (d *alertRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_rule"
}

func (d *alertRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *alertRuleDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
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

func (d *alertRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var alertRuleID string
	var alertRuleName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &alertRuleID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &alertRuleName)

	var alertRuleToLookup uptycs.AlertRule

	if len(alertRuleID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		alertRuleToLookup = uptycs.AlertRule{
			Name: alertRuleName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		alertRuleToLookup = uptycs.AlertRule{
			ID: alertRuleID,
		}
	}

	alertRuleResp, err := d.client.GetAlertRule(alertRuleToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
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
