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
	_ datasource.DataSource              = &customProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &customProfileDataSource{}
)

func CustomProfileDataSource() datasource.DataSource {
	return &customProfileDataSource{}
}

type customProfileDataSource struct {
	client *uptycs.Client
}

func (d *customProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_profile"
}

func (d *customProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *customProfileDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"query_schedules": {
				Type:     types.StringType,
				Optional: true,
			},
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (d *customProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var customProfileID string
	var customProfileName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &customProfileID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &customProfileName)

	var customProfileToLookup uptycs.CustomProfile

	if len(customProfileID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		customProfileToLookup = uptycs.CustomProfile{
			Name: customProfileName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		customProfileToLookup = uptycs.CustomProfile{
			ID: customProfileID,
		}
	}

	customProfileResp, err := d.client.GetCustomProfile(customProfileToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get customProfile with ID  "+customProfileID+": "+err.Error(),
		)
		return
	}

	queryScheduleJSON, err := json.MarshalIndent(customProfileResp.QuerySchedules, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = CustomProfile{
		ID:             types.String{Value: customProfileResp.ID},
		Name:           types.String{Value: customProfileResp.Name},
		Description:    types.String{Value: customProfileResp.Description},
		QuerySchedules: types.String{Value: string([]byte(queryScheduleJSON)) + "\n"},
		Priority:       customProfileResp.Priority,
		ResourceType:   types.String{Value: customProfileResp.ResourceType},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
