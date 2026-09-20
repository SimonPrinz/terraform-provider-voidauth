package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &ConfigDataSource{}
)

func NewConfigDataSource() datasource.DataSource {
	return &ConfigDataSource{}
}

type ConfigDataSource struct {
	client *VoidauthClient
}

type ConfigDataSourceModel struct {
	Domain            types.String `tfsdk:"domain"`
	AppName           types.String `tfsdk:"app_name"`
	ZxcvbnMin         types.Int32  `tfsdk:"zxcvbn_min"`
	EmailActive       types.Bool   `tfsdk:"email_active"`
	EmailVerification types.Bool   `tfsdk:"email_verification"`
	Registration      types.Bool   `tfsdk:"registration"`
	ContactEmail      types.String `tfsdk:"contact_email"`
	DefaultRedirect   types.String `tfsdk:"default_redirect"`
	MfaRequired       types.Bool   `tfsdk:"mfa_required"`
}

func (d *ConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *ConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Computed: true,
			},
			"app_name": schema.StringAttribute{
				Computed: true,
			},
			"zxcvbn_min": schema.Int32Attribute{
				Computed: true,
			},
			"email_active": schema.BoolAttribute{
				Computed: true,
			},
			"email_verification": schema.BoolAttribute{
				Computed: true,
			},
			"registration": schema.BoolAttribute{
				Computed: true,
			},
			"contact_email": schema.StringAttribute{
				Computed: true,
			},
			"default_redirect": schema.StringAttribute{
				Computed: true,
			},
			"mfa_required": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *ConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicConfig, err := d.client.client.PublicService.Config()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read public config: %s", err))
		return
	}

	if publicConfig.Domain != nil {
		config.Domain = types.StringValue(*publicConfig.Domain)
	}
	if publicConfig.AppName != nil {
		config.AppName = types.StringValue(*publicConfig.AppName)
	}
	if publicConfig.ZxcvbnMin != nil {
		config.ZxcvbnMin = types.Int32Value(*publicConfig.ZxcvbnMin)
	}
	if publicConfig.EmailActive != nil {
		config.EmailActive = types.BoolValue(*publicConfig.EmailActive)
	}
	if publicConfig.EmailVerification != nil {
		config.EmailVerification = types.BoolValue(*publicConfig.EmailVerification)
	}
	if publicConfig.Registration != nil {
		config.Registration = types.BoolValue(*publicConfig.Registration)
	}
	if publicConfig.ContactEmail != nil {
		config.ContactEmail = types.StringValue(*publicConfig.ContactEmail)
	}
	if publicConfig.DefaultRedirect != nil {
		config.DefaultRedirect = types.StringValue(*publicConfig.DefaultRedirect)
	}
	if publicConfig.MfaRequired != nil {
		config.MfaRequired = types.BoolValue(*publicConfig.MfaRequired)
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
