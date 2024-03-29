package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func AtcQueryDataSource() datasource.DataSource {
	return &atcQueryDataSource{}
}

type atcQueryDataSource struct {
	client *uptycs.Client
}

func (d *atcQueryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_atc_query"
}

func (d *atcQueryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *atcQueryDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"query":       schema.StringAttribute{Optional: true},
		},
	}
}

func (d *atcQueryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var atcQueryID string
	var atcQueryName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &atcQueryID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &atcQueryName)

	var atcQueryToLookup uptycs.AtcQuery

	if len(atcQueryID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		atcQueryToLookup = uptycs.AtcQuery{
			Name: atcQueryName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		atcQueryToLookup = uptycs.AtcQuery{
			ID: atcQueryID,
		}
	}

	atcQueryResp, err := d.client.GetAtcQuery(atcQueryToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get atcQuery with ID  "+atcQueryID+": "+err.Error(),
		)
		return
	}

	var result = AtcQuery{
		ID:          types.StringValue(atcQueryResp.ID),
		Name:        types.StringValue(atcQueryResp.Name),
		Description: types.StringValue(atcQueryResp.Description),
		Query:       types.StringValue(atcQueryResp.Query),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
