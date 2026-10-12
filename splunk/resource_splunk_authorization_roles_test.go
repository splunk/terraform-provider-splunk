package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/splunk/terraform-provider-splunk/client"
)

const authorizationRolesInput = `
resource "splunk_authorization_roles" "role" {
  name = "new-role"
  default_app = "search"
  imported_roles = ["user"]
  # capabilities are explicitly listed in a different order than returned by the Splunk API
  capabilities = ["change_authentication", "accelerate_datamodel"]
  search_indexes_allowed = ["_audit", "_internal"]
  search_indexes_default = ["_audit", "_internal"]
}
`

const updateAuthorizationRolesInput = `
resource "splunk_authorization_roles" "role" {
  name = "new-role"
  default_app = "search"
  imported_roles = ["power", "user"]
  # capabilities are explicitly listed in a different order than returned by the Splunk API
  capabilities = ["restart_splunkd", "change_authentication", "accelerate_datamodel"]
  search_indexes_allowed = ["_audit", "_internal", "main"]
  search_indexes_default = ["_audit", "_internal", "main"]
}
`

func TestAuthorizationRoleNameIsLowercased(t *testing.T) {
	stateFunc := authorizationRoles().Schema["name"].StateFunc
	if stateFunc == nil {
		t.Fatal("name StateFunc must lowercase role names before they are stored")
	}
	if got := stateFunc("MyRole"); got != "myrole" {
		t.Fatalf("StateFunc(%q) = %q, want myrole", "MyRole", got)
	}
	if got := canonicalAuthorizationRoleName("MyRole"); got != "myrole" {
		t.Fatalf("canonicalAuthorizationRoleName(%q) = %q, want myrole", "MyRole", got)
	}
}

func TestGetAuthorizationRolesByNameIgnoresCase(t *testing.T) {
	body := `{"entry":[{"name":"other"},{"name":"myrole","acl":{"owner":"nobody","app":"system"}}]}`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}

	entry, err := getAuthorizationRolesByName("MyRole", resp)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("expected to find myrole when looking up MyRole")
	}
	if entry.Name != "myrole" {
		t.Fatalf("entry.Name = %q, want myrole", entry.Name)
	}
	if entry.ACL.Owner != "nobody" || entry.ACL.App != "system" {
		t.Fatalf("acl = owner %q app %q, want nobody/system", entry.ACL.Owner, entry.ACL.App)
	}
}

func TestAuthorizationRolesReadRewritesMixedCaseID(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	const listPath = "/services/authorization/roles"
	const rolePath = "/services/authorization/roles/myrole"
	const roleBody = `{
		"entry": [{
			"name": "myrole",
			"acl": {"owner": "nobody", "app": "system"},
			"content": {
				"capabilities": ["change_authentication"],
				"defaultApp": "search",
				"imported_roles": ["user"]
			}
		}],
		"messages": []
	}`

	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.EscapedPath())
		if strings.Contains(r.URL.EscapedPath(), "MyRole") {
			http.Error(w, "mixed-case role path", http.StatusNotFound)
			return
		}
		switch {
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == rolePath:
			w.WriteHeader(http.StatusOK)
		case r.URL.EscapedPath() == listPath || r.URL.EscapedPath() == rolePath:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, roleBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	// The local test server does not check credentials.
	api, err := client.NewSplunkdClient("", [2]string{}, backend.Host, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	provider := &SplunkProvider{Client: api}

	d := schema.TestResourceDataRaw(t, authorizationRoles().Schema, map[string]interface{}{
		"name": "MyRole",
	})
	d.SetId("MyRole")

	if err := authorizationRolesRead(d, provider); err != nil {
		t.Fatalf("authorizationRolesRead: %v", err)
	}
	if got := d.Id(); got != "myrole" {
		t.Fatalf("id = %q, want myrole", got)
	}
	if got := d.Get("name").(string); got != "myrole" {
		t.Fatalf("name = %q, want myrole", got)
	}
	if got := d.Get("default_app").(string); got != "search" {
		t.Fatalf("default_app = %q, want search", got)
	}
	imported := d.Get("imported_roles").([]interface{})
	if len(imported) != 1 || imported[0] != "user" {
		t.Fatalf("imported_roles = %#v, want [user]", imported)
	}

	if err := authorizationRolesDelete(d, provider); err != nil {
		t.Fatalf("authorizationRolesDelete: %v", err)
	}

	want := []string{
		"GET " + listPath,
		"GET " + rolePath,
		"DELETE " + rolePath,
	}
	if strings.Join(requests, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests =\n%s\nwant =\n%s", strings.Join(requests, "\n"), strings.Join(want, "\n"))
	}
}

func TestAccSplunkAuthorizationRoles(t *testing.T) {
	resourceName := "splunk_authorization_roles.role"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAuthorizationRolesInputDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: authorizationRolesInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-role"),
					resource.TestCheckResourceAttr(resourceName, "default_app", "search"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.0", "user"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.1", "_internal"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.1", "_internal"),
				),
			},
			{
				Config: updateAuthorizationRolesInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-role"),
					resource.TestCheckResourceAttr(resourceName, "default_app", "search"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.0", "power"),
					resource.TestCheckResourceAttr(resourceName, "imported_roles.1", "user"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.1", "_internal"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_allowed.2", "main"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.0", "_audit"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.1", "_internal"),
					resource.TestCheckResourceAttr(resourceName, "search_indexes_default.2", "main"),
				),
			},
			{
				ResourceName:      "splunk_authorization_roles.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkAuthorizationRolesInputDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_authorization_roles":
			endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}
