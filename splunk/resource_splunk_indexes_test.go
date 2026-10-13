package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/terraform-provider-splunk/client"
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
				// Recreate an index deleted outside Terraform, then keep the in-place update below.
				PreConfig: func() {
					deleteAndWaitForNotFound(t,
						func(c *client.Client) (*http.Response, error) {
							return c.DeleteIndexObject("new-index", "nobody", "system")
						},
						func(c *client.Client) url.URL {
							return c.BuildSplunkURL(nil, "services", "data", "indexes", "new-index")
						},
					)
				},
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

func TestIndexReadDropsStateWhenIndexIsGone(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	var count string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/servicesNS/-/-/data/indexes" {
			t.Errorf("path = %s, want the full index collection", r.URL.EscapedPath())
		}
		count = r.URL.Query().Get("count")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"entry":[]}`)
	}))
	defer server.Close()

	d := indexResourceData(t, "new-index")
	if err := indexRead(d, httpTestProvider(t, server)); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "" {
		t.Fatalf("id = %q, want empty", d.Id())
	}
	if count != "-1" {
		t.Fatalf("count = %q, want -1 so a missing name is not a truncated page", count)
	}
}

func TestIndexReadDropsStateWhenDetailIsNotFound(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	const listBody = `{"entry":[{"name":"new-index","acl":{"owner":"nobody","app":"system"},"content":{}}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == "/servicesNS/-/-/data/indexes" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, listBody)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	d := indexResourceData(t, "new-index")
	if err := indexRead(d, httpTestProvider(t, server)); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "" {
		t.Fatalf("id = %q, want empty", d.Id())
	}
}

func TestIndexReadKeepsStateWhenCreateCannotSeeIndex(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"entry":[]}`)
	}))
	defer server.Close()

	d := indexResourceData(t, "new-index")
	d.MarkNewResource()
	err := indexRead(d, httpTestProvider(t, server))
	if err == nil {
		t.Fatal("indexRead error = nil, want a missing-index error during create")
	}
	if d.Id() != "new-index" {
		t.Fatalf("id = %q, want new-index", d.Id())
	}
}

func TestIndexReadKeepsStateOnListError(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":[{"type":"ERROR","text":"unavailable"}]}`, http.StatusInternalServerError)
	}))
	defer server.Close()

	d := indexResourceData(t, "new-index")
	if err := indexRead(d, httpTestProvider(t, server)); err == nil {
		t.Fatal("indexRead error = nil, want the list failure")
	}
	if d.Id() != "new-index" {
		t.Fatalf("id = %q, want new-index", d.Id())
	}
}

func indexResourceData(t *testing.T, name string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, index().Schema, map[string]interface{}{
		"name": name,
	})
	d.SetId(name)
	return d
}

func testAccSplunkIndexDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
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
