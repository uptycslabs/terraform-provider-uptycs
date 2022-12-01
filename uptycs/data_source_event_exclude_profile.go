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
	_ datasource.DataSource              = &eventExcludeProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &eventExcludeProfileDataSource{}
)

func EventExcludeProfileDataSource() datasource.DataSource {
	return &eventExcludeProfileDataSource{}
}

type eventExcludeProfileDataSource struct {
	client *uptycs.Client
}

func (d *eventExcludeProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_exclude_profile"
}

func (d *eventExcludeProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *eventExcludeProfileDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *eventExcludeProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var eventExcludeProfileID string
	var eventExcludeProfileName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &eventExcludeProfileID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &eventExcludeProfileName)

	var eventExcludeProfileToLookup uptycs.EventExcludeProfile

	if len(eventExcludeProfileID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		eventExcludeProfileToLookup = uptycs.EventExcludeProfile{
			Name: eventExcludeProfileName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		eventExcludeProfileToLookup = uptycs.EventExcludeProfile{
			ID: eventExcludeProfileID,
		}
	}

	eventExcludeProfileResp, err := d.client.GetEventExcludeProfile(eventExcludeProfileToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get eventExcludeProfile with ID  "+eventExcludeProfileID+": "+err.Error(),
		)
		return
	}

	metadataJSON, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJSON)) + "\n"},
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
