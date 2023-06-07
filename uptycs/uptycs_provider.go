package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
	"os"
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

func (p *UptycsProvider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host":        schema.StringAttribute{Optional: true},
			"api_key":     schema.StringAttribute{Optional: true},
			"api_secret":  schema.StringAttribute{Optional: true, Sensitive: true},
			"customer_id": schema.StringAttribute{Optional: true, Sensitive: true},
		},
	}
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
	if config.CustomerID.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as customerID",
		)
		return
	}

	if config.CustomerID.IsNull() {
		customerID = os.Getenv("UPTYCS_CUSTOMER_ID")
	} else {
		customerID = config.CustomerID.ValueString()
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
	if config.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as apiKey",
		)
		return
	}

	if config.APIKey.IsNull() {
		apiKey = os.Getenv("UPTYCS_API_KEY")
	} else {
		apiKey = config.APIKey.ValueString()
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
	if config.APISecret.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as APISecret",
		)
		return
	}

	if config.APISecret.IsNull() {
		apiSecret = os.Getenv("UPTYCS_API_SECRET")
	} else {
		apiSecret = config.APISecret.ValueString()
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
	if config.Host.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return
	}

	if config.Host.IsNull() {
		host = os.Getenv("UPTYCS_HOST")
	} else {
		host = config.Host.ValueString()
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
		ExceptionResource,
		FilePathGroupResource,
		FlagProfileResource,
		LookupTableResource,
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
		ExceptionDataSource,
		FilePathGroupDataSource,
		FlagProfileDataSource,
		LookupTableDataSource,
		ObjectGroupDataSource,
		QuerypackDataSource,
		RegistryPathDataSource,
		RoleDataSource,
		TagDataSource,
		TagRuleDataSource,
		UserDataSource,
		YaraGroupRuleDataSource,
	}
}
