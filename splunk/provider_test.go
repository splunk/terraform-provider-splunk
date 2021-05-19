package splunk

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
