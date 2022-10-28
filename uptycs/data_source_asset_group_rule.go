package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &assetGroupRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &assetGroupRuleDataSource{}
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

func (d *assetGroupRuleDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"query": {
				Type:     types.StringType,
				Optional: true,
			},
			"interval": {
				Type:     types.NumberType,
				Optional: true,
			},
			"osquery_version": {
				Type:     types.StringType,
				Optional: true,
			},
			"platform": {
				Type:     types.StringType,
				Optional: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Optional: true,
			},
		},
	}, nil
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
		ID:             types.String{Value: assetGroupRuleResp.ID},
		Name:           types.String{Value: assetGroupRuleResp.Name},
		Description:    types.String{Value: assetGroupRuleResp.Description},
		Query:          types.String{Value: assetGroupRuleResp.Query},
		Interval:       assetGroupRuleResp.Interval,
		OsqueryVersion: types.String{Value: assetGroupRuleResp.OsqueryVersion},
		Platform:       types.String{Value: assetGroupRuleResp.Platform},
		Enabled:        types.Bool{Value: assetGroupRuleResp.Enabled},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
