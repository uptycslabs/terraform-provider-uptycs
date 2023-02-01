package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func FilePathGroupResource() resource.Resource {
	return &filePathGroupResource{}
}

type filePathGroupResource struct {
	client *uptycs.Client
}

func (r *filePathGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_path_group"
}

func (r *filePathGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *filePathGroupResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true,
				Optional: true,
			},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"grouping": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString(""),
				},
			},
			"include_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"include_path_extensions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"exclude_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"check_signature": schema.BoolAttribute{Optional: true},
			"file_accesses":   schema.BoolAttribute{Optional: true},
			"exclude_process_names": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"priority_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"signatures": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"name":        schema.StringAttribute{Optional: true},
						"description": schema.StringAttribute{Optional: true},
						"paths": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"yara_group_rules": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"name": schema.StringAttribute{Computed: true,
							Optional: true,
						},
						"description": schema.StringAttribute{Computed: true,
							Optional: true,
						},
						"rules": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *filePathGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var filePathGroupID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &filePathGroupID)...)
	filePathGroupResp, err := r.client.GetFilePathGroup(uptycs.FilePathGroup{
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
		ID:                    types.StringValue(filePathGroupResp.ID),
		Name:                  types.StringValue(filePathGroupResp.Name),
		Description:           types.StringValue(filePathGroupResp.Description),
		Grouping:              types.StringValue(filePathGroupResp.Grouping),
		IncludePaths:          makeListStringAttribute(filePathGroupResp.IncludePaths),
		IncludePathExtensions: makeListStringAttribute(filePathGroupResp.IncludePathExtensions),
		ExcludePaths:          makeListStringAttribute(filePathGroupResp.ExcludePaths),
		CheckSignature:        types.BoolValue(filePathGroupResp.CheckSignature),
		FileAccesses:          types.BoolValue(filePathGroupResp.FileAccesses),
		ExcludeProcessNames:   makeListStringAttribute(filePathGroupResp.ExcludeProcessNames),
		PriorityPaths:         makeListStringAttribute(filePathGroupResp.PriorityPaths),
	}

	signatures := make([]FilePathGroupSignature, 0)
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.StringValue(s.ID),
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	yaraGroupRules := make([]YaraGroupRule, 0)
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.StringValue(ygr.ID),
			Name:        types.StringValue(ygr.Name),
			Description: types.StringValue(ygr.Description),
			Rules:       types.StringValue(ygr.Rules),
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *filePathGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
			ID:   s.ID.ValueString(),
			Name: s.Name.ValueString(),
		})
	}

	_yaraGroupRules := make([]uptycs.YaraGroupRule, 0)
	for _, yg := range plan.YaraGroupRules {
		_yaraGroupRules = append(_yaraGroupRules, uptycs.YaraGroupRule{
			ID:   yg.ID.ValueString(),
			Name: yg.Name.ValueString(),
		})
	}

	filePathGroupResp, err := r.client.CreateFilePathGroup(uptycs.FilePathGroup{
		Name:                  plan.Name.ValueString(),
		Description:           plan.Description.ValueString(),
		Grouping:              plan.Grouping.ValueString(),
		IncludePaths:          includePaths,
		IncludePathExtensions: includePathExtensions,
		ExcludePaths:          excludePaths,
		CheckSignature:        plan.CheckSignature.ValueBool(),
		FileAccesses:          plan.FileAccesses.ValueBool(),
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
		ID:                    types.StringValue(filePathGroupResp.ID),
		Name:                  types.StringValue(filePathGroupResp.Name),
		Description:           types.StringValue(filePathGroupResp.Description),
		Grouping:              types.StringValue(filePathGroupResp.Grouping),
		IncludePaths:          makeListStringAttribute(filePathGroupResp.IncludePaths),
		IncludePathExtensions: makeListStringAttribute(filePathGroupResp.IncludePathExtensions),
		ExcludePaths:          makeListStringAttribute(filePathGroupResp.ExcludePaths),
		CheckSignature:        types.BoolValue(filePathGroupResp.CheckSignature),
		FileAccesses:          types.BoolValue(filePathGroupResp.FileAccesses),
		ExcludeProcessNames:   makeListStringAttribute(filePathGroupResp.ExcludeProcessNames),
		PriorityPaths:         makeListStringAttribute(filePathGroupResp.PriorityPaths),
	}

	signatures := make([]FilePathGroupSignature, 0)
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.StringValue(s.ID),
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	yaraGroupRules := make([]YaraGroupRule, 0)
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.StringValue(ygr.ID),
			Name:        types.StringValue(ygr.Name),
			Description: types.StringValue(ygr.Description),
			Rules:       types.StringValue(ygr.Rules),
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *filePathGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state FilePathGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePathGroupID := state.ID.ValueString()

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
			ID:   s.ID.ValueString(),
			Name: s.Name.ValueString(),
		})
	}

	_yaraGroupRules := make([]uptycs.YaraGroupRule, 0)
	for _, yg := range plan.YaraGroupRules {
		_yaraGroupRules = append(_yaraGroupRules, uptycs.YaraGroupRule{
			ID:   yg.ID.ValueString(),
			Name: yg.Name.ValueString(),
		})
	}

	filePathGroupResp, err := r.client.UpdateFilePathGroup(uptycs.FilePathGroup{
		ID:                    filePathGroupID,
		Name:                  plan.Name.ValueString(),
		Description:           plan.Description.ValueString(),
		Grouping:              plan.Grouping.ValueString(),
		IncludePaths:          includePaths,
		IncludePathExtensions: includePathExtensions,
		ExcludePaths:          excludePaths,
		CheckSignature:        plan.CheckSignature.ValueBool(),
		FileAccesses:          plan.FileAccesses.ValueBool(),
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
		ID:                    types.StringValue(filePathGroupResp.ID),
		Name:                  types.StringValue(filePathGroupResp.Name),
		Description:           types.StringValue(filePathGroupResp.Description),
		Grouping:              types.StringValue(filePathGroupResp.Grouping),
		IncludePaths:          makeListStringAttribute(filePathGroupResp.IncludePaths),
		IncludePathExtensions: makeListStringAttribute(filePathGroupResp.IncludePathExtensions),
		ExcludePaths:          makeListStringAttribute(filePathGroupResp.ExcludePaths),
		CheckSignature:        types.BoolValue(filePathGroupResp.CheckSignature),
		FileAccesses:          types.BoolValue(filePathGroupResp.FileAccesses),
		ExcludeProcessNames:   makeListStringAttribute(filePathGroupResp.ExcludeProcessNames),
		PriorityPaths:         makeListStringAttribute(filePathGroupResp.PriorityPaths),
	}

	signatures := make([]FilePathGroupSignature, 0)
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.StringValue(s.ID),
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	yaraGroupRules := make([]YaraGroupRule, 0)
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.StringValue(ygr.ID),
			Name:        types.StringValue(ygr.Name),
			Description: types.StringValue(ygr.Description),
			Rules:       types.StringValue(ygr.Rules),
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *filePathGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FilePathGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePathGroupID := state.ID.ValueString()

	_, err := r.client.DeleteFilePathGroup(uptycs.FilePathGroup{
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

func (r *filePathGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
