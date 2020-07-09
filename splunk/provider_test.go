package splunk

import (
	"crypto/tls"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
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