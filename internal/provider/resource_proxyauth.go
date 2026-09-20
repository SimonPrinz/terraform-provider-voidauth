package provider

import (
	"context"
	"fmt"

	"github.com/SimonPrinz/terraform-provider-voidauth/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &ProxyAuthResource{}
	_ resource.ResourceWithConfigure = &ProxyAuthResource{}
)

func NewProxyAuthResource() resource.Resource {
	return &ProxyAuthResource{}
}

type ProxyAuthResource struct {
	client *VoidauthClient
}

type ProxyAuthResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	Domain           types.String   `tfsdk:"domain"`
	MfaRequired      types.Bool     `tfsdk:"mfa_required"`
	MaxSessionLength types.Int64    `tfsdk:"max_session_length"`
	Groups           []types.String `tfsdk:"groups"`
}

func stringSlice(vals []types.String) []string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		if v.IsNull() || v.IsUnknown() {
			continue
		}
		out = append(out, v.ValueString())
	}
	return out
}

func (r *ProxyAuthResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProxyAuthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxyauth"
}

func (r *ProxyAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "Domain (supports wildcards) matched for ProxyAuth. The server stores it in normalized wildcard form (app.example.com becomes app.example.com/*); the configured spelling is kept in state.",
			},
			"mfa_required": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"max_session_length": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(5, 525600),
				},
			},
			"groups": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *ProxyAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProxyAuthResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.client.AdminService.CreateProxyAuth(&client.AdminProxyAuthCreateRequest{
		Domain:           plan.Domain.ValueString(),
		MaxSessionLength: int(plan.MaxSessionLength.ValueInt64()),
		MfaRequired:      plan.MfaRequired.ValueBool(),
		Groups:           stringSlice(plan.Groups),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating admin proxyauth",
			fmt.Sprintf("Error creating admin proxyauth: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(*res.Id)
	if res.MfaRequired != nil {
		plan.MfaRequired = types.BoolValue(*res.MfaRequired)
	}
	if res.MaxSessionLength != nil {
		plan.MaxSessionLength = types.Int64Value(int64(*res.MaxSessionLength))
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ProxyAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProxyAuthResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := r.client.client.AdminService.GetProxyAuth(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Proxyauth not found",
			fmt.Sprintf("Proxyauth not found: %s", err),
		)
		return
	}

	if c.MfaRequired != nil {
		state.MfaRequired = types.BoolValue(*c.MfaRequired)
	}
	if c.MaxSessionLength != nil {
		state.MaxSessionLength = types.Int64Value(int64(*c.MaxSessionLength))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ProxyAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProxyAuthResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ProxyAuthResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.client.AdminService.CreateProxyAuth(&client.AdminProxyAuthCreateRequest{
		Id:               state.ID.ValueString(),
		Domain:           plan.Domain.ValueString(),
		MaxSessionLength: int(plan.MaxSessionLength.ValueInt64()),
		MfaRequired:      plan.MfaRequired.ValueBool(),
		Groups:           stringSlice(plan.Groups),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating admin proxyauth",
			fmt.Sprintf("Error updating admin proxyauth: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(*res.Id)
	if res.MfaRequired != nil {
		plan.MfaRequired = types.BoolValue(*res.MfaRequired)
	}
	if res.MaxSessionLength != nil {
		plan.MaxSessionLength = types.Int64Value(int64(*res.MaxSessionLength))
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ProxyAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProxyAuthResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.client.AdminService.DeleteProxyAuth(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete proxyauth",
			fmt.Sprintf("Could not delete proxyauth: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(diags...)
}
