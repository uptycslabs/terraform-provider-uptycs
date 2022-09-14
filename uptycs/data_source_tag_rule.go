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

type dataSourceTagRuleType struct {
	p Provider
}

func (r dataSourceTagRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"source": {
				Type:     types.StringType,
				Optional: true,
			},
			"run_once": {
				Type:     types.BoolType,
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
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Optional: true,
			},
			"system": {
				Type:     types.BoolType,
				Optional: true,
			},
		},
	}, nil
}

func (r dataSourceTagRuleType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceTagRuleType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceTagRuleType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tagRuleID string
	var tagRuleName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &tagRuleID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &tagRuleName)

	var tagRuleToLookup uptycs.TagRule

	if len(tagRuleID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		tagRuleToLookup = uptycs.TagRule{
			Name: tagRuleName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		tagRuleToLookup = uptycs.TagRule{
			ID: tagRuleID,
		}
	}

	tagRuleResp, err := r.p.client.GetTagRule(tagRuleToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get tagRule with ID  "+tagRuleID+": "+err.Error(),
		)
		return
	}
	var result = TagRule{
		ID:             types.String{Value: tagRuleResp.ID},
		Name:           types.String{Value: tagRuleResp.Name},
		Description:    types.String{Value: tagRuleResp.Description},
		Query:          types.String{Value: tagRuleResp.Query},
		Source:         types.String{Value: tagRuleResp.Source},
		RunOnce:        types.Bool{Value: tagRuleResp.RunOnce},
		Interval:       tagRuleResp.Interval,
		OSqueryVersion: types.String{Value: tagRuleResp.OSqueryVersion},
		Platform:       types.String{Value: tagRuleResp.Name},
		Enabled:        types.Bool{Value: tagRuleResp.Enabled},
		System:         types.Bool{Value: tagRuleResp.System},
		ResourceType:   types.String{Value: tagRuleResp.ResourceType},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
