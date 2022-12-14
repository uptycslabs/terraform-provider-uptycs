package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &flagProfileResource{}
	_ resource.ResourceWithConfigure   = &flagProfileResource{}
	_ resource.ResourceWithImportState = &flagProfileResource{}
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

func (r *flagProfileResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Type:     types.StringType,
				Optional: true,
			},
			"flags": {
				Type:     types.StringType,
				Required: true,
			},
			"os_flags": {
				Type:     types.StringType,
				Required: true,
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
		Priority:     flagProfileResp.Priority,
		Flags:        types.StringValue(string([]byte(flagsJSON)) + "\n"),
		OsFlags:      types.StringValue(string([]byte(osFlagsJSON)) + "\n"),
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
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		Flags:        uptycs.CustomJSONString(plan.Flags.Value),
		OsFlags:      uptycs.CustomJSONString(plan.OsFlags.Value),
		Priority:     plan.Priority,
		ResourceType: plan.ResourceType.Value,
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
		Flags:        types.StringValue(string([]byte(flagsJSON)) + "\n"),
		OsFlags:      types.StringValue(string([]byte(osFlagsJSON)) + "\n"),
		Priority:     flagProfileResp.Priority,
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

	flagProfileID := state.ID.Value

	// Retrieve values from plan
	var plan FlagProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flagProfileResp, err := r.client.UpdateFlagProfile(uptycs.FlagProfile{
		ID:           flagProfileID,
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		Flags:        uptycs.CustomJSONString(plan.Flags.Value),
		OsFlags:      uptycs.CustomJSONString(plan.OsFlags.Value),
		Priority:     plan.Priority,
		ResourceType: plan.ResourceType.Value,
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
		Flags:        types.StringValue(string([]byte(flagsJSON)) + "\n"),
		OsFlags:      types.StringValue(string([]byte(osFlagsJSON)) + "\n"),
		Priority:     flagProfileResp.Priority,
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

	flagProfileID := state.ID.Value

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
