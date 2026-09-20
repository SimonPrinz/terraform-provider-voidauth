package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &GroupDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

type GroupDataSource struct {
	client *VoidauthClient
}

type GroupDataSourceModel struct {
	ID           types.String               `tfsdk:"id"`
	Name         types.String               `tfsdk:"name"`
	MfaRequired  types.Bool                 `tfsdk:"mfa_required"`
	AutoAssign   types.Bool                 `tfsdk:"auto_assign"`
	Users        []GroupDataSourceUserModel `tfsdk:"users"`
	CustomClaims []CustomClaimModel         `tfsdk:"custom_claims"`
}

type GroupDataSourceUserModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
}

func (r *GroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *GroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"mfa_required": schema.BoolAttribute{
				Computed: true,
			},
			"auto_assign": schema.BoolAttribute{
				Computed: true,
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true},
						"username": schema.StringAttribute{Computed: true},
					},
				},
			},
			"custom_claims": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"claim": schema.StringAttribute{Computed: true},
						"value": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (r *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config GroupDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.client.AdminService.GetGroup(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Group not found",
			fmt.Sprintf("Group not found: %s", err),
		)
		return
	}

	if res.Name != nil {
		config.Name = types.StringValue(*res.Name)
	}
	if res.MfaRequired != nil {
		config.MfaRequired = types.BoolValue(*res.MfaRequired)
	}
	if res.AutoAssign != nil {
		config.AutoAssign = types.BoolValue(*res.AutoAssign)
	}
	if res.Users != nil {
		users := make([]GroupDataSourceUserModel, 0, len(*res.Users))
		for _, g := range *res.Users {
			um := GroupDataSourceUserModel{}
			if g.Id != nil {
				um.ID = types.StringValue(*g.Id)
			}
			if g.Username != nil {
				um.Username = types.StringValue(*g.Username)
			}
			users = append(users, um)
		}
		config.Users = users
	}
	if res.CustomClaims != nil {
		claims := make([]CustomClaimModel, 0, len(*res.CustomClaims))
		for _, g := range *res.CustomClaims {
			um := CustomClaimModel{}
			if g.Claim != nil {
				um.Claim = types.StringValue(*g.Claim)
			}
			if g.Value != nil {
				um.Value = types.StringValue(*g.Value)
			}
			claims = append(claims, um)
		}
		config.CustomClaims = claims
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
