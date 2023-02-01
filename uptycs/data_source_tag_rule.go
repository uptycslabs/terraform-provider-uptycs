package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func TagRuleDataSource() datasource.DataSource {
	return &tagRuleDataSource{}
}

type tagRuleDataSource struct {
	client *uptycs.Client
}

func (d *tagRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag_rule"
}

func (d *tagRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *tagRuleDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Optional: true},
			"name":            schema.StringAttribute{Optional: true},
			"description":     schema.StringAttribute{Optional: true},
			"query":           schema.StringAttribute{Optional: true},
			"source":          schema.StringAttribute{Optional: true},
			"run_once":        schema.BoolAttribute{Optional: true},
			"interval":        schema.Int64Attribute{Optional: true},
			"osquery_version": schema.StringAttribute{Optional: true},
			"platform":        schema.StringAttribute{Optional: true},
			"resource_type":   schema.StringAttribute{Optional: true},
			"enabled":         schema.BoolAttribute{Optional: true},
			"system":          schema.BoolAttribute{Optional: true},
		},
	}
}

func (d *tagRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	tagRuleResp, err := d.client.GetTagRule(tagRuleToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get tagRule with ID  "+tagRuleID+": "+err.Error(),
		)
		return
	}
	var result = TagRule{
		ID:             types.StringValue(tagRuleResp.ID),
		Name:           types.StringValue(tagRuleResp.Name),
		Description:    types.StringValue(tagRuleResp.Description),
		Query:          types.StringValue(tagRuleResp.Query),
		Source:         types.StringValue(tagRuleResp.Source),
		RunOnce:        types.BoolValue(tagRuleResp.RunOnce),
		Interval:       types.Int64Value(int64(tagRuleResp.Interval)),
		OSqueryVersion: types.StringValue(tagRuleResp.OSqueryVersion),
		Platform:       types.StringValue(tagRuleResp.Name),
		Enabled:        types.BoolValue(tagRuleResp.Enabled),
		System:         types.BoolValue(tagRuleResp.System),
		ResourceType:   types.StringValue(tagRuleResp.ResourceType),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
