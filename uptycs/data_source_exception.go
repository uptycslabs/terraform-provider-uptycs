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
	_ datasource.DataSource              = &exceptionDataSource{}
	_ datasource.DataSourceWithConfigure = &exceptionDataSource{}
)

func ExceptionDataSource() datasource.DataSource {
	return &exceptionDataSource{}
}

type exceptionDataSource struct {
	client *uptycs.Client
}

func (d *exceptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception"
}

func (d *exceptionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *exceptionDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"exception_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"table_name": {
				Type:     types.StringType,
				Optional: true,
			},
			"is_global": {
				Type:     types.BoolType,
				Optional: true,
			},
			"disabled": {
				Type:     types.BoolType,
				Optional: true,
			},
			"close_open_alerts": {
				Type:     types.BoolType,
				Optional: true,
			},
			"rule": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (d *exceptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var exceptionID string
	var exceptionName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &exceptionID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &exceptionName)

	var exceptionToLookup uptycs.Exception

	if len(exceptionID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		exceptionToLookup = uptycs.Exception{
			Name: exceptionName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		exceptionToLookup = uptycs.Exception{
			ID: exceptionID,
		}
	}
	exceptionResp, err := d.client.GetException(exceptionToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get exception with ID  "+exceptionID+": "+err.Error(),
		)
		return
	}

	ruleJSON, err := json.MarshalIndent(exceptionResp.Rule, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Exception{
		ID:              types.String{Value: exceptionResp.ID},
		Name:            types.String{Value: exceptionResp.Name},
		Description:     types.String{Value: exceptionResp.Description},
		ExceptionType:   types.String{Value: exceptionResp.ExceptionType},
		TableName:       types.String{Value: exceptionResp.TableName},
		IsGlobal:        types.Bool{Value: exceptionResp.IsGlobal},
		Disabled:        types.Bool{Value: exceptionResp.Disabled},
		CloseOpenAlerts: types.Bool{Value: exceptionResp.CloseOpenAlerts},
		Rule:            types.String{Value: string([]byte(ruleJSON)) + "\n"},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
