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
	_ datasource.DataSource              = &filePathGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &filePathGroupDataSource{}
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

func (d *filePathGroupDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Optional: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							Computed: true,
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
				Optional: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							Computed: true,
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
						"rules": {
							Type:     types.StringType,
							Optional: true,
						},
					},
				),
			},
		},
	}, nil
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

	var signatures []FilePathGroupSignature
	for _, s := range filePathGroupResp.Signatures {
		signatures = append(signatures, FilePathGroupSignature{
			ID:          types.String{Value: s.ID},
			Name:        types.String{Value: s.Name},
			Description: types.String{Value: s.Description},
			//Paths:       types.List{}, //TODO we dont have any signatures
		})
	}
	result.Signatures = signatures

	var yaraGroupRules []YaraGroupRule
	for _, ygr := range filePathGroupResp.YaraGroupRules {
		yaraGroupRules = append(yaraGroupRules, YaraGroupRule{
			ID:          types.String{Value: ygr.ID},
			Name:        types.String{Value: ygr.Name},
			Description: types.String{Value: ygr.Description},
			Rules:       types.String{Value: ygr.Rules},
		})
	}
	result.YaraGroupRules = yaraGroupRules

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
