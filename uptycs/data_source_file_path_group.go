package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceFilePathGroupType struct {
	p Provider
}

func (r dataSourceFilePathGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"custom": {
				Type:     types.BoolType,
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
						"custom": {
							Type:     types.BoolType,
							Optional: true,
						},
					},
				),
			},
		},
	}, nil
}

func (r dataSourceFilePathGroupType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceFilePathGroupType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceFilePathGroupType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var filePathGroupID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &filePathGroupID)...)

	filePathGroupResp, err := r.p.client.GetFilePathGroup(uptycs.FilePathGroup{
		ID: filePathGroupID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get filePathGroup with ID  "+filePathGroupID+": "+err.Error(),
		)
		return
	}

	//ID:               types.String{Value: filePathGroupResp.ID},
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
