package splunk

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestParseNamespacedRESTImportID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		resourceParts []string
		wantMatched   bool
		wantOwner     string
		wantApp       string
		wantName      string
		wantErr       bool
	}{
		{
			name:          "saved search path",
			id:            "/servicesNS/nobody/<app>/saved/searches/Example%20Saved%20Search",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantOwner:     "nobody",
			wantApp:       "<app>",
			wantName:      "Example Saved Search",
		},
		{
			name:          "view full url",
			id:            "https://example.invalid:8089/servicesNS/<owner>/<app>/data/ui/views/Example%20View?output_mode=json",
			resourceParts: []string{"data", "ui", "views"},
			wantMatched:   true,
			wantOwner:     "<owner>",
			wantApp:       "<app>",
			wantName:      "Example View",
		},
		{
			name:          "literal plus remains plus",
			id:            "/servicesNS/nobody/<app>/saved/searches/Example+Saved+Search",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantOwner:     "nobody",
			wantApp:       "<app>",
			wantName:      "Example+Saved+Search",
		},
		{
			name:          "encoded plus remains plus",
			id:            "/servicesNS/nobody/<app>/saved/searches/Example%2BSaved%2BSearch",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantOwner:     "nobody",
			wantApp:       "<app>",
			wantName:      "Example+Saved+Search",
		},
		{
			name:          "full url encoded plus remains plus",
			id:            "https://example.invalid:8089/servicesNS/nobody/<app>/saved/searches/Example%2BSaved%2BSearch",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantOwner:     "nobody",
			wantApp:       "<app>",
			wantName:      "Example+Saved+Search",
		},
		{
			name:          "encoded slash remains in name",
			id:            "/servicesNS/nobody/<app>/saved/searches/Example%2FSaved%20Search",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantOwner:     "nobody",
			wantApp:       "<app>",
			wantName:      "Example/Saved Search",
		},
		{
			name:          "bare name falls back to legacy import",
			id:            "Example Saved Search",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   false,
		},
		{
			name:          "bare servicesNS name falls back to legacy import",
			id:            "servicesNS",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   false,
		},
		{
			name:          "relative servicesNS path falls back to legacy import",
			id:            "prefix/servicesNS/nobody/<app>/saved/searches/Example%20Saved%20Search",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   false,
		},
		{
			name:          "wrong resource path errors",
			id:            "/servicesNS/nobody/<app>/data/ui/views/Example%20Saved%20Search",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantErr:       true,
		},
		{
			name:          "child endpoint errors",
			id:            "/servicesNS/nobody/<app>/saved/searches/Example%20Saved%20Search/acl",
			resourceParts: []string{"saved", "searches"},
			wantMatched:   true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, matched, err := parseNamespacedRESTImportID(tt.id, tt.resourceParts...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if matched != tt.wantMatched {
				t.Fatalf("matched = %v, want %v", matched, tt.wantMatched)
			}
			if !matched {
				return
			}
			if got.Owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", got.Owner, tt.wantOwner)
			}
			if got.App != tt.wantApp {
				t.Errorf("app = %q, want %q", got.App, tt.wantApp)
			}
			if got.Name != tt.wantName {
				t.Errorf("name = %q, want %q", got.Name, tt.wantName)
			}
		})
	}
}

func TestParseNamespacedConfigsConfImportID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		wantMatched bool
		wantOwner   string
		wantApp     string
		wantName    string
		wantErr     bool
	}{
		{
			name:        "configs conf path",
			id:          "/servicesNS/nobody/search_app_dashweb/configs/conf-props/httpevent",
			wantMatched: true,
			wantOwner:   "nobody",
			wantApp:     "search_app_dashweb",
			wantName:    "props/httpevent",
		},
		{
			name:        "configs conf full url with encoded stanza",
			id:          "https://example.invalid:8089/servicesNS/admin/search/configs/conf-inputs/sqs%3Aqueue",
			wantMatched: true,
			wantOwner:   "admin",
			wantApp:     "search",
			wantName:    "inputs/sqs:queue",
		},
		{
			name:        "bare name falls back to default namespace import",
			id:          "props/httpevent",
			wantMatched: false,
		},
		{
			name:        "wrong resource path errors",
			id:          "/servicesNS/nobody/search_app_dashweb/saved/conf-props/httpevent",
			wantMatched: true,
			wantErr:     true,
		},
		{
			name:        "missing configuration prefix errors",
			id:          "/servicesNS/nobody/search_app_dashweb/configs/props/httpevent",
			wantMatched: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, matched, err := parseNamespacedConfigsConfImportID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if matched != tt.wantMatched {
				t.Fatalf("matched = %v, want %v", matched, tt.wantMatched)
			}
			if !matched {
				return
			}
			if got.Owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", got.Owner, tt.wantOwner)
			}
			if got.App != tt.wantApp {
				t.Errorf("app = %q, want %q", got.App, tt.wantApp)
			}
			if got.Name != tt.wantName {
				t.Errorf("name = %q, want %q", got.Name, tt.wantName)
			}
		})
	}
}

