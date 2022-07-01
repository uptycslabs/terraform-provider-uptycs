package uptycs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceEventExcludeProfileType struct {
	p provider
}

func (r dataSourceEventExcludeProfileType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"platform": {
				Type:     types.StringType,
				Optional: true,
			},
			"metadata": {
				Optional: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

func (d dataSourceEventExcludeProfileType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceEventExcludeProfileType{
		p: *(p.(*provider)),
	}, nil
}

func (d dataSourceEventExcludeProfileType) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var eventExcludeProfileId string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &eventExcludeProfileId)...)

	eventExcludeProfileResp, err := d.p.client.GetEventExcludeProfile(uptycs.EventExcludeProfile{
		ID: eventExcludeProfileId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get eventExcludeProfile with ID  "+eventExcludeProfileId+": "+err.Error(),
		)
		return
	}

	metadataJson, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJson)) + "\n"},
		Priority:     eventExcludeProfileResp.Priority,
		ResourceType: types.String{Value: eventExcludeProfileResp.ResourceType},
		Platform:     types.String{Value: eventExcludeProfileResp.Platform},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
