package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const newUDPInput = `
resource "splunk_inputs_udp" "test" {
    name = "41000"
    source = "new"
    disabled = false
}
`

const updateUDPInput = `

resource "splunk_inputs_udp" "test" {
    name = "41000"
    index = "main"
    source = "new"
    sourcetype = "new"
    disabled = false
}
`

func TestAccSplunkUDPInput(t *testing.T) {
	resourceName := "splunk_inputs_udp.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkUDPInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newUDPInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "default"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "ip"),
					resource.TestCheckResourceAttr(resourceName, "queue", ""),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", ""),
					resource.TestCheckResourceAttr(resourceName, "no_appending_timestamp", "false"),
					resource.TestCheckResourceAttr(resourceName, "no_priority_stripping", "false"),
				),
			},
			{
				Config: updateUDPInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source", "new"),
					resource.TestCheckResourceAttr(resourceName, "index", "main"),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "new"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "queue", ""),
					resource.TestCheckResourceAttr(resourceName, "connection_host", "ip"),
					resource.TestCheckResourceAttr(resourceName, "restrict_to_host", ""),
					resource.TestCheckResourceAttr(resourceName, "no_priority_stripping", "false"),
				),
			},
			{
				ResourceName:      "splunk_inputs_udp.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkUDPInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_inputs_udp":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "search", "data", "inputs", "udp", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}

func TestInputsUDPDeleteReturnsClientError(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	tests := []struct {
		name string
		body string
	}{
		{
			name: "splunk message",
			body: `{"messages":[{"type":"ERROR","text":"Could not find object id=41000"}]}`,
		},
		{
			name: "empty messages",
			body: `{"messages":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.Method + " " + r.URL.EscapedPath()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer server.Close()

			d := udpDeleteResourceData(t)
			err := inputsUDPDelete(d, httpTestProvider(t, server))
			if err == nil {
				t.Fatal("inputsUDPDelete error = nil, want the client status error")
			}
			if !strings.Contains(err.Error(), "400 Bad Request") || !strings.Contains(err.Error(), tt.body) {
				t.Fatalf("inputsUDPDelete error = %q, want status and response body", err)
			}
			const want = "DELETE /servicesNS/nobody/search/data/inputs/udp/41000"
			if got != want {
				t.Fatalf("request = %q, want %q", got, want)
			}
		})
	}
}

func TestInputsUDPDeleteSucceeds(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	if err := inputsUDPDelete(udpDeleteResourceData(t), httpTestProvider(t, server)); err != nil {
		t.Fatal(err)
	}
}

func udpDeleteResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, inputsUDP().Schema, map[string]interface{}{
		"name": "41000",
		"acl": []interface{}{map[string]interface{}{
			"owner":   "nobody",
			"app":     "search",
			"sharing": "app",
		}},
	})
	d.SetId("41000")
	return d
}
