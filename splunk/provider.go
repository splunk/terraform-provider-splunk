package splunk

import (
	"github.com/splunk/terraform-provider-splunk/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type SplunkProvider struct {
	Client *client.Client
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
		"splunk_indexes":                     index(),
		"splunk_configs_conf":                configsConf(),
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise splunk client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	provider := &SplunkProvider{}
	var splunkdClient *client.Client

	if token, ok := d.GetOk("auth_token"); ok {
		splunkdClient = client.NewSplunkdClientWithAuthToken(token.(string),
			[2]string{d.Get("username").(string), d.Get("password").(string)},
			d.Get("url").(string),
			client.NewSplunkdHTTPClient(
				time.Duration(d.Get("timeout").(int))*time.Second,
				d.Get("insecure_skip_verify").(bool)))
	} else {
		splunkdClient = client.NewSplunkdClient("",
			[2]string{d.Get("username").(string), d.Get("password").(string)},
			d.Get("url").(string),
			client.NewSplunkdHTTPClient(
				time.Duration(d.Get("timeout").(int))*time.Second,
				d.Get("insecure_skip_verify").(bool)))

		// Login is required to get session key
		err := splunkdClient.Login()
		if err != nil {
			return splunkdClient, err
		}
	}

	provider.Client = splunkdClient
	return provider, nil
}