func TestConfigsConfImportStateNamespacedPath(t *testing.T) {
	d := schema.TestResourceDataRaw(t, configsConf().Schema, nil)
	d.SetId("/servicesNS/nobody/search_app_dashweb/configs/conf-props/httpevent")

	states, err := configsConfImportState(d, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(states) != 1 {
		t.Fatalf("state count = %d, want 1", len(states))
	}
	if got, want := d.Id(), "props/httpevent"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("name").(string), "props/httpevent"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}

	acl := getACLConfig(d.Get("acl").([]interface{}))
	if got, want := acl.Owner, "nobody"; got != want {
		t.Errorf("acl owner = %q, want %q", got, want)
	}
	if got, want := acl.App, "search_app_dashweb"; got != want {
		t.Errorf("acl app = %q, want %q", got, want)
	}
	if got, want := acl.Sharing, "app"; got != want {
		t.Errorf("acl sharing = %q, want %q", got, want)
	}
}

func TestConfigsConfImportStateBareName(t *testing.T) {
	d := schema.TestResourceDataRaw(t, configsConf().Schema, nil)
	d.SetId("props/httpevent")

	_, err := configsConfImportState(d, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got, want := d.Id(), "props/httpevent"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("name").(string), "props/httpevent"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
}

func TestSavedSearchesImportStateNamespacedPath(t *testing.T) {
	d := schema.TestResourceDataRaw(t, savedSearches().Schema, nil)
	d.SetId("/servicesNS/nobody/<app>/saved/searches/Example%20Saved%20Search")

	states, err := savedSearchesImportState(d, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(states) != 1 {
		t.Fatalf("state count = %d, want 1", len(states))
	}
	if got, want := d.Id(), "Example Saved Search"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("name").(string), "Example Saved Search"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}

	acl := getACLConfig(d.Get("acl").([]interface{}))
	if got, want := acl.Owner, "nobody"; got != want {
		t.Errorf("acl owner = %q, want %q", got, want)
	}
	if got, want := acl.App, "<app>"; got != want {
		t.Errorf("acl app = %q, want %q", got, want)
	}
	if got, want := acl.Sharing, "app"; got != want {
		t.Errorf("acl sharing = %q, want %q", got, want)
	}
}

func TestSavedSearchesImportStateBareName(t *testing.T) {
	d := schema.TestResourceDataRaw(t, savedSearches().Schema, nil)
	d.SetId("Example Saved Search")

	_, err := savedSearchesImportState(d, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got, want := d.Id(), "Example Saved Search"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("name").(string), "Example Saved Search"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
}

func TestSplunkDashboardsImportStateNamespacedURL(t *testing.T) {
	d := schema.TestResourceDataRaw(t, splunkDashboards().Schema, nil)
	d.SetId("https://example.invalid:8089/servicesNS/<owner>/<app>/data/ui/views/Example%20View?output_mode=json")

	states, err := splunkDashboardsImportState(d, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(states) != 1 {
		t.Fatalf("state count = %d, want 1", len(states))
	}
	if got, want := d.Id(), "Example View"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("name").(string), "Example View"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}

	acl := getACLConfig(d.Get("acl").([]interface{}))
	if got, want := acl.Owner, "<owner>"; got != want {
		t.Errorf("acl owner = %q, want %q", got, want)
	}
	if got, want := acl.App, "<app>"; got != want {
		t.Errorf("acl app = %q, want %q", got, want)
	}
	if got, want := acl.Sharing, "user"; got != want {
		t.Errorf("acl sharing = %q, want %q", got, want)
	}
}
