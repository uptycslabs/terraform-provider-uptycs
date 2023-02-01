package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *alertRuleDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"code":        schema.StringAttribute{Optional: true},
			"type":        schema.StringAttribute{Optional: true},
			"rule":        schema.StringAttribute{Optional: true},
			"grouping":    schema.StringAttribute{Optional: true},
			"enabled":     schema.BoolAttribute{Optional: true},
			"throttled":   schema.BoolAttribute{Optional: true},
			"is_internal": schema.BoolAttribute{Optional: true},
			"alert_tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"grouping_l2":     schema.StringAttribute{Optional: true},
			"grouping_l3":     schema.StringAttribute{Optional: true},
			"notify_interval": schema.Int64Attribute{Optional: true},
			"notify_count":    schema.Int64Attribute{Optional: true},
			"rule_exceptions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"destinations": schema.ListNestedAttribute{
				Optional: true,
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
