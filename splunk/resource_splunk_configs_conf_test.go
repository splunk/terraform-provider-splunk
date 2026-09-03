package splunk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	splunkclient "github.com/splunk/terraform-provider-splunk/client"
)

const newConfigsConf = `
resource "splunk_configs_conf" "tftest-stanza" {
	name = "tf_test/tftest_stanza"
	variables = {
        "disabled": "false"
		"key": "value"
	}
}
`

const updateConfigsConf = `
resource "splunk_configs_conf" "tftest-stanza" {
	name = "tf_test/tftest_stanza"
	variables = {
        "disabled": "false"
		"key": "new-value"
	}
}
`

func TestConfigsConfReadUsesConfiguredNamespace(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	var requestPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPaths = append(requestPaths, r.URL.EscapedPath())
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{
			"entry": [{
				"name": "httpevent",
				"acl": {
					"app": "search_app_dashweb",
					"owner": "nobody",
					"sharing": "app",
					"perms": {"read": ["*"], "write": ["admin"]}
				},
				"content": {
					"disabled": false,
					"EXTRACT-environment": "environment=(?P<environment>[^ ]+)"
				}
			}],
			"messages": []
		}`)
	}))
	defer server.Close()

	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	client, err := splunkclient.NewSplunkdClient(
		"",
		[2]string{"admin", "password"},
		backend.Host,
		"",
		server.Client(),
	)
	if err != nil {
		t.Fatal(err)
	}

	d := schema.TestResourceDataRaw(t, configsConf().Schema, map[string]interface{}{
		"name": "props/httpevent",
		"acl": []interface{}{map[string]interface{}{
			"owner":   "nobody",
			"app":     "search_app_dashweb",
			"sharing": "app",
		}},
	})
	d.SetId("props/httpevent")

	if err := configsConfRead(d, &SplunkProvider{Client: client}); err != nil {
		t.Fatalf("configsConfRead: %v", err)
	}

	const expectedPath = "/servicesNS/nobody/search_app_dashweb/configs/conf-props/httpevent"
	if len(requestPaths) == 0 {
		t.Fatal("expected at least one request")
	}
	for _, got := range requestPaths {
		if got != expectedPath {
			t.Errorf("request path = %q, want %q", got, expectedPath)
		}
	}

	acl := getACLConfig(d.Get("acl").([]interface{}))
	if got, want := acl.App, "search_app_dashweb"; got != want {
		t.Errorf("acl app = %q, want %q", got, want)
	}
	if got, want := acl.Owner, "nobody"; got != want {
		t.Errorf("acl owner = %q, want %q", got, want)
	}
}

func TestGetResourceDataConfigsConfACL(t *testing.T) {
	tests := []struct {
		name        string
		raw         map[string]interface{}
		wantOwner   string
		wantApp     string
		wantSharing string
	}{
		{
			name:        "default namespace",
			wantOwner:   "nobody",
			wantApp:     "search",
			wantSharing: "app",
		},
		{
			name: "app shared namespace uses nobody",
			raw: map[string]interface{}{
				"acl": []interface{}{map[string]interface{}{
					"owner":   "admin",
					"app":     "custom_app",
					"sharing": "app",
				}},
			},
			wantOwner:   "nobody",
			wantApp:     "custom_app",
			wantSharing: "app",
		},
		{
			name: "user shared namespace keeps owner",
			raw: map[string]interface{}{
				"acl": []interface{}{map[string]interface{}{
					"owner":   "alice",
					"app":     "custom_app",
					"sharing": "user",
				}},
			},
			wantOwner:   "alice",
			wantApp:     "custom_app",
			wantSharing: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, configsConf().Schema, tt.raw)
			got := getResourceDataConfigsConfACL(d)
			if got.Owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", got.Owner, tt.wantOwner)
			}
			if got.App != tt.wantApp {
				t.Errorf("app = %q, want %q", got.App, tt.wantApp)
			}
			if got.Sharing != tt.wantSharing {
				t.Errorf("sharing = %q, want %q", got.Sharing, tt.wantSharing)
			}
		})
	}
}

func TestAccCreateSplunkConfigsConf(t *testing.T) {
	resourceName := "splunk_configs_conf.tftest-stanza"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkConfigsConfDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newConfigsConf,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.key", "value"),
				),
			},
			{
				Config: updateConfigsConf,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.key", "new-value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const newConfigsConfSpecialChars = `
resource "splunk_configs_conf" "tftest-stanza-special-chars" {
	name = "tf_test/sqs://tftest_stanza_special_chars"
	variables = {
        "disabled": "false"
		"key": "value"
	}
}
`

const updateConfigsConfSpecialChars = `
resource "splunk_configs_conf" "tftest-stanza-special-chars" {
	name = "tf_test/sqs://tftest_stanza_special_chars"
	variables = {
        "disabled": "false"
		"key": "new-value"
	}
}
`

func TestAccCreateSplunkConfigsConfSpecialChars(t *testing.T) {
	resourceName := "splunk_configs_conf.tftest-stanza-special-chars"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkConfigsConfDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: newConfigsConfSpecialChars,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.key", "value"),
				),
			},
			{
				Config: updateConfigsConfSpecialChars,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.key", "new-value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkConfigsConfDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_configs_conf":
			endpoint := client.BuildSplunkURL(nil, "services", "configs", "conf", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
