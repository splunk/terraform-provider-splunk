package splunk

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestFederatedSearchSchemas(t *testing.T) {
	providerSchema := federatedProviders().Schema
	if !providerSchema["password"].Sensitive {
		t.Fatal("federated provider password must be sensitive")
	}
	if providerSchema["password"].ForceNew {
		t.Fatal("federated provider password changes must be updated in place")
	}
	if providerSchema["password"].StateFunc == nil {
		t.Fatal("federated provider password must be hashed in state")
	}
	if !providerSchema["name"].ForceNew {
		t.Fatal("federated provider name must require replacement")
	}
	if _, errs := providerSchema["password"].ValidateFunc("", "password"); len(errs) == 0 {
		t.Fatal("empty passwords must be rejected before an API request")
	}
	if _, errs := providerSchema["password"].ValidateFunc("valid-test-password", "password"); len(errs) != 0 {
		t.Fatal("nonempty passwords must be accepted")
	}

	indexSchema := federatedIndexes().Schema
	for _, name := range []string{"remote-main", "remote_main", "0remote", strings.Repeat("a", 2048)} {
		if _, errors := indexSchema["name"].ValidateFunc(name, "name"); len(errors) != 0 {
			t.Errorf("expected %q to be valid: %v", name, errors)
		}
	}
	for _, name := range []string{"", "federated:remote-main", "Remote", "kvstore", "_remote", "remote/main", strings.Repeat("a", 2049)} {
		if _, errors := indexSchema["name"].ValidateFunc(name, "name"); len(errors) == 0 {
			t.Errorf("expected %q to be invalid", name)
		}
	}
}

func TestAccFederatedSearch(t *testing.T) {
	remoteHostPort := os.Getenv("SPLUNK_FEDERATED_PROVIDER_HOST_PORT")
	remoteUsername := os.Getenv("SPLUNK_FEDERATED_PROVIDER_USERNAME")
	remotePassword := os.Getenv("SPLUNK_FEDERATED_PROVIDER_PASSWORD")
	if remoteHostPort == "" || remoteUsername == "" || remotePassword == "" {
		// Like other environment-specific acceptance tests, this is opt-in and
		// does not add credentials or infrastructure requirements to normal CI.
		t.Skip("set SPLUNK_FEDERATED_PROVIDER_HOST_PORT, SPLUNK_FEDERATED_PROVIDER_USERNAME, and SPLUNK_FEDERATED_PROVIDER_PASSWORD to run this test")
	}

	const providerResource = "splunk_federated_providers.remote"
	const indexResource = "splunk_federated_indexes.remote_main"
	config := `
resource "splunk_federated_providers" "remote" {
	app_context = "search"
  name            = "tf-acc-remote"
  host_port       = ` + strconv.Quote(remoteHostPort) + `
  service_account = ` + strconv.Quote(remoteUsername) + `
  password        = ` + strconv.Quote(remotePassword) + `
}

resource "splunk_federated_indexes" "remote_main" {
  name               = "tf-acc-remote-main"
  federated_provider = splunk_federated_providers.remote.name
  dataset_name       = "main"
}
`

	// Exercise both update endpoints using built-in Splunk apps and indexes.
	updatedConfig := strings.Replace(config, "app_context = \"search\"", "app_context = \"launcher\"", 1)
	updatedConfig = strings.Replace(updatedConfig, "\"main\"", "\"_internal\"", 1)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFederatedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(providerResource, "mode", "standard"),
					resource.TestCheckResourceAttr(providerResource, "app_context", "search"),
					resource.TestCheckResourceAttr(indexResource, "dataset_type", "index"),
					resource.TestCheckResourceAttr(indexResource, "dataset_name", "main"),
					resource.TestCheckResourceAttr(indexResource, "name", "tf-acc-remote-main"),
					resource.TestCheckResourceAttr(indexResource, "id", "federated:tf-acc-remote-main"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(providerResource, "app_context", "launcher"),
					resource.TestCheckResourceAttr(indexResource, "dataset_type", "index"),
					resource.TestCheckResourceAttr(indexResource, "dataset_name", "_internal"),
				),
			},
			{
				ResourceName:            providerResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				ResourceName:      indexResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccFederatedSearchDestroy verifies deletion rather than treating an
// authentication or transport failure as evidence that a resource is absent.
func testAccFederatedSearchDestroy(state *terraform.State) error {
	c, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range state.RootModule().Resources {
		var resp *http.Response
		switch rs.Type {
		case "splunk_federated_providers":
			resp, err = c.ReadFederatedProvider(rs.Primary.ID)
		case "splunk_federated_indexes":
			resp, err = c.ReadFederatedIndex(rs.Primary.ID)
		default:
			continue
		}
		if resp != nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				continue
			}
		}
		if err != nil {
			return fmt.Errorf("verify deletion of %s: %w", rs.Type, err)
		}
		return fmt.Errorf("%s %q still exists after destroy", rs.Type, rs.Primary.ID)
	}
	return nil
}
