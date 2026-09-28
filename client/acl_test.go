package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/splunk/terraform-provider-splunk/client/models"
)

// TestGetAcl_CloudMode_QueryStringIncludesOwnerAndSharing verifies owner and sharing appear on the GET query when ACLGetMode is cloud.
func TestGetAcl_CloudMode_QueryStringIncludesOwnerAndSharing(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	var got string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"entry":[{"name":"acl","acl":{"app":"system","owner":"nobody","sharing":"app","perms":{"read":[],"write":[]}}}],"messages":[]}`))
	}))
	defer ts.Close()

	backend, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewSplunkdClient("", defaultAuth, backend.Host, "", ts.Client())
	if err != nil {
		t.Fatal(err)
	}
	c.ACLGetMode = ACLGetModeCloud

	resp, err := c.GetAcl("nobody", "system", "myapp", "app", "apps", "local")
	if err != nil {
		t.Fatalf("GetAcl: %v", err)
	}
	defer resp.Body.Close()

	q, err := url.ParseQuery(got)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if q.Get("owner") != "nobody" {
		t.Errorf("query owner: got %q, want nobody", q.Get("owner"))
	}
	if q.Get("sharing") != "app" {
		t.Errorf("query sharing: got %q, want app", q.Get("sharing"))
	}
	if q.Get("output_mode") != "json" {
		t.Errorf("query output_mode: got %q, want json", q.Get("output_mode"))
	}
	if !strings.Contains(got, "count=") {
		t.Errorf("query should include count: %q", got)
	}
}

// TestGetAcl_NonCloudMode_OmitsOwnerSharingFromQuery verifies owner/sharing are not query parameters unless mode is cloud.
func TestGetAcl_NonCloudMode_OmitsOwnerSharingFromQuery(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	var got string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"entry":[{"name":"acl","acl":{"app":"system","owner":"nobody","sharing":"app","perms":{"read":[],"write":[]}}}],"messages":[]}`))
	}))
	defer ts.Close()

	backend, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewSplunkdClient("", defaultAuth, backend.Host, "", ts.Client())
	if err != nil {
		t.Fatal(err)
	}
	c.ACLGetMode = ""

	resp, err := c.GetAcl("nobody", "system", "myapp", "app", "apps", "local")
	if err != nil {
		t.Fatalf("GetAcl: %v", err)
	}
	defer resp.Body.Close()

	q, err := url.ParseQuery(got)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if q.Get("owner") != "" {
		t.Errorf("non-cloud mode should not set owner query param, got %q", q.Get("owner"))
	}
	if q.Get("sharing") != "" {
		t.Errorf("non-cloud mode should not set sharing query param, got %q", q.Get("sharing"))
	}
}

func TestAclURLsEscapeResourceName(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	acl := &models.ACLObject{
		Sharing: "app",
		Perms: models.Permissions{
			Read:  []string{"*"},
			Write: []string{"admin"},
		},
	}

	tests := []struct {
		name string
		want string
		call func(*Client) error
	}{
		{
			name: "get saved search acl",
			want: "/servicesNS/nobody/app/saved/searches/A%2FB+C%20D/acl",
			call: func(c *Client) error {
				resp, err := c.GetAcl("nobody", "app", "A/B+C D", "app", "saved", "searches")
				if resp != nil {
					defer resp.Body.Close()
				}
				return err
			},
		},
		{
			name: "update saved search acl",
			want: "/servicesNS/nobody/app/saved/searches/A%2FB+C%20D/acl",
			call: func(c *Client) error {
				return c.UpdateAcl("nobody", "app", "A/B+C D", acl, "saved", "searches")
			},
		},
		{
			name: "update dashboard view acl",
			want: "/servicesNS/nobody/app/data/ui/views/A%2FB+C%20D/acl",
			call: func(c *Client) error {
				return c.UpdateAcl("nobody", "app", "A/B+C D", acl, "data", "ui", "views")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.URL.EscapedPath()
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer server.Close()

			backend, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			client, err := NewSplunkdClient("", defaultAuth, backend.Host, "", server.Client())
			if err != nil {
				t.Fatal(err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("ACL call failed: %s", err)
			}

			if got != tt.want {
				t.Errorf("escaped path invalid, got %s, want %s", got, tt.want)
			}
		})
	}
}
