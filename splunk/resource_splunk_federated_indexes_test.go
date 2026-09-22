package splunk

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/terraform-provider-splunk/client"
)

// TestFederatedIndexEntry verifies missing, valid, ambiguous, and malformed responses.
func TestFederatedIndexEntry(t *testing.T) {
	for _, test := range []struct {
		name     string
		body     string
		wantName string
		wantErr  bool
	}{
		{name: "empty", body: `{"entry":[]}`},
		{name: "single", body: `{"entry":[{"name":"federated:remote-main"}]}`, wantName: "federated:remote-main"},
		{name: "multiple", body: `{"entry":[{"name":"federated:first"},{"name":"federated:second"}]}`, wantErr: true},
		{name: "invalid JSON", body: `{"entry":`, wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp := &http.Response{Body: io.NopCloser(strings.NewReader(test.body))}
			defer func() { _ = resp.Body.Close() }()
			entry, err := federatedIndexEntry(resp)
			if (err != nil) != test.wantErr {
				t.Fatalf("error = %v, want error: %t", err, test.wantErr)
			}
			if test.wantName == "" {
				if entry != nil {
					t.Fatalf("expected no entry, got %+v", entry)
				}
			} else if entry == nil || entry.Name != test.wantName {
				t.Fatalf("entry = %+v, want name %q", entry, test.wantName)
			}
		})
	}
}

// TestFederatedIndexNameLifecycle verifies that configuration uses a short name,
// while reads, updates, deletes, and imports address Splunk's prefixed name.
func TestFederatedIndexNameLifecycle(t *testing.T) {
	const shortName = "remote-main"
	const fullName = "federated:remote-main"
	const collection = "/services/data/federated/index"
	var dataset string
	var created, deleted bool
	transport := federatedPasswordTransport(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodPost {
			data, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			form, err := url.ParseQuery(string(data))
			if err != nil {
				t.Fatal(err)
			}
			if !created {
				if req.URL.Path != collection || form.Get("name") != shortName {
					t.Fatal("create must submit the unprefixed name to the collection")
				}
				created = true
			} else if req.URL.Path != collection+"/"+fullName {
				t.Fatal("update must address the prefixed name")
			}
			dataset = form.Get("federated.dataset")
		} else {
			if req.URL.Path != collection+"/"+fullName {
				t.Fatal("read/delete must address the prefixed name")
			}
			if req.Method == http.MethodDelete {
				deleted = true
			} else if req.Method != http.MethodGet {
				t.Fatalf("unexpected method %s", req.Method)
			}
		}
		body, err := json.Marshal(map[string]interface{}{
			"entry": []interface{}{map[string]interface{}{
				"name":    fullName,
				"content": map[string]string{"federated.provider": "remote", "federated.dataset": dataset},
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
	r := federatedIndexes()
	state := &terraform.InstanceState{}
	for _, step := range []struct{ datasetType, name string }{
		{"", "main"}, // Omitted type defaults to index.
		{"index", "_internal"},
		{"savedsearch", "Daily report: errors"},
		{"lastjob", "Daily report: errors"},
		{"datamodel", "Network_Traffic"},
	} {
		raw := map[string]interface{}{
			"name": shortName, "federated_provider": "remote", "dataset_name": step.name,
		}
		expectedType := step.datasetType
		if expectedType == "" {
			expectedType = "index"
		} else {
			raw["dataset_type"] = expectedType
		}
		config := terraform.NewResourceConfigRaw(raw)
		diff, err := r.Diff(state, config, meta)
		if err != nil {
			t.Fatal(err)
		}
		if state.ID != "" && diff.RequiresNew() {
			t.Fatal("dataset updates must not replace the index")
		}
		state, err = r.Apply(state, diff, meta)
		if err != nil {
			t.Fatal(err)
		}
		if state.ID != fullName || state.Attributes["name"] != shortName {
			t.Fatal("state must preserve the short name and prefixed ID")
		}
		if dataset != expectedType+":"+step.name || state.Attributes["dataset_type"] != expectedType || state.Attributes["dataset_name"] != step.name {
			t.Fatal("dataset fields must round-trip through Splunk's combined mapping")
		}
		diff, err = r.Diff(state, config, meta)
		if err != nil {
			t.Fatal(err)
		}
		if diff != nil && !diff.Empty() {
			t.Fatal("refresh must not introduce a name diff")
		}
	}
	for _, id := range []string{shortName, fullName} {
		d := r.Data(nil)
		d.SetId(id)
		imported, err := r.Importer.State(d, meta)
		if err != nil {
			t.Fatal(err)
		}
		if err := r.Read(imported[0], meta); err != nil {
			t.Fatal(err)
		}
		if imported[0].Id() != fullName || imported[0].Get("name") != shortName {
			t.Fatal("import must normalize the ID and read the short name")
		}
		if imported[0].Get("dataset_type") != "datamodel" || imported[0].Get("dataset_name") != "Network_Traffic" {
			t.Fatal("import must populate the separate dataset fields")
		}
	}
	if _, err := r.Apply(state, &terraform.InstanceDiff{Destroy: true}, meta); err != nil {
		t.Fatal(err)
	}
	if !deleted {
		t.Fatal("expected a DELETE request")
	}
}

// TestFederatedIndexDatasetTypes checks the exact supported API values.
func TestFederatedIndexDatasetTypes(t *testing.T) {
	s := federatedIndexes().Schema["dataset_type"]
	if s.Default != "index" || s.ForceNew {
		t.Fatal("dataset_type must default to index and support in-place updates")
	}
	for _, value := range []string{"index", "savedsearch", "lastjob", "datamodel"} {
		if _, errs := s.ValidateFunc(value, "dataset_type"); len(errs) != 0 {
			t.Fatalf("valid type %q rejected: %v", value, errs)
		}
	}
	for _, value := range []string{"", "saved_searches", "INDEX", "unknown"} {
		if _, errs := s.ValidateFunc(value, "dataset_type"); len(errs) == 0 {
			t.Fatalf("invalid type %q accepted", value)
		}
	}
}
