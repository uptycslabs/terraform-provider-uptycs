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

type resourceFilePathGroupType struct{}

func (r resourceFilePathGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
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
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"include_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"include_path_extensions": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"exclude_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"check_signature": {
				Type:     types.BoolType,
				Optional: true,
			},
			"file_accesses": {
				Type:     types.BoolType,
				Optional: true,
			},
			"exclude_process_names": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"priority_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"signatures": {
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							Computed: true,
							Optional: true,
							Type:     types.StringType,
						},
						"name": {
							Type:     types.StringType,
							Optional: true,
						},
						"description": {
							Type:     types.StringType,
							Optional: true,
						},
						"paths": {
							Type:     types.ListType{ElemType: types.StringType},
							Optional: true,
						},
					},
				),
			},
			"yara_group_rules": {
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							Computed: true,
							Optional: true,
							Type:     types.StringType,
						},
						"name": {
							Type:     types.StringType,
							Computed: true,
							Optional: true,
						},
						"description": {
							Type:     types.StringType,
							Computed: true,
							Optional: true,
						},
						"rules": {
							Computed: true,
							Type:     types.StringType,
							Optional: true,
						},
						"custom": {
							Type:     types.BoolType,
							Optional: true,
							Computed: true,
						},
					},
				),
			},
		},
	}, nil
}

// New resource instance
func (r resourceFilePathGroupType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceFilePathGroup{
		p: *(p.(*Provider)),
	}, nil
}

type resourceFilePathGroup struct {
	p Provider
}

