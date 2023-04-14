package windowsca

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &windowscaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &windowscaProvider{}
}

type windowscaProviderModel struct {
	WinrmHostname types.String `tfsdk:"winrm_hostname"`
	WinrmUsername types.String `tfsdk:"winrm_username"`
	WinrmPassword types.String `tfsdk:"winrm_password"`
}

// windowscaProvider is the provider implementation.
type windowscaProvider struct{}

// Metadata returns the provider type name.
func (p *windowscaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "windowsca"
}

// Schema defines the provider-level schema for configuration data.
func (p *windowscaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"winrm_hostname": schema.StringAttribute{
				Optional: true,
			},
			"winrm_username": schema.StringAttribute{
				Optional: true,
			},
			"winrm_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a Windows WinRM client for data sources and resources.
func (p *windowscaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config windowscaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.WinrmHostname.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("winrm_hostname"),
			"Unknown Windows WinRM Host",
			"The provider cannot create the Windows WinRM client as there is an unknown configuration value for the Windows WinRM host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WINDOWSCA_HOSTNAME environment variable.",
		)
	}

	if config.WinrmUsername.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("winrm_username"),
			"Unknown Windows WinRM Username",
			"The provider cannot create the Windows WinRM client as there is an unknown configuration value for the Windows WinRM username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WINDOWSCA_USERNAME environment variable.",
		)
	}

	if config.WinrmPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("winrm_password"),
			"Unknown Windows WinRM Password",
			"The provider cannot create the Windows WinRM client as there is an unknown configuration value for the Windows WinRM password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WINDOWSCA_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("WINDOWSCA_HOSTNAME")
	username := os.Getenv("WINDOWSCA_USERNAME")
	password := os.Getenv("WINDOWSCA_PASSWORD")

	if !config.WinrmHostname.IsNull() {
		host = config.WinrmHostname.ValueString()
	}

	if !config.WinrmUsername.IsNull() {
		username = config.WinrmUsername.ValueString()
	}

	if !config.WinrmPassword.IsNull() {
		password = config.WinrmPassword.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Windows WinRM Host",
			"The provider cannot create the Windows WinRM client as there is a missing or empty value for the Windows WinRM host. "+
				"Set the host value in the configuration or use the WINDOWSCA_HOSTNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Windows WinRM Username",
			"The provider cannot create the Windows WinRM client as there is a missing or empty value for the Windows WinRM username. "+
				"Set the username value in the configuration or use the WINDOWSCA_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Windows WinRM Password",
			"The provider cannot create the Windows WinRM client as there is a missing or empty value for the Windows WinRM password. "+
				"Set the password value in the configuration or use the WINDOWSCA_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new HashiCups client using the configuration values
	client, err := windowsca.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Windows WinRM Client",
			"An unexpected error occurred when creating the Windows WinRM client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Windows WinRM Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *windowscaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *windowscaProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
