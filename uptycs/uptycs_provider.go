package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ provider.Provider = &UptycsProvider{}
)

func New() provider.Provider {
	return &UptycsProvider{}
}

type UptycsProvider struct{} //revive:disable-line:exported

type uptycsProviderData struct {
	Host       types.String `tfsdk:"host"`
	CustomerID types.String `tfsdk:"customer_id"`
	APIKey     types.String `tfsdk:"api_key"`
	APISecret  types.String `tfsdk:"api_secret"`
}

func (p *UptycsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uptycs"
}

func (p *UptycsProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"api_key": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"api_secret": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"customer_id": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

func (p *UptycsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Uptycs client")

	// Retrieve provider data from configuration
	var config uptycsProviderData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var customerID string
	if config.CustomerID.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as customerID",
		)
		return
	}

	if config.CustomerID.Null {
		customerID = os.Getenv("UPTYCS_CUSTOMER_ID")
	} else {
		customerID = config.CustomerID.Value
	}

	if customerID == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find customer id",
			"CustomerID cannot be an empty string",
		)
		return
	}

	var apiKey string
	if config.APIKey.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as apiKey",
		)
		return
	}

	if config.APIKey.Null {
		apiKey = os.Getenv("UPTYCS_API_KEY")
	} else {
		apiKey = config.APIKey.Value
	}

	if apiKey == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find api key",
			"APIKey cannot be an empty string",
		)
		return
	}

	// User must provide an api secret to the provider
	var apiSecret string
	if config.APISecret.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as APISecret",
		)
		return
	}

	if config.APISecret.Null {
		apiSecret = os.Getenv("UPTYCS_API_SECRET")
	} else {
		apiSecret = config.APISecret.Value
	}

	if apiSecret == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find APISecret",
			"APISecret cannot be an empty string",
		)
		return
	}

	// User must specify a host
	var host string
	if config.Host.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return
	}

	if config.Host.Null {
		host = os.Getenv("UPTYCS_HOST")
	} else {
		host = config.Host.Value
	}

	if host == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Host cannot be an empty string",
		)
		return
	}

	// Create a new uptycs client and set it to the provider.client
	client, err := uptycs.NewClient(uptycs.Config{
		Host:       host,
		APIKey:     apiKey,
		APISecret:  apiSecret,
		CustomerID: customerID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create uptycs client:\n\n"+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Uptycs client", map[string]any{"success": true})

}

func (p *UptycsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		AlertRuleResource,
		ComplianceProfileResource,
		CustomProfileResource,
		DestinationResource,
		EventExcludeProfileResource,
		EventRuleResource,
		FilePathGroupResource,
		FlagProfileResource,
		QuerypackResource,
		RegistryPathResource,
		RoleResource,
		TagResource,
		TagRuleResource,
		UserResource,
		YaraGroupRuleResource,
	}
}

func (p *UptycsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		AlertRuleDataSource,
		AlertRuleCategoryDataSource,
		AssetGroupRuleDataSource,
		AtcQueryDataSource,
		AuditConfigurationDataSource,
		ComplianceProfileDataSource,
		CustomProfileDataSource,
		DestinationDataSource,
		EventRuleDataSource,
		EventExcludeProfileDataSource,
		FilePathGroupDataSource,
		FlagProfileDataSource,
		QuerypackDataSource,
		RegistryPathDataSource,
		RoleDataSource,
		TagDataSource,
		TagRuleDataSource,
		UserDataSource,
		YaraGroupRuleDataSource,
	}
}
