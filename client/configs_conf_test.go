package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/splunk/terraform-provider-splunk/client/models"
)

func TestSplitConfStanza(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		wantConf   string
		wantStanza string
	}{
		{"props/custom_stanza", "props", "custom_stanza"},
		{"inputs/monitor:///data/syslog/", "inputs", "monitor:///data/syslog/"},
		{"inputs/splunktcp://9997", "inputs", "splunktcp://9997"},
		{"props/source::/data/json/...", "props", "source::/data/json/..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, stanza := client.SplitConfStanza(tt.name)
			if conf != tt.wantConf {
				t.Errorf("conf = %q, want %q", conf, tt.wantConf)
			}
			if stanza != tt.wantStanza {
				t.Errorf("stanza = %q, want %q", stanza, tt.wantStanza)
			}
		})
	}
}

func TestUpdateConfigsConfObjectEscapesStanza(t *testing.T) {
	os.Setenv(envVarHTTPScheme, "http")
	defer os.Unsetenv(envVarHTTPScheme)

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

	conf := &models.ConfigsConfObject{Variables: map[string]string{"disabled": "false"}}
	if err := client.UpdateConfigsConfObject("inputs/splunktcp://9997", "nobody", "system", conf); err != nil {
		t.Fatalf("update failed: %s", err)
	}

	if want := "/servicesNS/nobody/system/configs/conf-inputs/splunktcp:%2F%2F9997"; got != want {
		t.Errorf("escaped path invalid, got %s, want %s", got, want)
	}
}
