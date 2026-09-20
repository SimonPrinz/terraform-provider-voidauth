package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &voidauthProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &voidauthProvider{
			version: version,
		}
	}
}

type voidauthProvider struct {
	version string
}

type voidauthProviderModel struct {
	Url      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *voidauthProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "voidauth"
	resp.Version = p.version
}

func (p *voidauthProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *voidauthProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config voidauthProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown URL",
			"URL must be set",
		)
		return
	}

	url := config.Url.ValueString()
	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing URL",
			"URL must be set",
		)
		return
	}
	username := ""
	password := ""
	if !config.Username.IsUnknown() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsUnknown() {
		password = config.Password.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := NewVoidauthClient(url, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create API client",
			fmt.Sprintf("Could not create voidauth API client: %s", err.Error()),
		)
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *voidauthProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProxyAuthResource,
		NewGroupResource,
	}
}

func (p *voidauthProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConfigDataSource,
		NewMeDataSource,
		NewUserDataSource,
		NewGroupDataSource,
	}
}
