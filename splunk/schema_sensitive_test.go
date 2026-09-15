package splunk

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestCredentialAttributesAreSensitive(t *testing.T) {
	p := Provider().(*schema.Provider)

	providerAttrs := []string{"password", "auth_token"}
	for _, name := range providerAttrs {
		if !p.Schema[name].Sensitive {
			t.Errorf("provider schema attribute %q must be Sensitive", name)
		}
	}

	resourceAttrs := map[string][]string{
		"splunk_inputs_http_event_collector": {"token"},
		"splunk_inputs_tcp_splunk_tcp_token": {"token"},
		"splunk_outputs_tcp_group":           {"token"},
		"splunk_saved_searches": {
			"action_email_auth_password",
			"action_pagerduty_integration_key",
			"action_pagerduty_integration_key_override",
			"action_slack_param_webhook_url_override",
			"action_victorops_param_routing_key_override",
			"action_webhook_param_url",
		},
		"splunk_apps_local": {"auth", "session"},
		// Existing sensitive attributes should stay marked.
		"splunk_authentication_users": {"password"},
		"splunk_outputs_tcp_server":   {"ssl_password"},
		"splunk_inputs_tcp_ssl":      {"password"},
	}

	for resourceName, attrs := range resourceAttrs {
		res, ok := p.ResourcesMap[resourceName]
		if !ok {
			t.Errorf("resource %q not registered", resourceName)
			continue
		}
		for _, name := range attrs {
			s, ok := res.Schema[name]
			if !ok {
				t.Errorf("resource %q is missing schema attribute %q", resourceName, name)
				continue
			}
			if !s.Sensitive {
				t.Errorf("resource %q attribute %q must be Sensitive", resourceName, name)
			}
		}
	}
}
