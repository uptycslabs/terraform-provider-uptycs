package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &registryPathResource{}
	_ resource.ResourceWithConfigure   = &registryPathResource{}
	_ resource.ResourceWithImportState = &registryPathResource{}
)

func RegistryPathResource() resource.Resource {
	return &registryPathResource{}
}

type registryPathResource struct {
	client *uptycs.Client
}

func (r *registryPathResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_path"
}

func (r *registryPathResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *registryPathResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"grouping": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"include_registry_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"reg_accesses": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"exclude_registry_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
		},
	}, nil
}

func (r *registryPathResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var registryPathID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &registryPathID)...)
	registryPathResp, err := r.client.GetRegistryPath(uptycs.RegistryPath{
		ID: registryPathID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
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

func (r *registryPathResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RegistryPath
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var includeRegistryPaths []string
	plan.IncludeRegistryPaths.ElementsAs(ctx, &includeRegistryPaths, false)

	var excludeRegistryPaths []string
	plan.ExcludeRegistryPaths.ElementsAs(ctx, &excludeRegistryPaths, false)

	registryPathResp, err := r.client.CreateRegistryPath(uptycs.RegistryPath{
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Grouping:             plan.Grouping.Value,
		IncludeRegistryPaths: includeRegistryPaths,
		RegAccesses:          plan.RegAccesses.Value,
		ExcludeRegistryPaths: excludeRegistryPaths,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create registryPath, unexpected error: "+err.Error(),
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
	}

	for _, _irp := range registryPathResp.IncludeRegistryPaths {
		result.IncludeRegistryPaths.Elems = append(result.IncludeRegistryPaths.Elems, types.String{Value: _irp})
	}

	for _, _erp := range registryPathResp.ExcludeRegistryPaths {
		result.ExcludeRegistryPaths.Elems = append(result.ExcludeRegistryPaths.Elems, types.String{Value: _erp})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *registryPathResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state RegistryPath
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registryPathID := state.ID.Value

	// Retrieve values from plan
	var plan RegistryPath
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var includeRegistryPaths []string
	plan.IncludeRegistryPaths.ElementsAs(ctx, &includeRegistryPaths, false)

	var excludeRegistryPaths []string
	plan.ExcludeRegistryPaths.ElementsAs(ctx, &excludeRegistryPaths, false)

	registryPathResp, err := r.client.UpdateRegistryPath(uptycs.RegistryPath{
		ID:                   registryPathID,
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Grouping:             plan.Grouping.Value,
		IncludeRegistryPaths: includeRegistryPaths,
		RegAccesses:          plan.RegAccesses.Value,
		ExcludeRegistryPaths: excludeRegistryPaths,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create registryPath, unexpected error: "+err.Error(),
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
	}

	for _, _irp := range registryPathResp.IncludeRegistryPaths {
		result.IncludeRegistryPaths.Elems = append(result.IncludeRegistryPaths.Elems, types.String{Value: _irp})
	}

	for _, _erp := range registryPathResp.ExcludeRegistryPaths {
		result.ExcludeRegistryPaths.Elems = append(result.ExcludeRegistryPaths.Elems, types.String{Value: _erp})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *registryPathResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RegistryPath
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registryPathID := state.ID.Value

	_, err := r.client.DeleteRegistryPath(uptycs.RegistryPath{
		ID: registryPathID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete registryPath with ID  "+registryPathID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r *registryPathResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
