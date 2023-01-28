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

func (d *customProfileDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Optional: true},
			"name":            schema.StringAttribute{Optional: true},
			"description":     schema.StringAttribute{Optional: true},
			"query_schedules": schema.StringAttribute{Optional: true},
			"priority":        schema.NumberAttribute{Optional: true},
			"resource_type":   schema.StringAttribute{Optional: true},
		},
	}
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
		ID:             types.StringValue(customProfileResp.ID),
		Name:           types.StringValue(customProfileResp.Name),
		Description:    types.StringValue(customProfileResp.Description),
		QuerySchedules: types.StringValue(string([]byte(queryScheduleJSON)) + "\n"),
		Priority:       customProfileResp.Priority,
		ResourceType:   types.StringValue(customProfileResp.ResourceType),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
