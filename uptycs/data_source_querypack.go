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
				Optional: true,
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
						"query": {
							Type:     types.StringType,
							Optional: true,
						},
						"removed": {
							Type:     types.BoolType,
							Optional: true,
						},
						"version": {
							Type:     types.StringType,
							Optional: true,
						},
						"interval": {
							Type:     types.NumberType,
							Optional: true,
						},
						"platform": {
							Type:     types.StringType,
							Optional: true,
						},
						"snapshot": {
							Type:     types.BoolType,
							Optional: true,
						},
						"run_now": {
							Type:     types.BoolType,
							Optional: true,
						},
						"value": {
							Type:     types.StringType,
							Optional: true,
						},
						"querypack_id": {
							Type:     types.StringType,
							Optional: true,
						},
						"table_name": {
							Type:     types.StringType,
							Optional: true,
						},
						"data_types": {
							Type:     types.StringType,
							Optional: true,
						},
						"verified": {
							Type:     types.BoolType,
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

	var queries []Query
	for _, q := range querypackResp.Queries {
		dataTypesJSON, err := json.MarshalIndent(q.DataTypes, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		queries = append(queries, Query{
			ID:          types.String{Value: q.ID},
			Name:        types.String{Value: q.Name},
			Description: types.String{Value: q.Description},
			Query:       types.String{Value: q.Query},
			Removed:     types.Bool{Value: q.Removed},
			Version:     types.String{Value: q.Version},
			Interval:    q.Interval,
			Platform:    types.String{Value: q.Platform},
			Snapshot:    types.Bool{Value: q.Snapshot},
			RunNow:      types.Bool{Value: q.RunNow},
			Value:       types.String{Value: q.Value},
			QuerypackID: types.String{Value: q.QuerypackID},
			TableName:   types.String{Value: q.TableName},
			DataTypes:   types.String{Value: string([]byte(dataTypesJSON)) + "\n"},
			Verified:    types.Bool{Value: q.Verified},
		})
	}
	result.Queries = queries

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
