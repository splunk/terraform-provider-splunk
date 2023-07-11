package splunk

import (
	"fmt"
	"net/url"
	"time"

	"github.com/splunk/terraform-provider-splunk/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/go-splunk-client/pkg/authenticators"
	externalclient "github.com/splunk/go-splunk-client/pkg/client"
)

const (
	useClientUnset    = ""
	useClientLegacy   = "legacy"
	useClientExternal = "external"
)

type SplunkProvider struct {
	Client           *client.Client
	ExternalClient   *externalclient.Client
	useClientDefault string
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: providerDataSources(),
		ResourcesMap:   providerResources(),
		ConfigureFunc:  providerConfigure,
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_URL", nil),
			Description: "Splunk instance URL",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_USERNAME", nil),
			Description: "Splunk instance admin username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_PASSWORD", nil),
			Description: "Splunk instance password",
		},
		"auth_token": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_AUTH_TOKEN", nil),
			Description: "Authentication tokens, also known as JSON Web Tokens (JWT), are a method for authenticating " +
				"Splunk platform users into the Splunk platform",
		},
		"insecure_skip_verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_INSECURE_SKIP_VERIFY", true),
			Description: "insecure skip verification flag",
		},
		"timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_TIMEOUT", 60),
			Description: "Timeout when making calls to Splunk server. Defaults to 60 seconds",
		},
		"use_client_default": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      useClientLegacy,
			ValidateFunc: validateClientValueFunc(false),
			Description: "Determines the default behavior for resources that implement use_client. Permitted values are legacy and external. " +
				"Currently defaults to legacy, but will default to external in a future version. " +
				"The legacy client is being replaced by a standalone Splunk client with improved error and drift handling. The legacy client will be deprecated in a future version.",
		},
	}
}

// Returns a map of splunk resources for configuration
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"splunk_admin_saml_groups":           adminSAMLGroups(),
		"splunk_apps_local":                  appsLocal(),
		"splunk_authentication_users":        authenticationUsers(),
		"splunk_authorization_roles":         authorizationRoles(),
		"splunk_global_http_event_collector": globalHttpEventCollector(),
		"splunk_generic_acl":                 genericAcl(),
		"splunk_inputs_http_event_collector": inputsHttpEventCollector(),
		"splunk_inputs_script":               inputsScript(),
		"splunk_inputs_monitor":              inputsMonitor(),
		"splunk_inputs_udp":                  inputsUDP(),
		"splunk_inputs_tcp_raw":              inputsTCPRaw(),
		"splunk_inputs_tcp_cooked":           inputsTCPCooked(),
		"splunk_inputs_tcp_splunk_tcp_token": inputsTCPSplunkTCPToken(),
		"splunk_inputs_tcp_ssl":              inputsTCPSSL(),
		"splunk_outputs_tcp_default":         outputsTCPDefault(),
		"splunk_outputs_tcp_server":          outputsTCPServer(),
		"splunk_outputs_tcp_group":           outputsTCPGroup(),
		"splunk_outputs_tcp_syslog":          outputsTCPSyslog(),
		"splunk_saved_searches":              savedSearches(),
		"splunk_sh_indexes_manager":          shIndexesManager(),
		"splunk_indexes":                     index(),
		"splunk_configs_conf":                configsConf(),
		"splunk_data_ui_views":               splunkDashboards(),
	}
}

// canonicalizeSplunkURL returns a URL string from originalURL that includes a default https scheme.
// go-splunk-client requires a full URL, but this provider has historically permitted the URL to be missing
// its scheme.
func canonicalizeSplunkURL(originalURL string) (string, error) {
	externalClientURL := originalURL

	parsedURL, err := url.Parse(externalClientURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		externalClientURL = fmt.Sprintf("https://%s", externalClientURL)
		if _, err := url.Parse(externalClientURL); err != nil {
			return "", fmt.Errorf("splunk: unable to determine valid splunkd URL from %q", originalURL)
		}
	}

	return externalClientURL, nil
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise splunk client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	provider := &SplunkProvider{}
	var splunkdClient *client.Client

	externalClientURL, err := canonicalizeSplunkURL(d.Get("url").(string))
	if err != nil {
		return nil, err
	}

	externalClient := externalclient.Client{
		URL: externalClientURL,
	}

	httpClient, err := client.NewSplunkdHTTPClient(
		time.Duration(d.Get("timeout").(int))*time.Second,
		d.Get("insecure_skip_verify").(bool))
	if err != nil {
		return nil, err
	}

	externalClient.TLSInsecureSkipVerify = d.Get("insecure_skip_verify").(bool)
	externalClient.Timeout = time.Duration(d.Get("timeout").(int)) * time.Second

	if token, ok := d.GetOk("auth_token"); ok {
		splunkdClient, err = client.NewSplunkdClientWithAuthToken(token.(string),
			[2]string{d.Get("username").(string), d.Get("password").(string)},
			d.Get("url").(string),
			httpClient)
		if err != nil {
			return splunkdClient, err
		}

		externalClient.Authenticator = authenticators.Token{
			Token: token.(string),
		}
	} else {
		splunkdClient, err = client.NewSplunkdClient("",
			[2]string{d.Get("username").(string), d.Get("password").(string)},
			d.Get("url").(string),
			httpClient)
		if err != nil {
			return splunkdClient, err
		}
		// Login is required to get session key
		err = splunkdClient.Login()
		if err != nil {
			return splunkdClient, err
		}

		externalClient.Authenticator = &authenticators.Password{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		}
	}

	provider.Client = splunkdClient
	provider.ExternalClient = &externalClient
	provider.useClientDefault = d.Get("use_client_default").(string)

	return provider, nil
}

// validateClientValueFunc returns a schema.SchemaValidateFunc that validates the value
// of use_client (or use_client_default).
func validateClientValueFunc(allowEmpty bool) schema.SchemaValidateFunc {
	return func(v interface{}, name string) ([]string, []error) {
		clientV := v.(string)
		switch clientV {
		default:
			return nil, []error{fmt.Errorf("%s invalid value %q", name, clientV)}
		case useClientUnset:
			if !allowEmpty {
				return nil, []error{fmt.Errorf("%s invalid empty value", name)}
			}
		case useClientLegacy, useClientExternal:
		}

		return nil, nil
	}
}
