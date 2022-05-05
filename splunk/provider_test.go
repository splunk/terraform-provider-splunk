package splunk

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/go-splunk-client/pkg/authenticators"
	externalClient "github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/terraform-provider-splunk/client"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"splunk": testAccProvider,
	}
}

func newTestClient() (*client.Client, error) {
	http, err := client.NewSplunkdHTTPClient(30*time.Second, true)
	if err != nil {
		return nil, err
	}
	return client.NewSplunkdClient(
		"",
		[2]string{os.Getenv("SPLUNK_USERNAME"),
			os.Getenv("SPLUNK_PASSWORD")},
		os.Getenv("SPLUNK_URL"),

		http)
}

func newTestExternalClient() (*externalClient.Client, error) {
	externalURL, err := canonicalizeSplunkURL(os.Getenv("SPLUNK_URL"))
	if err != nil {
		return nil, err
	}

	return &externalClient.Client{
		URL: externalURL,
		Authenticator: &authenticators.Password{
			Username: os.Getenv("SPLUNK_USERNAME"),
			Password: os.Getenv("SPLUNK_PASSWORD"),
		},
		TLSInsecureSkipVerify: true,
		Timeout:               30 * time.Second,
	}, nil
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SPLUNK_HOME"); v == "" {
		t.Fatal("SPLUNK_HOME must be set for acceptance tests")
	}
	if v := os.Getenv("SPLUNK_URL"); v == "" {
		t.Fatal("SPLUNK_URL must be set for acceptance tests")
	}
	if v := os.Getenv("SPLUNK_USERNAME"); v == "" {
		t.Fatal("SPLUNK_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("SPLUNK_PASSWORD"); v == "" {
		t.Fatal("SPLUNK_PASSWORD must be set for acceptance tests")
	}
}

func TestProvider_canonicalizeSplunkURL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{
			name:  "local",
			input: "localhost:8089",
			want:  "https://localhost:8089",
		},
		{
			name:  "domain, port 8089",
			input: "splunk.example.com:8089",
			want:  "https://splunk.example.com:8089",
		},
		{
			name:  "domain, no port",
			input: "splunk.example.com",
			want:  "https://splunk.example.com",
		},
		{
			name:  "local with scheme",
			input: "https://localhost:8089",
			want:  "https://localhost:8089",
		},
		{
			name:  "domain with scheme",
			input: "https://splunk.example.com:8089",
			want:  "https://splunk.example.com:8089",
		},
	}

	for _, test := range tests {
		got, err := canonicalizeSplunkURL(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: canonicalizeSplunkURL() returned error? %v (%s)", test.name, gotError, err)
		}

		if got != test.want {
			t.Errorf("%s: canonicalizeSplunkURL() got\n%s, want\n%s", test.name, got, test.want)
		}
	}
}
