package splunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const newLookupTableFile = `
resource "splunk_lookup_table_file" "test" {
    app = "search"
	owner = "nobody"
	file_name = "lookup.csv"
	file_contents = <<-EOT
[
  ["status", "status_description", "status_type"],
  ["100", "Continue", "Informational"],
  ["101", "Switching Protocols", "Informational"]
]
EOT
}
`

const updatedLookupTableFile = `
resource "splunk_lookup_table_file" "test" {
    app = "search"
	owner = "nobody"
	file_name = "lookup.csv"
	file_contents = <<-EOT
[
  ["status", "status_description", "status_type"],
  ["100", "Continue", "Informational"],
  ["101", "Switching Protocols", "Informational"],
  ["200", "OK", "Successful"]
]
EOT
}
`

func TestAccSplunkLookupTableFile(t *testing.T) {
	resourceName := "splunk_lookup_table_file.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkLookupTableFileDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newLookupTableFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app", "search"),
					resource.TestCheckResourceAttr(resourceName, "owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "file_name", "lookup.csv"),
					resource.TestCheckResourceAttr(resourceName, "file_contents", "[\n  [\"status\", \"status_description\", \"status_type\"],\n  [\"100\", \"Continue\", \"Informational\"],\n  [\"101\", \"Switching Protocols\", \"Informational\"]\n]\n"),
				),
			},
			{
				Config: updatedLookupTableFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app", "search"),
					resource.TestCheckResourceAttr(resourceName, "owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "file_name", "lookup.csv"),
					resource.TestCheckResourceAttr(resourceName, "file_contents", "[\n  [\"status\", \"status_description\", \"status_type\"],\n  [\"100\", \"Continue\", \"Informational\"],\n  [\"101\", \"Switching Protocols\", \"Informational\"],\n  [\"200\", \"OK\", \"Successful\"]\n]\n"),
				),
			},
		},
	})
}

func testAccSplunkLookupTableFileDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_lookup_table_file":
			resp, err := client.ReadLookupTableFile("lookup.csv", "nobody", "search")
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
