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
	_ datasource.DataSource              = &yaraGroupRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &yaraGroupRuleDataSource{}
)

func YaraGroupRuleDataSource() datasource.DataSource {
	return &yaraGroupRuleDataSource{}
}

type yaraGroupRuleDataSource struct {
	client *uptycs.Client
}

func (d *yaraGroupRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_yara_group_rule"
}

func (d *yaraGroupRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *yaraGroupRuleDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"rules": {
				Type:     types.StringType,
				Optional: true,
			},
			"custom": {
				Type:     types.BoolType,
				Optional: true,
			},
		},
	}, nil
}

func (d *yaraGroupRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var yaraGroupRuleID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &yaraGroupRuleID)...)

	yaraGroupRuleResp, err := d.client.GetYaraGroupRule(uptycs.YaraGroupRule{
		ID: yaraGroupRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get yaraGroupRule with ID  "+yaraGroupRuleID+": "+err.Error(),
		)
		return
	}

	var result = YaraGroupRule{
		ID:          types.String{Value: yaraGroupRuleResp.ID},
		Name:        types.String{Value: yaraGroupRuleResp.Name},
		Description: types.String{Value: yaraGroupRuleResp.Description},
		Rules:       types.String{Value: yaraGroupRuleResp.Rules},
		Custom:      types.Bool{Value: yaraGroupRuleResp.Custom},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
