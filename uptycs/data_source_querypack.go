package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceQuerypackType struct {
	p Provider
}

func (r dataSourceQuerypackType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
			},
			"sha": {
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
			"queries": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							Computed: true,
							Type:     types.StringType,
						},
						"name": {
							Type:     types.StringType,
							Optional: true,
						},
						"description": {
							Type:     types.StringType,
							Optional: true,
						},
					},
				),
			},
			"conf": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (r dataSourceQuerypackType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceQuerypackType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceQuerypackType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var querypackID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &querypackID)...)

	querypackResp, err := r.p.client.GetQuerypack(uptycs.Querypack{
		ID: querypackID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get querypack with ID  "+querypackID+": "+err.Error(),
		)
		return
	}
	queryPackConfJson, err := json.MarshalIndent(querypackResp.Conf, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	queries := Query{}
	querypackResp.Queries.(types.Object).As(ctx, &queries, types.ObjectAsOptions{})

	var result = Querypack{
		ID:               types.String{Value: querypackResp.ID},
		Sha:              types.String{Value: querypackResp.Sha},
		Name:             types.String{Value: querypackResp.Name},
		Description:      types.String{Value: querypackResp.Description},
		Type:             types.String{Value: querypackResp.Type},
		AdditionalLogger: types.Bool{Value: querypackResp.AdditionalLogger},
		Custom:           types.Bool{Value: querypackResp.Custom},
		IsInternal:       types.Bool{Value: querypackResp.IsInternal},
		ResourceType:     types.String{Value: querypackResp.ResourceType},
		//Queries: types.List{
		//	ElemType: types.StringType,
		//	Elems:    make([]attr.Value, 0),
		//},
		Conf: types.String{Value: string([]byte(queryPackConfJson)) + "\n"},
	}
	panic("fuck")

	//for _, q := range querypackResp.Queries {
	//	result.Queries.Elems = append(result.Queries.Elems, types.String{Value: "foo"})
	//}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
