package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &querypackDataSource{}
	_ datasource.DataSourceWithConfigure = &querypackDataSource{}
)

func QuerypackDataSource() datasource.DataSource {
	return &querypackDataSource{}
}

type querypackDataSource struct {
	client *uptycs.Client
}

func (d *querypackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_querypack"
}

func (d *querypackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *querypackDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"type": {
				Type:     types.StringType,
				Optional: true,
			},
			"additional_logger": {
				Type:     types.BoolType,
				Optional: true,
			},
			"custom": {
				Type:     types.BoolType,
				Optional: true,
			},
			"is_internal": {
				Type:     types.BoolType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"conf": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (d *querypackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var querypackID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &querypackID)...)

	querypackResp, err := d.client.GetQuerypack(uptycs.Querypack{
		ID: querypackID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get querypack with ID  "+querypackID+": "+err.Error(),
		)
		return
	}
	queryPackConfJSON, err := json.MarshalIndent(querypackResp.Conf, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Querypack{
		ID:               types.String{Value: querypackResp.ID},
		Name:             types.String{Value: querypackResp.Name},
		Description:      types.String{Value: querypackResp.Description},
		Type:             types.String{Value: querypackResp.Type},
		AdditionalLogger: types.Bool{Value: querypackResp.AdditionalLogger},
		Custom:           types.Bool{Value: querypackResp.Custom},
		IsInternal:       types.Bool{Value: querypackResp.IsInternal},
		ResourceType:     types.String{Value: querypackResp.ResourceType},
		Conf:             types.String{Value: string([]byte(queryPackConfJSON)) + "\n"},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
