package uptycs

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/uptycs-client-go/uptycs"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *uptycs.Client
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

// Provider schema struct
type providerData struct {
	Host       types.String `tfsdk:"host"`
	CustomerID types.String `tfsdk:"customer_id"`
	ApiKey     types.String `tfsdk:"api_key"`
	ApiSecret  types.String `tfsdk:"api_secret"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var customerId string
	if config.CustomerID.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as customerId",
		)
		return
	}

	if config.CustomerID.Null {
		customerId = os.Getenv("UPTYCS_CUSTOMER_ID")
	} else {
		customerId = config.CustomerID.Value
	}

	if customerId == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find customer id",
			"CustomerID cannot be an empty string",
		)
		return
	}

	var apiKey string
	if config.ApiKey.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as apiKey",
		)
		return
	}

	if config.ApiKey.Null {
		apiKey = os.Getenv("UPTYCS_API_KEY")
	} else {
		apiKey = config.ApiKey.Value
	}

	if apiKey == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find api key",
			"ApiKey cannot be an empty string",
		)
		return
	}

	// User must provide an api secret to the provider
	var apiSecret string
	if config.ApiSecret.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as ApiSecret",
		)
		return
	}

	if config.ApiSecret.Null {
		apiSecret = os.Getenv("UPTYCS_API_SECRET")
	} else {
		apiSecret = config.ApiSecret.Value
	}

	if apiSecret == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find ApiSecret",
			"ApiSecret cannot be an empty string",
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

	// Create a new uptycs client and set it to the provider client
	c, err := uptycs.NewClient(uptycs.UptycsConfig{
		Host:       host,
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
		CustomerID: customerId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create uptycs client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"uptycs_alert_rule": resourceAlertRuleType{},
		"uptycs_event_rule": resourceEventRuleType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
