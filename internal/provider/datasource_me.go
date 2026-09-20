package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &MeDataSource{}
	_ datasource.DataSourceWithConfigure = &MeDataSource{}
)

func NewMeDataSource() datasource.DataSource {
	return &MeDataSource{}
}

type MeDataSource struct {
	client *VoidauthClient
}

type MeDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	IsAdmin       types.Bool   `tfsdk:"is_admin"`
	EmailVerified types.Bool   `tfsdk:"email_verified"`
	HasTotp       types.Bool   `tfsdk:"has_totp"`
	HasPasskeys   types.Bool   `tfsdk:"has_passkeys"`
	//ExpiresAt                 interface{} `tfsdk:"expires_at"`
	Approved types.Bool `tfsdk:"approved"`
	//Amr                       []string    `tfsdk:"amr"`
	CanLogin types.Bool `tfsdk:"can_login"`
	HasEmail types.Bool `tfsdk:"has_email"`
}

func (d *MeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_me"
}

func (d *MeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Computed: true},
			"is_admin":       schema.BoolAttribute{Computed: true},
			"email_verified": schema.BoolAttribute{Computed: true},
			"has_totp":       schema.BoolAttribute{Computed: true},
			"has_passkeys":   schema.BoolAttribute{Computed: true},
			"approved":       schema.BoolAttribute{Computed: true},
			"can_login":      schema.BoolAttribute{Computed: true},
			"has_email":      schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *MeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MeDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	me, err := d.client.client.UserService.Me()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read user me: %s", err))
		return
	}

	if me.Id != nil {
		config.Id = types.StringValue(*me.Id)
	}
	if me.IsAdmin != nil {
		config.IsAdmin = types.BoolValue(*me.IsAdmin)
	}
	if me.EmailVerified != nil {
		config.EmailVerified = types.BoolValue(*me.EmailVerified)
	}
	if me.HasTotp != nil {
		config.HasTotp = types.BoolValue(*me.HasTotp)
	}
	if me.HasPasskeys != nil {
		config.HasPasskeys = types.BoolValue(*me.HasPasskeys)
	}
	if me.Approved != nil {
		config.Approved = types.BoolValue(*me.Approved)
	}
	if me.CanLogin != nil {
		config.CanLogin = types.BoolValue(*me.CanLogin)
	}
	if me.HasEmail != nil {
		config.HasEmail = types.BoolValue(*me.HasEmail)
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
