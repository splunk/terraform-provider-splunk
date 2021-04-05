package splunk

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const newIndex = `
resource "splunk_indexes" "new-index" {
    name = "new-index"
}
`

const updateIndex = `
resource "splunk_indexes" "new-index" {
	name = "new-index"
    max_time_unreplicated_no_acks = 301
}
`

func TestAccCreateSplunkIndex(t *testing.T) {
	resourceName := "splunk_indexes.new-index"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkIndexDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "max_time_unreplicated_no_acks", "300"),
				),
			},
			{
				Config: updateIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "max_time_unreplicated_no_acks", "301"),
				),
			},
			{
				ResourceName:      "splunk_indexes.new-index",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkIndexDestroyResources(s *terraform.State) error {
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_indexes":
			endpoint := client.BuildSplunkURL(nil, "services", "data", "indexes", rs.Primary.ID)
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

func testResourceExampleInstanceStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"max_hot_buckets": 3,
	}
}

func testResourceExampleInstanceStateDataV1() map[string]interface{} {
	v0 := testResourceExampleInstanceStateDataV0()
	return map[string]interface{}{
		"max_hot_buckets": strconv.Itoa(v0["max_hot_buckets"].(int)),
	}
}

func TestResourceExampleInstanceStateUpgradeV0(t *testing.T) {
	expected := testResourceExampleInstanceStateDataV1()
	actual, err := resourceExampleInstanceStateUpgradeV0(testResourceExampleInstanceStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

