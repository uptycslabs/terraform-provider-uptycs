package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &registryPathDataSource{}
	_ datasource.DataSourceWithConfigure = &registryPathDataSource{}
)

func RegistryPathDataSource() datasource.DataSource {
	return &registryPathDataSource{}
}

type registryPathDataSource struct {
	client *uptycs.Client
}

func (d *registryPathDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_path"
}

func (d *registryPathDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *registryPathDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"grouping": {
				Type:     types.StringType,
				Optional: true,
			},
			"include_registry_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"reg_accesses": {
				Type:     types.BoolType,
				Optional: true,
			},
			"exclude_registry_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"custom": {
				Type:     types.BoolType,
				Optional: true,
			},
		},
	}, nil
}

func (d *registryPathDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var registryPathID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &registryPathID)...)

	registryPathResp, err := d.client.GetRegistryPath(uptycs.RegistryPath{
		ID: registryPathID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get registryPath with ID  "+registryPathID+": "+err.Error(),
		)
		return
	}

	var result = RegistryPath{
		ID:          types.String{Value: registryPathResp.ID},
		Name:        types.String{Value: registryPathResp.Name},
		Description: types.String{Value: registryPathResp.Description},
		Grouping:    types.String{Value: registryPathResp.Grouping},
		IncludeRegistryPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		RegAccesses: types.Bool{Value: registryPathResp.RegAccesses},
		ExcludeRegistryPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom: types.Bool{Value: registryPathResp.Custom},
	}

	for _, _irp := range registryPathResp.IncludeRegistryPaths {
		result.IncludeRegistryPaths.Elems = append(result.IncludeRegistryPaths.Elems, types.String{Value: _irp})
	}

	for _, _erp := range registryPathResp.ExcludeRegistryPaths {
		result.ExcludeRegistryPaths.Elems = append(result.ExcludeRegistryPaths.Elems, types.String{Value: _erp})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
