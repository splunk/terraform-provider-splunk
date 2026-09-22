package splunk

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/terraform-provider-splunk/client"
)

// federatedPasswordTransport checks real client requests without a live Splunk.
type federatedPasswordTransport func(*http.Request) (*http.Response, error)

func (f federatedPasswordTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// TestFederatedProviderPasswordLifecycle exercises SDK diff/apply, HTTP requests,
// and persisted state so hashing cannot silently break creation or rotation.
func TestFederatedProviderPasswordLifecycle(t *testing.T) {
	var posts int
	var wantPassword, wantPath, appContext, mode string
	transport := federatedPasswordTransport(func(req *http.Request) (*http.Response, error) {
		switch req.Method {
		case http.MethodPost:
			posts++
			if req.URL.Path != wantPath {
				t.Errorf("POST path = %q, want %q", req.URL.Path, wantPath)
			}
			data, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			form, err := url.ParseQuery(string(data))
			if err != nil {
				t.Fatal(err)
			}
			if form.Get("password") != wantPassword {
				t.Error("request did not contain the expected original password")
			}
			if _, exists := form["password"]; wantPassword == "" && exists {
				t.Error("unchanged password must be omitted")
			}
			appContext = form.Get("appContext")
			mode = form.Get("mode")
			if _, exists := form["appContext"]; mode == "transparent" && exists {
				t.Error("transparent mode must omit appContext")
			}
		case http.MethodGet:
		default:
			t.Errorf("unexpected method %s: password rotation must not replace the resource", req.Method)
		}
		body, err := json.Marshal(map[string]interface{}{
			"entry": []interface{}{map[string]interface{}{
				"name": "remote",
				"content": map[string]string{
					"type": "splunk", "mode": mode, "appContext": appContext,
					"hostPort": "remote:8089", "serviceAccount": "federated",
				},
			}},
		})
		if err != nil {
			t.Fatal(err)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}, nil
	})
	c, err := client.NewSplunkdClient("", [2]string{}, "local:8089", "", &http.Client{Transport: transport})
	if err != nil {
		t.Fatal(err)
	}
	meta := &SplunkProvider{Client: c}
	resource := federatedProviders()
	state := &terraform.InstanceState{}
	for _, step := range []struct {
		name, password, app, mode string
		sendPassword              bool
	}{
		{"create", "first-test-password", "search", "standard", true},
		{"rotate", "second-test-password", "search", "standard", true},
		{"unrelated update", "second-test-password", "other_app", "standard", false},
		{"transparent mode", "second-test-password", "other_app", "transparent", false},
		{"restore standard mode", "second-test-password", "other_app", "standard", false},
	} {
		t.Run(step.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name": "remote", "host_port": "remote:8089", "service_account": "federated",
				"password": step.password, "app_context": step.app, "mode": step.mode,
			})
			diff, err := resource.Diff(state, config, meta)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil || diff.Empty() {
				t.Fatal("expected a configuration change")
			}
			if step.name != "create" && diff.RequiresNew() {
				t.Fatal("update must not require replacement")
			}
			wantPassword = ""
			if step.sendPassword {
				wantPassword = step.password
			}
			wantPath = "/services/data/federated/provider/remote"
			if step.name == "create" {
				wantPath = "/services/data/federated/provider"
			}
			before := posts
			state, err = resource.Apply(state, diff, meta)
			if err != nil {
				t.Fatal(err)
			}
			if posts != before+1 || state.ID != "remote" {
				t.Fatal("expected one write and a stable resource ID")
			}
			expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(step.password)))
			if state.Attributes["password"] != expectedHash {
				t.Fatal("state must store exactly the password hash")
			}
			diff, err = resource.Diff(state, config, meta)
			if err != nil {
				t.Fatal(err)
			}
			if diff != nil && !diff.Empty() {
				t.Fatal("unchanged configuration must produce no diff")
			}
		})
	}
}
