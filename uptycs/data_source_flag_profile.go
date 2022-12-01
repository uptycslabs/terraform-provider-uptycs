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
	_ datasource.DataSource              = &flagProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &flagProfileDataSource{}
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

func (d *flagProfileDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"flags": {
				Type:     types.StringType,
				Optional: true,
			},
			"os_flags": {
				Type:     types.StringType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
		},
	}, nil
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
		ID:           types.String{Value: flagProfileResp.ID},
		Name:         types.String{Value: flagProfileResp.Name},
		Description:  types.String{Value: flagProfileResp.Description},
		Priority:     flagProfileResp.Priority,
		Flags:        types.String{Value: string([]byte(flagsJSON)) + "\n"},
		OsFlags:      types.String{Value: string([]byte(osFlagsJSON)) + "\n"},
		ResourceType: types.String{Value: flagProfileResp.ResourceType},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
