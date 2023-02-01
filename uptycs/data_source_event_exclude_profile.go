package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *eventExcludeProfileDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Optional: true},
			"name":          schema.StringAttribute{Optional: true},
			"description":   schema.StringAttribute{Optional: true},
			"priority":      schema.Int64Attribute{Optional: true},
			"resource_type": schema.StringAttribute{Optional: true},
			"platform":      schema.StringAttribute{Optional: true},
			"metadata":      schema.StringAttribute{Optional: true},
		},
	}
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
		ID:           types.StringValue(eventExcludeProfileResp.ID),
		Name:         types.StringValue(eventExcludeProfileResp.Name),
		Description:  types.StringValue(eventExcludeProfileResp.Description),
		Metadata:     types.StringValue(string(metadataJSON) + "\n"),
		Priority:     types.Int64Value(int64(eventExcludeProfileResp.Priority)),
		ResourceType: types.StringValue(eventExcludeProfileResp.ResourceType),
		Platform:     types.StringValue(eventExcludeProfileResp.Platform),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
