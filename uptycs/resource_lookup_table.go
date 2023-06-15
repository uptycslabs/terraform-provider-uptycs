package uptycs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
	"golang.org/x/exp/slices"
)

func LookupTableResource() resource.Resource {
	return &lookupTableResource{}
}

type lookupTableResource struct {
	client *uptycs.Client
}

func (r *lookupTableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lookup_table"
}

func (r *lookupTableResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *lookupTableResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"id_field": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString(""),
				},
			},
			"data_rows": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *lookupTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var lookupTableID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &lookupTableID)...)
	lookupTableResp, err := r.client.GetLookupTable(uptycs.LookupTable{
		ID: lookupTableID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get lookupTable with ID  "+lookupTableID+": "+err.Error(),
		)
		return
	}
	var result = LookupTable{
		ID:          types.StringValue(lookupTableResp.ID),
		Name:        types.StringValue(lookupTableResp.Name),
		Description: types.StringValue(lookupTableResp.Description),
		IDField:     types.StringValue(lookupTableResp.IDField),
		DataRows:    makeListStringAttributeFn(lookupTableResp.DataRows, func(v uptycs.LookupTableDataRow) (string, bool) { return string(v.Data), true }),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *lookupTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan LookupTable
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lookupTableResp, err := r.client.CreateLookupTable(uptycs.LookupTable{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IDField:     plan.IDField.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create lookupTable, unexpected error: "+err.Error(),
		)
		return
	}

	var dataRows []string
	plan.DataRows.ElementsAs(ctx, &dataRows, false)

	// Gather all the data rows from the plan
	for _, _dr := range dataRows {
		_, _ = r.client.CreateLookupTableDataRow(
			lookupTableResp,
			uptycs.LookupTableDataRow{
				Data: uptycs.CustomJSONString(fmt.Sprintf("[%s]", _dr)),
			},
		)
	}

	updatedLookupTableResp, _ := r.client.GetLookupTable(uptycs.LookupTable{
		ID: lookupTableResp.ID,
	})
	_dataRows := make([]string, len(updatedLookupTableResp.DataRows))
	for ind := range updatedLookupTableResp.DataRows {
		_dataRows[ind] = string(updatedLookupTableResp.DataRows[ind].Data)
	}

	var result = LookupTable{
		ID:          types.StringValue(lookupTableResp.ID),
		Name:        types.StringValue(lookupTableResp.Name),
		Description: types.StringValue(lookupTableResp.Description),
		IDField:     types.StringValue(lookupTableResp.IDField),
		DataRows:    makeListStringAttribute(_dataRows),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *lookupTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state LookupTable
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _actualState LookupTable
	resp.State.Get(ctx, &_actualState)

	lookupTableID := state.ID.ValueString()

	// Retrieve values from plan
	var plan LookupTable
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lookupTableResp, err := r.client.UpdateLookupTable(uptycs.LookupTable{
		ID:          lookupTableID,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IDField:     plan.IDField.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create lookupTable, unexpected error: "+err.Error(),
		)
		return
	}

	var planDataRows []string
	plan.DataRows.ElementsAs(ctx, &planDataRows, false)

	var stateDataRows []string
	_actualState.DataRows.ElementsAs(ctx, &stateDataRows, false)

	var keysInPlan []string
	var keysInState []string

	for _, _dr := range planDataRows {
		_, _tempStr, _ := getKeyValueFromRawJSON(_dr, plan.IDField.ValueString())
		keysInPlan = append(keysInPlan, _tempStr)
	}
	for _, _dr := range stateDataRows {
		_, _tempStr, _ := getKeyValueFromRawJSON(_dr, plan.IDField.ValueString())
		keysInState = append(keysInState, _tempStr)
	}

	_completed := make([]string, 0)
	for _, _dr := range append(planDataRows, stateDataRows...) {
		_, _tempStr, _ := getKeyValueFromRawJSON(_dr, plan.IDField.ValueString())
		if slices.Contains(difference(keysInState, keysInPlan), _tempStr) {
			//delete
			_, err = r.client.DeleteLookupTableDataRow(
				lookupTableResp,
				uptycs.LookupTableDataRow{
					IDFieldValue: _tempStr,
				},
			)
			if err != nil {
				panic(fmt.Sprintf("Error deleting lookup table row '%s': %s", _tempStr, err))
			}
		} else if slices.Contains(difference(keysInPlan, keysInState), _tempStr) {
			//add
			_, err = r.client.CreateLookupTableDataRow(
				lookupTableResp,
				uptycs.LookupTableDataRow{
					IDFieldValue: _tempStr,
					Data:         uptycs.CustomJSONString(fmt.Sprintf("[%s]", _dr)),
				},
			)
			if err != nil {
				panic(fmt.Sprintf("Error adding lookup table row '%s': %s", _tempStr, err))
			}

		} else if slices.Contains(interSection(keysInPlan, keysInState), _tempStr) {
			//update
			if !slices.Contains(_completed, _tempStr) {
				_updatedRow, err := r.client.UpdateLookupTableDataRow(
					lookupTableResp,
					uptycs.LookupTableDataRow{
						IDFieldValue: _tempStr,
						Data:         uptycs.CustomJSONString(fmt.Sprintf("[%s]", _dr)),
					},
				)
				if err != nil {
					panic(fmt.Sprintf("Error updating lookup table row '%s' with data '%+v': %s", _tempStr, _updatedRow.Data, err))
				}
				_completed = append(_completed, _tempStr)
			}
		}

	}

	updatedLookupTableResp, _ := r.client.GetLookupTable(uptycs.LookupTable{
		ID: lookupTableResp.ID,
	})
	_dataRows := make([]string, len(updatedLookupTableResp.DataRows))
	for ind := range updatedLookupTableResp.DataRows {
		_dataRows[ind] = string(updatedLookupTableResp.DataRows[ind].Data)
	}

	var result = LookupTable{
		ID:          types.StringValue(lookupTableResp.ID),
		Name:        types.StringValue(lookupTableResp.Name),
		Description: types.StringValue(lookupTableResp.Description),
		IDField:     types.StringValue(lookupTableResp.IDField),
		DataRows:    makeListStringAttribute(_dataRows),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *lookupTableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LookupTable
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lookupTableID := state.ID.ValueString()

	_, err := r.client.DeleteLookupTable(uptycs.LookupTable{
		ID: lookupTableID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete lookupTable with ID  "+lookupTableID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *lookupTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
