package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *registryPathDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"grouping":    schema.StringAttribute{Optional: true},
			"include_registry_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"reg_accesses": schema.BoolAttribute{Optional: true},
			"exclude_registry_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
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
		ID:          types.StringValue(registryPathResp.ID),
		Name:        types.StringValue(registryPathResp.Name),
		Description: types.StringValue(registryPathResp.Description),
		Grouping:    types.StringValue(registryPathResp.Grouping),
		IncludeRegistryPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		RegAccesses: types.BoolValue(registryPathResp.RegAccesses),
		ExcludeRegistryPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
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