// Read resource information
func (r resourceFilePathGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var filePathGroupID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &filePathGroupID)...)
	filePathGroupResp, err := r.p.client.GetFilePathGroup(uptycs.FilePathGroup{
		ID: filePathGroupID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get filePathGroup with ID  "+filePathGroupID+": "+err.Error(),
		)
		return
	}

	var result = FilePathGroup{
		ID:          types.String{Value: filePathGroupResp.ID},
		Name:        types.String{Value: filePathGroupResp.Name},
		Description: types.String{Value: filePathGroupResp.Description},
		Grouping:    types.String{Value: filePathGroupResp.Grouping},
		IncludePaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		IncludePathExtensions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		ExcludePaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom:         types.Bool{Value: filePathGroupResp.Custom},
		CheckSignature: types.Bool{Value: filePathGroupResp.CheckSignature},
		FileAccesses:   types.Bool{Value: filePathGroupResp.FileAccesses},
		ExcludeProcessNames: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		PriorityPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, ip := range filePathGroupResp.IncludePaths {
		result.IncludePaths.Elems = append(result.IncludePaths.Elems, types.String{Value: ip})
	}

	for _, ipe := range filePathGroupResp.IncludePathExtensions {
		result.IncludePathExtensions.Elems = append(result.IncludePathExtensions.Elems, types.String{Value: ipe})
	}

	for _, ep := range filePathGroupResp.ExcludePaths {
		result.ExcludePaths.Elems = append(result.ExcludePaths.Elems, types.String{Value: ep})
	}

	for _, epn := range filePathGroupResp.ExcludeProcessNames {
		result.ExcludeProcessNames.Elems = append(result.ExcludeProcessNames.Elems, types.String{Value: epn})
	}

	for _, pp := range filePathGroupResp.PriorityPaths {
		result.PriorityPaths.Elems = append(result.PriorityPaths.Elems, types.String{Value: pp})
	}

	signatures := make([]FilePathGroupSignature, 0)
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.String{Value: s.ID},
			Name:        types.String{Value: s.Name},
			Description: types.String{Value: s.Description},
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	yaraGroupRules := make([]YaraGroupRule, 0)
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.String{Value: ygr.ID},
			Name:        types.String{Value: ygr.Name},
			Description: types.String{Value: ygr.Description},
			Rules:       types.String{Value: ygr.Rules},
			Custom:      types.Bool{Value: ygr.Custom},
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Create a new resource
func (r resourceFilePathGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan FilePathGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var includePaths []string
	plan.IncludePaths.ElementsAs(ctx, &includePaths, false)

	var includePathExtensions []string
	plan.IncludePathExtensions.ElementsAs(ctx, &includePathExtensions, false)

	var excludePaths []string
	plan.ExcludePaths.ElementsAs(ctx, &excludePaths, false)

	var excludeProcessNames []string
	plan.ExcludeProcessNames.ElementsAs(ctx, &excludeProcessNames, false)

	var priorityPaths []string
	plan.PriorityPaths.ElementsAs(ctx, &priorityPaths, false)

	_signatures := make([]uptycs.FilePathGroupSignature, 0)
	for _, s := range plan.Signatures {
		_signatures = append(_signatures, uptycs.FilePathGroupSignature{
			ID:   s.ID.Value,
			Name: s.Name.Value,
		})
	}

	_yaraGroupRules := make([]uptycs.YaraGroupRule, 0)
	for _, yg := range plan.YaraGroupRules {
		_yaraGroupRules = append(_yaraGroupRules, uptycs.YaraGroupRule{
			ID:   yg.ID.Value,
			Name: yg.Name.Value,
		})
	}

	filePathGroupResp, err := r.p.client.CreateFilePathGroup(uptycs.FilePathGroup{
		Name:                  plan.Name.Value,
		Description:           plan.Description.Value,
		Grouping:              plan.Grouping.Value,
		IncludePaths:          includePaths,
		IncludePathExtensions: includePathExtensions,
		ExcludePaths:          excludePaths,
		Custom:                plan.Custom.Value,
		CheckSignature:        plan.CheckSignature.Value,
		FileAccesses:          plan.FileAccesses.Value,
		ExcludeProcessNames:   excludeProcessNames,
		PriorityPaths:         priorityPaths,
		Signatures:            _signatures,
		YaraGroupRules:        _yaraGroupRules,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create filePathGroup, unexpected error: "+err.Error(),
		)
		return
	}

	var result = FilePathGroup{
		ID:          types.String{Value: filePathGroupResp.ID},
		Name:        types.String{Value: filePathGroupResp.Name},
		Description: types.String{Value: filePathGroupResp.Description},
		Grouping:    types.String{Value: filePathGroupResp.Grouping},
		IncludePaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		IncludePathExtensions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		ExcludePaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom:         types.Bool{Value: filePathGroupResp.Custom},
		CheckSignature: types.Bool{Value: filePathGroupResp.CheckSignature},
		FileAccesses:   types.Bool{Value: filePathGroupResp.FileAccesses},
		ExcludeProcessNames: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		PriorityPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, ip := range filePathGroupResp.IncludePaths {
		result.IncludePaths.Elems = append(result.IncludePaths.Elems, types.String{Value: ip})
	}

	for _, ipe := range filePathGroupResp.IncludePathExtensions {
		result.IncludePathExtensions.Elems = append(result.IncludePathExtensions.Elems, types.String{Value: ipe})
	}

	for _, ep := range filePathGroupResp.ExcludePaths {
		result.ExcludePaths.Elems = append(result.ExcludePaths.Elems, types.String{Value: ep})
	}

	for _, epn := range filePathGroupResp.ExcludeProcessNames {
		result.ExcludeProcessNames.Elems = append(result.ExcludeProcessNames.Elems, types.String{Value: epn})
	}

	for _, pp := range filePathGroupResp.PriorityPaths {
		result.PriorityPaths.Elems = append(result.PriorityPaths.Elems, types.String{Value: pp})
	}

	signatures := make([]FilePathGroupSignature, 0)
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.String{Value: s.ID},
			Name:        types.String{Value: s.Name},
			Description: types.String{Value: s.Description},
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	yaraGroupRules := make([]YaraGroupRule, 0)
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.String{Value: ygr.ID},
			Name:        types.String{Value: ygr.Name},
			Description: types.String{Value: ygr.Description},
			Rules:       types.String{Value: ygr.Rules},
			Custom:      types.Bool{Value: ygr.Custom},
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceFilePathGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state FilePathGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePathGroupID := state.ID.Value

	// Retrieve values from plan
	var plan FilePathGroup
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var includePaths []string
	plan.IncludePaths.ElementsAs(ctx, &includePaths, false)

	var includePathExtensions []string
	plan.IncludePathExtensions.ElementsAs(ctx, &includePathExtensions, false)

	var excludePaths []string
	plan.ExcludePaths.ElementsAs(ctx, &excludePaths, false)

	var excludeProcessNames []string
	plan.ExcludeProcessNames.ElementsAs(ctx, &excludeProcessNames, false)

	var priorityPaths []string
	plan.PriorityPaths.ElementsAs(ctx, &priorityPaths, false)

	_signatures := make([]uptycs.FilePathGroupSignature, 0)
	for _, s := range plan.Signatures {
		_signatures = append(_signatures, uptycs.FilePathGroupSignature{
			ID:   s.ID.Value,
			Name: s.Name.Value,
		})
	}

	_yaraGroupRules := make([]uptycs.YaraGroupRule, 0)
	for _, yg := range plan.YaraGroupRules {
		_yaraGroupRules = append(_yaraGroupRules, uptycs.YaraGroupRule{
			ID:   yg.ID.Value,
			Name: yg.Name.Value,
		})
	}

	filePathGroupResp, err := r.p.client.UpdateFilePathGroup(uptycs.FilePathGroup{
		ID:                    filePathGroupID,
		Name:                  plan.Name.Value,
		Description:           plan.Description.Value,
		Grouping:              plan.Grouping.Value,
		IncludePaths:          includePaths,
		IncludePathExtensions: includePathExtensions,
		ExcludePaths:          excludePaths,
		Custom:                plan.Custom.Value,
		CheckSignature:        plan.CheckSignature.Value,
		FileAccesses:          plan.FileAccesses.Value,
		ExcludeProcessNames:   excludeProcessNames,
		PriorityPaths:         priorityPaths,
		Signatures:            _signatures,
		YaraGroupRules:        _yaraGroupRules,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create filePathGroup, unexpected error: "+err.Error(),
		)
		return
	}

	var result = FilePathGroup{
		ID:          types.String{Value: filePathGroupResp.ID},
		Name:        types.String{Value: filePathGroupResp.Name},
		Description: types.String{Value: filePathGroupResp.Description},
		Grouping:    types.String{Value: filePathGroupResp.Grouping},
		IncludePaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		IncludePathExtensions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		ExcludePaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom:         types.Bool{Value: filePathGroupResp.Custom},
		CheckSignature: types.Bool{Value: filePathGroupResp.CheckSignature},
		FileAccesses:   types.Bool{Value: filePathGroupResp.FileAccesses},
		ExcludeProcessNames: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		PriorityPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, ip := range filePathGroupResp.IncludePaths {
		result.IncludePaths.Elems = append(result.IncludePaths.Elems, types.String{Value: ip})
	}

	for _, ipe := range filePathGroupResp.IncludePathExtensions {
		result.IncludePathExtensions.Elems = append(result.IncludePathExtensions.Elems, types.String{Value: ipe})
	}

	for _, ep := range filePathGroupResp.ExcludePaths {
		result.ExcludePaths.Elems = append(result.ExcludePaths.Elems, types.String{Value: ep})
	}

	for _, epn := range filePathGroupResp.ExcludeProcessNames {
		result.ExcludeProcessNames.Elems = append(result.ExcludeProcessNames.Elems, types.String{Value: epn})
	}

	for _, pp := range filePathGroupResp.PriorityPaths {
		result.PriorityPaths.Elems = append(result.PriorityPaths.Elems, types.String{Value: pp})
	}

	signatures := make([]FilePathGroupSignature, 0)
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.String{Value: s.ID},
			Name:        types.String{Value: s.Name},
			Description: types.String{Value: s.Description},
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	yaraGroupRules := make([]YaraGroupRule, 0)
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.String{Value: ygr.ID},
			Name:        types.String{Value: ygr.Name},
			Description: types.String{Value: ygr.Description},
			Rules:       types.String{Value: ygr.Rules},
			Custom:      types.Bool{Value: ygr.Custom},
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceFilePathGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FilePathGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePathGroupID := state.ID.Value

	_, err := r.p.client.DeleteFilePathGroup(uptycs.FilePathGroup{
		ID: filePathGroupID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete filePathGroup with ID  "+filePathGroupID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceFilePathGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
