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
	_ datasource.DataSource              = &atcQueryDataSource{}
	_ datasource.DataSourceWithConfigure = &atcQueryDataSource{}
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

func (d *atcQueryDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
		},
	}, nil
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
		ID:          types.String{Value: atcQueryResp.ID},
		Name:        types.String{Value: atcQueryResp.Name},
		Description: types.String{Value: atcQueryResp.Description},
		Query:       types.String{Value: atcQueryResp.Query},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
