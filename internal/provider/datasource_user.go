package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *VoidauthClient
}

type UserDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Username      types.String `tfsdk:"username"`
	Email         types.String `tfsdk:"email"`
	Name          types.String `tfsdk:"name"`
	EmailVerified types.Bool   `tfsdk:"email_verified"`
	Approved      types.Bool   `tfsdk:"approved"`
	MfaRequired   types.Bool   `tfsdk:"mfa_required"`
	// ToDo: expiresAt
	HasPassword  types.Bool                 `tfsdk:"has_password"`
	HasEmail     types.Bool                 `tfsdk:"has_email"`
	IsAdmin      types.Bool                 `tfsdk:"is_admin"`
	Groups       []UserDataSourceGroupModel `tfsdk:"groups"`
	CustomClaims []CustomClaimModel         `tfsdk:"custom_claims"`
	HasTotp      types.Bool                 `tfsdk:"has_totp"`
	HasPasskeys  types.Bool                 `tfsdk:"has_passkeys"`
	HasMfaGroups types.Bool                 `tfsdk:"has_mfa_groups"`
}

type UserDataSourceGroupModel struct {
	Id           types.String       `tfsdk:"id"`
	Name         types.String       `tfsdk:"name"`
	CustomClaims []CustomClaimModel `tfsdk:"custom_claims"`
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"username":       schema.StringAttribute{Computed: true},
			"email":          schema.StringAttribute{Computed: true},
			"name":           schema.StringAttribute{Computed: true},
			"email_verified": schema.BoolAttribute{Computed: true},
			"approved":       schema.BoolAttribute{Computed: true},
			"mfa_required":   schema.BoolAttribute{Computed: true},
			"has_password":   schema.BoolAttribute{Computed: true},
			"has_email":      schema.BoolAttribute{Computed: true},
			"is_admin":       schema.BoolAttribute{Computed: true},
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
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
			"has_totp":       schema.BoolAttribute{Computed: true},
			"has_passkeys":   schema.BoolAttribute{Computed: true},
			"has_mfa_groups": schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = c
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	me, err := d.client.client.AdminService.GetUser(config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read admin user: %s", err))
		return
	}

	if me.Id != nil {
		config.Id = types.StringValue(*me.Id)
	}
	if me.Username != nil {
		config.Username = types.StringValue(*me.Username)
	}
	if me.Email != nil {
		config.Email = types.StringValue(*me.Email)
	}
	if me.Name != nil {
		config.Name = types.StringValue(*me.Name)
	}
	if me.EmailVerified != nil {
		config.EmailVerified = types.BoolValue(*me.EmailVerified)
	}
	if me.Approved != nil {
		config.Approved = types.BoolValue(*me.Approved)
	}
	if me.MfaRequired != nil {
		config.MfaRequired = types.BoolValue(*me.MfaRequired)
	}
	if me.HasPassword != nil {
		config.HasPassword = types.BoolValue(*me.HasPassword)
	}
	if me.HasEmail != nil {
		config.HasEmail = types.BoolValue(*me.HasEmail)
	}
	if me.IsAdmin != nil {
		config.IsAdmin = types.BoolValue(*me.IsAdmin)
	}
	if me.Groups != nil {
		groups := make([]UserDataSourceGroupModel, 0, len(*me.Groups))
		for _, g := range *me.Groups {
			gm := UserDataSourceGroupModel{}
			if g.Id != nil {
				gm.Id = types.StringValue(*g.Id)
			}
			if g.Name != nil {
				gm.Name = types.StringValue(*g.Name)
			}
			gm.CustomClaims = claimsToModels(g.CustomClaims)
			groups = append(groups, gm)
		}
		config.Groups = groups
	}
	if me.CustomClaims != nil {
		config.CustomClaims = claimsToModels(me.CustomClaims)
	}
	if me.HasTotp != nil {
		config.HasTotp = types.BoolValue(*me.HasTotp)
	}
	if me.HasPasskeys != nil {
		config.HasPasskeys = types.BoolValue(*me.HasPasskeys)
	}
	if me.HasMfaGroup != nil {
		config.HasMfaGroups = types.BoolValue(*me.HasMfaGroup)
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
