package splunk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	splunkclient "github.com/splunk/terraform-provider-splunk/client"
)

const newLookupTableFile = `
resource "splunk_lookup_table_file" "test" {
    app = "search"
	owner = "nobody"
	file_name = "lookup.csv"
	file_contents = [
		["status", "status_description", "status_type"],
		["100", "Continue", "Informational"],
		["101", "Switching Protocols", "Informational"]
	]
}
`

const updatedLookupTableFile = `
resource "splunk_lookup_table_file" "test" {
    app = "search"
	owner = "nobody"
	file_name = "lookup.csv"
	file_contents = [
		["status", "status_description", "status_type"],
		["100", "Continue", "Informational"],
		["101", "Switching Protocols", "Informational"],
		["200", "OK", "Successful"]
	]
}
`

func TestLookupTableFileContentsFromPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lookup.csv")
	if err := os.WriteFile(path, []byte("id,message\n1,\"hello, world\"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	d := schema.TestResourceDataRaw(t, lookupTableFile().Schema, map[string]interface{}{
		"app":       "search",
		"owner":     "nobody",
		"file_name": "lookup.csv",
		"file_path": path,
	})

	got, err := lookupTableFileContents(d)
	if err != nil {
		t.Fatal(err)
	}
	if want := `[["id","message"],["1","hello, world"]]`; got != want {
		t.Fatalf("lookupTableFileContents() = %s, want %s", got, want)
	}
}

func TestLookupTableFileReadClearsInactiveState(t *testing.T) {
	t.Setenv("HTTPScheme", "http")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `[["id","value"],["1","hello"]]`)
	}))
	defer server.Close()

	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	client, err := splunkclient.NewSplunkdClient("", [2]string{"admin", "password"}, backend.Host, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("path clears contents", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "lookup.csv")
		if err := os.WriteFile(path, []byte("id,value\n1,hello\n"), 0600); err != nil {
			t.Fatal(err)
		}
		d := schema.TestResourceDataRaw(t, lookupTableFile().Schema, map[string]interface{}{
			"app": "search", "owner": "nobody", "file_name": "lookup.csv", "file_path": path,
		})
		if err := d.Set("file_contents", [][]string{{"old", "state"}}); err != nil {
			t.Fatal(err)
		}
		d.SetId("lookup.csv")

		if err := lookupTableFileRead(d, &SplunkProvider{Client: client}); err != nil {
			t.Fatal(err)
		}
		if got := len(d.Get("file_contents").([]interface{})); got != 0 {
			t.Errorf("file_contents still has %d rows, want it cleared", got)
		}
		if got := d.Get("file_contents_hash").(string); got == "" {
			t.Error("file_contents_hash is empty, want hash of remote contents")
		}
	})

	t.Run("inline contents clears hash", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, lookupTableFile().Schema, map[string]interface{}{
			"app": "search", "owner": "nobody", "file_name": "lookup.csv",
			"file_contents": []interface{}{
				[]interface{}{"id", "value"},
				[]interface{}{"1", "hello"},
			},
		})
		if err := d.Set("file_contents_hash", "stale-hash"); err != nil {
			t.Fatal(err)
		}
		d.SetId("lookup.csv")

		if err := lookupTableFileRead(d, &SplunkProvider{Client: client}); err != nil {
			t.Fatal(err)
		}
		if got := d.Get("file_contents_hash").(string); got != "" {
			t.Errorf("file_contents_hash = %q, want empty", got)
		}
	})
}

func TestAccSplunkLookupTableFile(t *testing.T) {
	// The splunk_lookup_table_file resource uses the lookup_edit API (e.g. from the
	// Splunk App for Lookup File Editing), which is not part of core Splunk. If that
	// API is not available, the provider gets 404. Skip this test unless the env var
	// is set so CI and default runs do not fail.
	if os.Getenv("SPLUNK_TEST_LOOKUP_TABLE_FILE") == "" {
		t.Skip("Skipping lookup table file test; set SPLUNK_TEST_LOOKUP_TABLE_FILE=1 when Splunk has the Lookup File Editing API available (e.g. Splunk App for Lookup File Editing installed)")
	}
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
					resource.TestCheckResourceAttr(resourceName, "file_contents.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.0.0", "status"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.0.1", "status_description"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.0.2", "status_type"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.1.0", "100"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.1.1", "Continue"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.1.2", "Informational"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.2.0", "101"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.2.1", "Switching Protocols"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.2.2", "Informational"),
				),
			},
			{
				Config: updatedLookupTableFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app", "search"),
					resource.TestCheckResourceAttr(resourceName, "owner", "nobody"),
					resource.TestCheckResourceAttr(resourceName, "file_name", "lookup.csv"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.0.0", "status"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.0.1", "status_description"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.0.2", "status_type"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.1.0", "100"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.1.1", "Continue"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.1.2", "Informational"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.2.0", "101"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.2.1", "Switching Protocols"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.2.2", "Informational"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.3.0", "200"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.3.1", "OK"),
					resource.TestCheckResourceAttr(resourceName, "file_contents.3.2", "Successful"),
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
