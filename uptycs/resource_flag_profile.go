package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func FlagProfileResource() resource.Resource {
	return &flagProfileResource{}
}

type flagProfileResource struct {
	client *uptycs.Client
}

func (r *flagProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flag_profile"
}

func (r *flagProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *flagProfileResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"name":          schema.StringAttribute{Optional: true},
			"description":   schema.StringAttribute{Optional: true},
			"flags":         schema.StringAttribute{Required: true},
			"os_flags":      schema.StringAttribute{Required: true},
			"resource_type": schema.StringAttribute{Optional: true},
			"priority":      schema.Int64Attribute{Optional: true},
		},
	}
}

func (r *flagProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var flagProfileID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &flagProfileID)...)
	flagProfileResp, err := r.client.GetFlagProfile(uptycs.FlagProfile{
		ID: flagProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
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
		Priority:     types.Int64Value(int64(flagProfileResp.Priority)),
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

func (r *flagProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan FlagProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flagProfileResp, err := r.client.CreateFlagProfile(uptycs.FlagProfile{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		Flags:        uptycs.CustomJSONString(plan.Flags.ValueString()),
		OsFlags:      uptycs.CustomJSONString(plan.OsFlags.ValueString()),
		Priority:     int(plan.Priority.ValueInt64()),
		ResourceType: plan.ResourceType.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create flagProfile, unexpected error: "+err.Error(),
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
		Flags:        types.StringValue(string(flagsJSON) + "\n"),
		OsFlags:      types.StringValue(string(osFlagsJSON) + "\n"),
		Priority:     types.Int64Value(int64(flagProfileResp.Priority)),
		ResourceType: types.StringValue(flagProfileResp.ResourceType),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *flagProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state FlagProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flagProfileID := state.ID.ValueString()

	// Retrieve values from plan
	var plan FlagProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flagProfileResp, err := r.client.UpdateFlagProfile(uptycs.FlagProfile{
		ID:           flagProfileID,
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		Flags:        uptycs.CustomJSONString(plan.Flags.ValueString()),
		OsFlags:      uptycs.CustomJSONString(plan.OsFlags.ValueString()),
		Priority:     int(plan.Priority.ValueInt64()),
		ResourceType: plan.ResourceType.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create flagProfile, unexpected error: "+err.Error(),
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
		Flags:        types.StringValue(string(flagsJSON) + "\n"),
		OsFlags:      types.StringValue(string(osFlagsJSON) + "\n"),
		Priority:     types.Int64Value(int64(flagProfileResp.Priority)),
		ResourceType: types.StringValue(flagProfileResp.ResourceType),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *flagProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FlagProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flagProfileID := state.ID.ValueString()

	_, err := r.client.DeleteFlagProfile(uptycs.FlagProfile{
		ID: flagProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete flagProfile with ID  "+flagProfileID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *flagProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
