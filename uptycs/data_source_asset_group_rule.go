package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func AssetGroupRuleDataSource() datasource.DataSource {
	return &assetGroupRuleDataSource{}
}

type assetGroupRuleDataSource struct {
	client *uptycs.Client
}

func (d *assetGroupRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_group_rule"
}

func (d *assetGroupRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *assetGroupRuleDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Optional: true},
			"name":            schema.StringAttribute{Optional: true},
			"description":     schema.StringAttribute{Optional: true},
			"query":           schema.StringAttribute{Optional: true},
			"interval":        schema.Int64Attribute{Optional: true},
			"osquery_version": schema.StringAttribute{Optional: true},
			"platform":        schema.StringAttribute{Optional: true},
			"enabled":         schema.BoolAttribute{Optional: true},
		},
	}
}

func (d *assetGroupRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var assetGroupRuleID string
	var assetGroupRuleName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &assetGroupRuleID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &assetGroupRuleName)

	var assetGroupRuleToLookup uptycs.AssetGroupRule

	if len(assetGroupRuleID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		assetGroupRuleToLookup = uptycs.AssetGroupRule{
			Name: assetGroupRuleName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		assetGroupRuleToLookup = uptycs.AssetGroupRule{
			ID: assetGroupRuleID,
		}
	}

	assetGroupRuleResp, err := d.client.GetAssetGroupRule(assetGroupRuleToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get assetGroupRule with ID  "+assetGroupRuleID+": "+err.Error(),
		)
		return
	}

	var result = AssetGroupRule{
		ID:             types.StringValue(assetGroupRuleResp.ID),
		Name:           types.StringValue(assetGroupRuleResp.Name),
		Description:    types.StringValue(assetGroupRuleResp.Description),
		Query:          types.StringValue(assetGroupRuleResp.Query),
		Interval:       types.Int64Value(int64(assetGroupRuleResp.Interval)),
		OsqueryVersion: types.StringValue(assetGroupRuleResp.OsqueryVersion),
		Platform:       types.StringValue(assetGroupRuleResp.Platform),
		Enabled:        types.BoolValue(assetGroupRuleResp.Enabled),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
