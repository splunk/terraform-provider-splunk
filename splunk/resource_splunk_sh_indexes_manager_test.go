package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const newShIndexesManager = `
resource "splunk_sh_indexes_manager" "tf-index" {
    name = "tf-test-index-0"
    datatype = "event"
    frozen_time_period_in_secs = "94608000"
    max_global_raw_data_size_mb = "1000"
}
`

func TestAccCreateSplunkShIndexesManager(t *testing.T) {
	t.Skip("Test requires a Splunk Cloud instance")
	resourceName := "splunk_sh_indexes_manager.tf-index"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newShIndexesManager,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "datatype", "event"),
					resource.TestCheckResourceAttr(resourceName, "frozen_time_period_in_secs", "94608000"),
					resource.TestCheckResourceAttr(resourceName, "max_global_raw_data_size_mb", "1000"),
				),
			},
			{
				ResourceName:      "splunk_sh_indexes_manager.tf-index",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
