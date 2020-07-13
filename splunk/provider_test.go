package splunk

import (
	"crypto/tls"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"os"
	"testing"
	"time"
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


func newTestClient() *SplunkClient {
	client := &SplunkClient{}
	client.HttpClient = &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client.Url = "https://localhost:8089"
	client.Username = "admin"
	client.Password = "changeme"
	return client
}

func testAccPreCheck(t *testing.T) {
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