package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func AlertRuleCategoryDataSource() datasource.DataSource {
	return &alertRuleCategoryDataSource{}
}

type alertRuleCategoryDataSource struct {
	client *uptycs.Client
}

func (d *alertRuleCategoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_rule_category"
}

func (d *alertRuleCategoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *alertRuleCategoryDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Optional: true},
			"rule_id": schema.StringAttribute{Optional: true},
			"name":    schema.StringAttribute{Optional: true},
		},
	}
}

func (d *alertRuleCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var alertRuleCategoryID string
	var alertRuleCategoryName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &alertRuleCategoryID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &alertRuleCategoryName)

	var alertRuleCategoryToLookup uptycs.AlertRuleCategory

	if len(alertRuleCategoryID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		alertRuleCategoryToLookup = uptycs.AlertRuleCategory{
			Name: alertRuleCategoryName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		alertRuleCategoryToLookup = uptycs.AlertRuleCategory{
			ID: alertRuleCategoryID,
		}
	}

	alertRuleCategoryResp, err := d.client.GetAlertRuleCategory(alertRuleCategoryToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get alertRuleCategory with ID  "+alertRuleCategoryID+": "+err.Error(),
		)
		return
	}

	var result = AlertRuleCategory{
		ID:     types.StringValue(alertRuleCategoryResp.ID),
		RuleID: types.StringValue(alertRuleCategoryResp.RuleID),
		Name:   types.StringValue(alertRuleCategoryResp.Name),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
