package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceYaraGroupRuleType struct {
	p Provider
}

func (r dataSourceYaraGroupRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r dataSourceYaraGroupRuleType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceYaraGroupRuleType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceYaraGroupRuleType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var yaraGroupRuleID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &yaraGroupRuleID)...)

	yaraGroupRuleResp, err := r.p.client.GetYaraGroupRule(uptycs.YaraGroupRule{
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
