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

func FlagProfileDataSource() datasource.DataSource {
	return &flagProfileDataSource{}
}

type flagProfileDataSource struct {
	client *uptycs.Client
}

func (d *flagProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flag_profile"
}

func (d *flagProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *flagProfileDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Optional: true},
			"name":          schema.StringAttribute{Optional: true},
			"description":   schema.StringAttribute{Optional: true},
			"flags":         schema.StringAttribute{Optional: true},
			"os_flags":      schema.StringAttribute{Optional: true},
			"resource_type": schema.StringAttribute{Optional: true},
			"priority":      schema.NumberAttribute{Optional: true},
		},
	}
}

func (d *flagProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var flagProfileID string
	var flagProfileName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &flagProfileID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &flagProfileName)

	var flagProfileToLookup uptycs.FlagProfile

	if len(flagProfileID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		flagProfileToLookup = uptycs.FlagProfile{
			Name: flagProfileName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		flagProfileToLookup = uptycs.FlagProfile{
			ID: flagProfileID,
		}
	}

	flagProfileResp, err := d.client.GetFlagProfile(flagProfileToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get flagProfile with ID  "+flagProfileID+": "+err.Error(),
		)
		return
	}

	flagsJSON, err := json.MarshalIndent(flagProfileResp.Flags, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	osFlagsJSON, err := json.MarshalIndent(flagProfileResp.OsFlags, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = FlagProfile{
		ID:           types.StringValue(flagProfileResp.ID),
		Name:         types.StringValue(flagProfileResp.Name),
		Description:  types.StringValue(flagProfileResp.Description),
		Priority:     flagProfileResp.Priority,
		Flags:        types.StringValue(string(flagsJSON) + "\n"),
		OsFlags:      types.StringValue(string(osFlagsJSON) + "\n"),
		ResourceType: types.StringValue(flagProfileResp.ResourceType),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
