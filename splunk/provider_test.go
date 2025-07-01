package splunk

import (
	"net/url"
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

func getTestProviderURL(inputData map[string]interface{}, t *testing.T) url.URL {
	providerSchema := providerSchema()
	// Create a mock *schema.ResourceData
	resourceData := schema.TestResourceDataRaw(t, providerSchema, inputData)
	// Call the providerConfigure method
	providerInterface, err := providerConfigure(resourceData)
	if err != nil {
		t.Error(err)
	}
	provider, ok := providerInterface.(*SplunkProvider) // type assertion
	if !ok {
		t.Error(ok)
	}

	url := provider.Client.BuildSplunkURL(nil)
	return url
}

func TestProviderConfigure(t *testing.T) {
	// Define the input data for the ResourceData
	inputData := map[string]interface{}{
		"url":                  "localhost",
		"timeout":              60,
		"insecure_skip_verify": true,
		"auth_token":           "aa",
	}

	url := getTestProviderURL(inputData, t)
	if got, want := url.Host, "localhost"; got != want {
		t.Errorf("hostname invalid, got %s, want %s", got, want)
	}

	inputData["url"] = "localhost:8089"
	url = getTestProviderURL(inputData, t)
	if got, want := url.Host, "localhost:8089"; got != want {
		t.Errorf("url.Host invalid, got %s, want %s", got, want)
	}

	inputData["url"] = "https://localhost:8089"
	url = getTestProviderURL(inputData, t)
	if got, want := url.Host, "localhost:8089"; got != want {
		t.Errorf("url.Host invalid, got %s, want %s", got, want)
	}
	if got, want := url.Scheme, "https"; got != want {
		t.Errorf("url.Scheme invalid, got %s, want %s", got, want)
	}
	// Test URL with Path
	inputData["url"] = "https://localhost:8089/test/path"
	url = getTestProviderURL(inputData, t)
	if got, want := url.Path, "/test/path"; got != want {
		t.Errorf("url path invalid, got %s, want %s", got, want)
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
	// Parse splunk url host
	SPLUNK_URL := (os.Getenv("SPLUNK_URL"))
	if !hasScheme(SPLUNK_URL) {
		SPLUNK_URL = "http://" + SPLUNK_URL // add http scheme so url.Parse works
	}
	u, err := url.Parse(SPLUNK_URL)
	if err != nil {
		return nil, err
	}
	host := u.Host

	return client.NewSplunkdClient(
		"",
		[2]string{os.Getenv("SPLUNK_USERNAME"),
			os.Getenv("SPLUNK_PASSWORD")},
		host,
		"",
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
