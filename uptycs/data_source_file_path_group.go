package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func FilePathGroupDataSource() datasource.DataSource {
	return &filePathGroupDataSource{}
}

type filePathGroupDataSource struct {
	client *uptycs.Client
}

func (d *filePathGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_path_group"
}

func (d *filePathGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *filePathGroupDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"grouping":    schema.StringAttribute{Optional: true},
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
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
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
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
						"name":        schema.StringAttribute{Optional: true},
						"description": schema.StringAttribute{Optional: true},
						"rules":       schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (d *filePathGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var filePathGroupID string
	var filePathGroupName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &filePathGroupID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &filePathGroupName)

	var filePathGroupToLookup uptycs.FilePathGroup

	if len(filePathGroupID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		filePathGroupToLookup = uptycs.FilePathGroup{
			Name: filePathGroupName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		filePathGroupToLookup = uptycs.FilePathGroup{
			ID: filePathGroupID,
		}
	}

	filePathGroupResp, err := d.client.GetFilePathGroup(filePathGroupToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
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

	var signatures []FilePathGroupSignature
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.StringValue(s.ID),
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	var yaraGroupRules []YaraGroupRule
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
