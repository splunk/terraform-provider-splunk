package client

import (
	"errors"
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

func TestUpdateAclReturnsHTTPStatus(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	acl := &models.ACLObject{Sharing: "app", Owner: "admin", App: "search"}
	tests := []struct {
		name  string
		cloud bool
		code  int
	}{
		{"enterprise 404", false, http.StatusNotFound},
		{"enterprise 400", false, http.StatusBadRequest},
		{"cloud 404", true, http.StatusNotFound},
		{"cloud 403", true, http.StatusForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "denied", tt.code)
			}))
			defer server.Close()

			backend, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			c, err := NewSplunkdClient("", defaultAuth, backend.Host, "", server.Client())
			if err != nil {
				t.Fatal(err)
			}
			if tt.cloud {
				c.ACLGetMode = ACLGetModeCloud
			}

			err = c.UpdateAcl("admin", "search", "alert", acl, "saved", "searches")
			var statusErr *HTTPStatusError
			if !errors.As(err, &statusErr) || statusErr.StatusCode != tt.code {
				t.Fatalf("err = %v, want HTTPStatusError %d", err, tt.code)
			}
			if !strings.Contains(err.Error(), "POST failed for endpoint") {
				t.Fatalf("err = %q, want the existing POST failed message", err)
			}
		})
	}
}

func TestUpdateAclTransportErrorHasNoStatus(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	c, err := NewSplunkdClient("", defaultAuth, "splunk.example:8089", "", &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("dial failed")
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = c.UpdateAcl("admin", "search", "alert", &models.ACLObject{Sharing: "app"}, "saved", "searches")
	var statusErr *HTTPStatusError
	if err == nil || errors.As(err, &statusErr) {
		t.Fatalf("err = %v, want a transport error without an HTTP status", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
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
			want: "/servicesNS/nobody/app/saved/searches/A%2FB%2BC%20D/acl",
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
			want: "/servicesNS/nobody/app/saved/searches/A%2FB%2BC%20D/acl",
			call: func(c *Client) error {
				return c.UpdateAcl("nobody", "app", "A/B+C D", acl, "saved", "searches")
			},
		},
		{
			name: "update dashboard view acl",
			want: "/servicesNS/nobody/app/data/ui/views/A%2FB%2BC%20D/acl",
			call: func(c *Client) error {
				return c.UpdateAcl("nobody", "app", "A/B+C D", acl, "data", "ui", "views")
			},
		},
		{
			name: "update conf stanza acl with trailing slash",
			want: "/servicesNS/nobody/system/configs/conf-inputs/monitor:%2F%2F%2Fvar%2Flog%2F/acl",
			call: func(c *Client) error {
				return c.UpdateAcl("nobody", "system", "monitor:///var/log/", acl, "configs", "conf-inputs")
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
