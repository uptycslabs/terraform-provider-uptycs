package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *registryPathResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringDefault(""),
				},
			},
			"grouping": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringDefault(""),
				},
			},
			"include_registry_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"reg_accesses": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolDefault(false),
				},
			},
			"exclude_registry_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
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
		ID:                   types.StringValue(registryPathResp.ID),
		Name:                 types.StringValue(registryPathResp.Name),
		Description:          types.StringValue(registryPathResp.Description),
		Grouping:             types.StringValue(registryPathResp.Grouping),
		IncludeRegistryPaths: makeListStringAttribute(registryPathResp.IncludeRegistryPaths),
		RegAccesses:          types.BoolValue(registryPathResp.RegAccesses),
		ExcludeRegistryPaths: makeListStringAttribute(registryPathResp.ExcludeRegistryPaths),
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
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Grouping:             plan.Grouping.ValueString(),
		IncludeRegistryPaths: includeRegistryPaths,
		RegAccesses:          plan.RegAccesses.ValueBool(),
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
		ID:                   types.StringValue(registryPathResp.ID),
		Name:                 types.StringValue(registryPathResp.Name),
		Description:          types.StringValue(registryPathResp.Description),
		Grouping:             types.StringValue(registryPathResp.Grouping),
		IncludeRegistryPaths: makeListStringAttribute(registryPathResp.IncludeRegistryPaths),
		RegAccesses:          types.BoolValue(registryPathResp.RegAccesses),
		ExcludeRegistryPaths: makeListStringAttribute(registryPathResp.ExcludeRegistryPaths),
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

	registryPathID := state.ID.ValueString()

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
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Grouping:             plan.Grouping.ValueString(),
		IncludeRegistryPaths: includeRegistryPaths,
		RegAccesses:          plan.RegAccesses.ValueBool(),
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
		ID:                   types.StringValue(registryPathResp.ID),
		Name:                 types.StringValue(registryPathResp.Name),
		Description:          types.StringValue(registryPathResp.Description),
		Grouping:             types.StringValue(registryPathResp.Grouping),
		IncludeRegistryPaths: makeListStringAttribute(registryPathResp.IncludeRegistryPaths),
		RegAccesses:          types.BoolValue(registryPathResp.RegAccesses),
		ExcludeRegistryPaths: makeListStringAttribute(registryPathResp.ExcludeRegistryPaths),
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

	registryPathID := state.ID.ValueString()

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
