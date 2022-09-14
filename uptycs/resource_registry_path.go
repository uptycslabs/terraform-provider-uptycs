package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceRegistryPathType struct{}

// Alert Rule Resource schema
func (r resourceRegistryPathType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
		},
	}, nil
}

// New resource instance
func (r resourceRegistryPathType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceRegistryPath{
		p: *(p.(*Provider)),
	}, nil
}

type resourceRegistryPath struct {
	p Provider
}

// Read resource information
func (r resourceRegistryPath) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var registryPathID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &registryPathID)...)
	registryPathResp, err := r.p.client.GetRegistryPath(uptycs.RegistryPath{
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

// Create a new resource
func (r resourceRegistryPath) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

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

	registryPathResp, err := r.p.client.CreateRegistryPath(uptycs.RegistryPath{
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Grouping:             plan.Grouping.Value,
		IncludeRegistryPaths: includeRegistryPaths,
		RegAccesses:          plan.RegAccesses.Value,
		ExcludeRegistryPaths: excludeRegistryPaths,
		Custom:               plan.Custom.Value,
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
		Custom: types.Bool{Value: registryPathResp.Custom},
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

// Update resource
func (r resourceRegistryPath) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	registryPathResp, err := r.p.client.UpdateRegistryPath(uptycs.RegistryPath{
		ID:                   registryPathID,
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Grouping:             plan.Grouping.Value,
		IncludeRegistryPaths: includeRegistryPaths,
		RegAccesses:          plan.RegAccesses.Value,
		ExcludeRegistryPaths: excludeRegistryPaths,
		Custom:               plan.Custom.Value,
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
		Custom: types.Bool{Value: registryPathResp.Custom},
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

// Delete resource
func (r resourceRegistryPath) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RegistryPath
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registryPathID := state.ID.Value

	_, err := r.p.client.DeleteRegistryPath(uptycs.RegistryPath{
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

// Import resource
func (r resourceRegistryPath) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
