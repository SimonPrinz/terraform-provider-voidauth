package provider

import (
	"context"
	"fmt"

	"github.com/SimonPrinz/terraform-provider-voidauth/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &GroupResource{}
	_ resource.ResourceWithConfigure = &GroupResource{}
)

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

type GroupResource struct {
	client *VoidauthClient
}

type GroupResourceModel struct {
	ID           types.String             `tfsdk:"id"`
	Name         types.String             `tfsdk:"name"`
	MfaRequired  types.Bool               `tfsdk:"mfa_required"`
	AutoAssign   types.Bool               `tfsdk:"auto_assign"`
	Users        []GroupResourceUserModel `tfsdk:"users"`
	CustomClaims []CustomClaimModel       `tfsdk:"custom_claims"`
}

type GroupResourceUserModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
}

func modelsToUserRefs(users []GroupResourceUserModel) []client.AdminUserRef {
	out := make([]client.AdminUserRef, 0, len(users))
	for _, u := range users {
		if u.ID.IsNull() || u.ID.IsUnknown() {
			continue
		}
		out = append(out, client.AdminUserRef{
			Id:       u.ID.ValueString(),
			Username: u.Username.ValueString(),
		})
	}
	return out
}

func (r *GroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*VoidauthClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *voidauthClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"mfa_required": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"auto_assign": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"users": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Required: true},
						"username": schema.StringAttribute{Required: true},
					},
				},
			},
			"custom_claims": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"claim": schema.StringAttribute{Required: true},
						"value": schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.client.AdminService.CreateGroup(&client.AdminGroupCreateRequest{
		Name:         plan.Name.ValueString(),
		MfaRequired:  plan.MfaRequired.ValueBool(),
		AutoAssign:   plan.AutoAssign.ValueBool(),
		Users:        modelsToUserRefs(plan.Users),
		CustomClaims: modelsToClaims(plan.CustomClaims),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating admin group",
			fmt.Sprintf("Error creating admin group: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(*res.Id)
	if res.Name != nil {
		plan.Name = types.StringValue(*res.Name)
	}
	if res.MfaRequired != nil {
		plan.MfaRequired = types.BoolValue(*res.MfaRequired)
	}
	if res.AutoAssign != nil {
		plan.AutoAssign = types.BoolValue(*res.AutoAssign)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.client.AdminService.GetGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Group not found",
			fmt.Sprintf("Group not found: %s", err),
		)
		return
	}

	if res.Name != nil {
		state.Name = types.StringValue(*res.Name)
	}
	if res.MfaRequired != nil {
		state.MfaRequired = types.BoolValue(*res.MfaRequired)
	}
	if res.AutoAssign != nil {
		state.AutoAssign = types.BoolValue(*res.AutoAssign)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state GroupResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.client.AdminService.CreateGroup(&client.AdminGroupCreateRequest{
		Id:           state.ID.ValueString(),
		Name:         plan.Name.ValueString(),
		MfaRequired:  plan.MfaRequired.ValueBool(),
		AutoAssign:   plan.AutoAssign.ValueBool(),
		Users:        modelsToUserRefs(plan.Users),
		CustomClaims: modelsToClaims(plan.CustomClaims),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating admin group",
			fmt.Sprintf("Error updating admin group: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(*res.Id)
	if res.Name != nil {
		plan.Name = types.StringValue(*res.Name)
	}
	if res.MfaRequired != nil {
		plan.MfaRequired = types.BoolValue(*res.MfaRequired)
	}
	if res.AutoAssign != nil {
		plan.AutoAssign = types.BoolValue(*res.AutoAssign)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.client.AdminService.DeleteGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete group",
			fmt.Sprintf("Could not delete group: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(diags...)
}
