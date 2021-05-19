package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const newMetricsIndex = `
resource "splunk_indexes" "new-metrics-index" {
    name = "new-metrics-index"
	datatype = "metric"
}
`

const updateMetricsIndex = `
resource "splunk_indexes" "new-metrics-index" {
	name = "new-metrics-index"
    max_time_unreplicated_no_acks = 301
}
`

func TestAccCreateSplunkMetricsIndex(t *testing.T) {
	resourceName := "splunk_indexes.new-metrics-index"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkMetricsIndexDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newMetricsIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "max_time_unreplicated_no_acks", "300"),
				),
			},
			{
				Config: updateMetricsIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "max_time_unreplicated_no_acks", "301"),
				),
			},
			{
				ResourceName:      "splunk_indexes.new-metrics-index",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkMetricsIndexDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_indexes":
			queryValues := url.Values{}
			queryValues.Add("datatype", "all")

			endpoint := client.BuildSplunkURL(queryValues, "services", "data", "indexes", rs.Primary.ID)
			// Index delete is asynchronous - brief wait to ensure deletion has completed
			time.Sleep(2 * time.Second)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
