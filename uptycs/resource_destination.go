package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func DestinationResource() resource.Resource {
	return &destinationResource{}
}

type destinationResource struct {
	client *uptycs.Client
}

func (r *destinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (r *destinationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *destinationResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true},
			"name":    schema.StringAttribute{Optional: true},
			"type":    schema.StringAttribute{Optional: true},
			"address": schema.StringAttribute{Optional: true},
			"enabled": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(true),
				},
			},
			"config": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"sender": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"method": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"username": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"password": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data_key": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"token": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"slack_attachment": schema.BoolAttribute{
						Required: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"headers": schema.StringAttribute{Required: true},
				},
			},
			"template": schema.StringAttribute{Optional: true},
		},
	}
}

func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var destinationID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &destinationID)...)
	destinationResp, err := r.client.GetDestination(uptycs.Destination{
		ID: destinationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}

	headersJSON, err := json.MarshalIndent(destinationResp.Config.Headers, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
		Config: DestinationConfig{
			Sender:   types.StringValue(destinationResp.Config.Sender),
			Method:   types.StringValue(destinationResp.Config.Method),
			Username: types.StringValue(destinationResp.Config.Username),
			//Password:        types.StringValue("**********"), //They will never give this back in a response
			DataKey:         types.StringValue(destinationResp.Config.DataKey),
			Token:           types.StringValue(destinationResp.Config.Token),
			SlackAttachment: types.BoolValue(destinationResp.Config.SlackAttachment),
			Headers:         types.StringValue(string(headersJSON) + "\n"),
		},
		Template: types.StringValue(destinationResp.Template.Template),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Destination
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationResp, err := r.client.CreateDestination(uptycs.Destination{
		Name:    plan.Name.ValueString(),
		Type:    plan.Type.ValueString(),
		Address: plan.Address.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
		Config: uptycs.DestinationConfig{
			Sender:          plan.Config.Sender.ValueString(),
			Method:          plan.Config.Method.ValueString(),
			Username:        plan.Config.Username.ValueString(),
			Password:        plan.Config.Password.ValueString(),
			DataKey:         plan.Config.DataKey.ValueString(),
			Token:           plan.Config.Token.ValueString(),
			SlackAttachment: plan.Config.SlackAttachment.ValueBool(),
			Headers:         uptycs.CustomJSONString(plan.Config.Headers.ValueString() + "\n"),
		},
		Template: struct {
			Template string `json:"template,omitempty"`
		}{
			Template: plan.Template.ValueString(),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	headersJSON, err := json.MarshalIndent(destinationResp.Config.Headers, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
		Config: DestinationConfig{
			Sender:          types.StringValue(destinationResp.Config.Sender),
			Method:          types.StringValue(destinationResp.Config.Method),
			Username:        types.StringValue(destinationResp.Config.Username),
			Password:        plan.Config.Password, //They will never give this back in a response
			DataKey:         types.StringValue(destinationResp.Config.DataKey),
			Token:           types.StringValue(destinationResp.Config.Token),
			SlackAttachment: types.BoolValue(destinationResp.Config.SlackAttachment),
			Headers:         types.StringValue(string(headersJSON) + "\n"),
		},
		Template: types.StringValue(destinationResp.Template.Template),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.ValueString()

	// Retrieve values from plan
	var plan Destination
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destination := uptycs.Destination{
		ID:      destinationID,
		Name:    plan.Name.ValueString(),
		Type:    plan.Type.ValueString(),
		Address: plan.Address.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
		Config: uptycs.DestinationConfig{
			Sender:          plan.Config.Sender.ValueString(),
			Method:          plan.Config.Method.ValueString(),
			Username:        plan.Config.Username.ValueString(),
			DataKey:         plan.Config.DataKey.ValueString(),
			Token:           plan.Config.Token.ValueString(),
			SlackAttachment: plan.Config.SlackAttachment.ValueBool(),
			Headers:         uptycs.CustomJSONString(plan.Config.Headers.ValueString() + "\n"),
		},
		Template: struct {
			Template string `json:"template,omitempty"`
		}{
			Template: plan.Template.ValueString(),
		},
	}
	if plan.Config.Password.ValueString() != "" {
		destination.Config.Password = plan.Config.Password.ValueString()
	}

	destinationResp, err := r.client.UpdateDestination(destination)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	headersJSON, err := json.MarshalIndent(destinationResp.Config.Headers, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
		Config: DestinationConfig{
			Sender:          types.StringValue(destinationResp.Config.Sender),
			Method:          types.StringValue(destinationResp.Config.Method),
			Username:        types.StringValue(destinationResp.Config.Username),
			Password:        types.StringValue("**********"), //They will never give this back in a response
			DataKey:         types.StringValue(destinationResp.Config.DataKey),
			Token:           types.StringValue(destinationResp.Config.Token),
			SlackAttachment: types.BoolValue(destinationResp.Config.SlackAttachment),
			Headers:         types.StringValue(string(headersJSON) + "\n"),
		},
		Template: types.StringValue(destinationResp.Template.Template),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.ValueString()

	_, err := r.client.DeleteDestination(uptycs.Destination{
		ID: destinationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *destinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
