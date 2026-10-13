package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestIndexReadAcceptsNumberAndStringJSON(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	// The collection and the detail response use opposite JSON types so a
	// decode failure on either response fails the read.
	const listBody = `{"entry":[{"name":"new-index","acl":{"owner":"nobody","app":"system"},"content":{
		"bucketRebuildMemoryHint":"16MB",
		"maxHotBuckets":6,
		"repFactor":"auto"
	}}]}`
	const detailBody = `{"entry":[{"name":"new-index","acl":{"owner":"nobody","app":"system"},"content":{
		"bucketRebuildMemoryHint":16777216,
		"maxHotBuckets":"auto",
		"repFactor":0
	}}]}`

	d := indexResourceData(t, "new-index")
	if err := indexRead(d, indexServer(t, listBody, detailBody)); err != nil {
		t.Fatal(err)
	}
	if got := d.Get("bucket_rebuild_memory_hint").(string); got != "16777216" {
		t.Errorf("bucket_rebuild_memory_hint = %q, want the integer Splunk returned", got)
	}
	if got := d.Get("max_hot_buckets").(int); got != 3 {
		t.Errorf("max_hot_buckets = %d, want 3 from auto", got)
	}
	if got := d.Get("rep_factor").(string); got != "0" {
		t.Errorf("rep_factor = %q, want 0", got)
	}
}

func TestIndexReadKeepsConfiguredSpellings(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	const listBody = `{"entry":[{"name":"new-index","acl":{"owner":"nobody","app":"system"},"content":{
		"bucketRebuildMemoryHint":"auto",
		"maxHotBuckets":"6",
		"repFactor":"0"
	}}]}`
	const detailBody = `{"entry":[{"name":"new-index","acl":{"owner":"nobody","app":"system"},"content":{
		"bucketRebuildMemoryHint":16777216,
		"maxHotBuckets":"6",
		"repFactor":0
	}}]}`

	d := schema.TestResourceDataRaw(t, index().Schema, map[string]interface{}{
		"name":                       "new-index",
		"bucket_rebuild_memory_hint": "16MB",
		"max_hot_buckets":            6,
		"rep_factor":                 "auto",
	})
	d.SetId("new-index")
	if err := indexRead(d, indexServer(t, listBody, detailBody)); err != nil {
		t.Fatal(err)
	}
	if got := d.Get("bucket_rebuild_memory_hint").(string); got != "16MB" {
		t.Errorf("bucket_rebuild_memory_hint = %q, want 16MB", got)
	}
	if got := d.Get("max_hot_buckets").(int); got != 6 {
		t.Errorf("max_hot_buckets = %d, want 6", got)
	}
	if got := d.Get("rep_factor").(string); got != "auto" {
		t.Errorf("rep_factor = %q, want auto", got)
	}
}

func TestIndexReadDecodeErrorNamesField(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	const listBody = `{"entry":[{"name":"new-index","acl":{"owner":"nobody","app":"system"},"content":{
		"maxHotBuckets":{"unexpected":true}
	}}]}`
	d := indexResourceData(t, "new-index")
	err := indexRead(d, indexServer(t, listBody, `{"entry":[]}`))
	if err == nil {
		t.Fatal("indexRead error = nil, want a decode error")
	}
	if !strings.Contains(err.Error(), "maxHotBuckets") {
		t.Fatalf("indexRead error = %q, want it to name maxHotBuckets", err)
	}
	if d.Id() != "new-index" {
		t.Fatalf("id = %q, want new-index", d.Id())
	}
}

func TestGetIndexConfigByNameErrorResponses(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{name: "splunk message", status: 400, body: `{"messages":[{"type":"ERROR","text":"index is in use"}]}`, want: "index is in use"},
		{name: "empty messages", status: 500, body: `{}`, want: "HTTP 500"},
		{name: "invalid json", status: 502, body: `not-json`, want: "HTTP 502"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.status,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}
			_, err := getIndexConfigByName("new-index", resp)
			if err == nil {
				t.Fatal("error = nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestRepFactorAndBucketHintDiffSuppression(t *testing.T) {
	if suppressRepFactorDiff("rep_factor", "", "auto", nil) {
		t.Error("empty to auto should plan so the first apply sends auto")
	}
	if suppressRepFactorDiff("rep_factor", "0", "auto", nil) != true {
		t.Error("0 and auto should not plan a change")
	}
	if suppressRepFactorDiff("rep_factor", "auto", "1", nil) {
		t.Error("auto and 1 should plan a change")
	}
	if suppressBucketRebuildMemoryHintDiff("bucket_rebuild_memory_hint", "16MB", "16777216", nil) != true {
		t.Error("16MB and 16777216 should not plan a change")
	}
	if suppressBucketRebuildMemoryHintDiff("bucket_rebuild_memory_hint", "auto", "67108864", nil) {
		t.Error("auto and a byte count should plan a change")
	}
}

func indexServer(t *testing.T, listBody, detailBody string) *SplunkProvider {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.EscapedPath() {
		case "/servicesNS/-/-/data/indexes":
			_, _ = io.WriteString(w, listBody)
		case "/servicesNS/nobody/system/data/indexes/new-index":
			_, _ = io.WriteString(w, detailBody)
		default:
			t.Errorf("unexpected path %s", r.URL.EscapedPath())
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return httpTestProvider(t, server)
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
