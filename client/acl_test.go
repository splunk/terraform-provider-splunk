package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestGetAcl_CloudMode_QueryStringIncludesOwnerAndSharing verifies that with ACLGetModeCloud,
// GET .../acl includes owner and sharing query parameters (Splunk Cloud / issue #224).
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

// TestGetAcl_EnterpriseMode_OmitsOwnerSharingFromQuery verifies default Enterprise behavior:
// no owner/sharing in the query string (only output_mode / count from BuildSplunkURL).
func TestGetAcl_EnterpriseMode_OmitsOwnerSharingFromQuery(t *testing.T) {
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
	c.ACLGetMode = ACLGetModeEnterprise

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
		t.Errorf("enterprise mode should not set owner query param, got %q", q.Get("owner"))
	}
	if q.Get("sharing") != "" {
		t.Errorf("enterprise mode should not set sharing query param, got %q", q.Get("sharing"))
	}
}
