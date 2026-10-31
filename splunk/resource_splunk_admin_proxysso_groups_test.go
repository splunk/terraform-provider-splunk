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

const adminProxyssoGroupsInput = `
resource "splunk_admin_proxysso_groups" "proxy_sso_group" {
  name = "new-proxysso-group"
  roles = ["admin", "power"]
}
`

const updateAdminProxyssoGroupsInput = `
resource "splunk_admin_proxysso_groups" "proxy_sso_group" {
  name = "new-proxysso-group"
  roles = ["admin", "power", "user"]
}
`

func TestAdminProxyssoGroupsSchema(t *testing.T) {
	res := adminProxyssoGroups()
	name := res.Schema["name"]
	if !name.Required || !name.ForceNew || name.Type != schema.TypeString {
		t.Fatalf("name schema = required %v forceNew %v type %v", name.Required, name.ForceNew, name.Type)
	}
	roles := res.Schema["roles"]
	if !roles.Required || roles.Optional || roles.Computed || roles.Type != schema.TypeSet {
		t.Fatalf("roles schema = required %v optional %v computed %v type %v", roles.Required, roles.Optional, roles.Computed, roles.Type)
	}
	if roles.Set == nil || roles.Set("admin") != schema.HashString("admin") {
		t.Fatal("roles must use schema.HashString so acceptance checks find roles.<hash>")
	}

	p := Provider().(*schema.Provider)
	if _, ok := p.ResourcesMap["splunk_admin_proxysso_groups"]; !ok {
		t.Fatal("splunk_admin_proxysso_groups is not registered")
	}
}

func TestAdminProxyssoGroupsCreateReadUpdateDelete(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	const groupName = "A/B+C D"
	const entityPath = "/services/admin/ProxySSO-groups/A%2FB%2BC%20D"
	const collectionPath = "/services/admin/ProxySSO-groups"

	roles := []string{"power", "admin"}
	readBody := func() string {
		return `{"entry":[{"name":"A/B+C D","content":{"roles":["` + strings.Join(roles, `","`) + `"]}}]}`
	}

	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests = append(requests, r.Method+" "+r.URL.EscapedPath()+" "+string(body))
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.EscapedPath() == collectionPath:
			values, _ := url.ParseQuery(string(body))
			if values.Get("name") != groupName {
				t.Errorf("create name = %q", values.Get("name"))
			}
			if !sameRoles(values["roles"], []string{"admin", "power"}) {
				t.Errorf("create roles = %#v", values["roles"])
			}
			roles = append([]string(nil), values["roles"]...)
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, readBody())
		case r.Method == http.MethodGet && r.URL.EscapedPath() == entityPath:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, readBody())
		case r.Method == http.MethodPost && r.URL.EscapedPath() == entityPath:
			values, _ := url.ParseQuery(string(body))
			if values.Get("name") != "" {
				t.Errorf("update included name %q", values.Get("name"))
			}
			if !sameRoles(values["roles"], []string{"user", "admin"}) {
				t.Errorf("update roles = %#v", values["roles"])
			}
			roles = append([]string(nil), values["roles"]...)
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, readBody())
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == entityPath:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{}`)
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.EscapedPath(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	provider := proxyssoTestProvider(t, server)
	d := schema.TestResourceDataRaw(t, adminProxyssoGroups().Schema, map[string]interface{}{
		"name":  groupName,
		"roles": []interface{}{"admin", "power"},
	})

	if err := adminProxyssoGroupsCreate(d, provider); err != nil {
		t.Fatalf("create: %v", err)
	}
	if d.Id() != groupName {
		t.Fatalf("id = %q, want %q", d.Id(), groupName)
	}
	assertRoles(t, d, "admin", "power")

	if err := d.Set("roles", []interface{}{"user", "admin"}); err != nil {
		t.Fatal(err)
	}
	if err := adminProxyssoGroupsUpdate(d, provider); err != nil {
		t.Fatalf("update: %v", err)
	}
	assertRoles(t, d, "user", "admin")

	if err := adminProxyssoGroupsDelete(d, provider); err != nil {
		t.Fatalf("delete: %v", err)
	}

	for _, req := range requests {
		if strings.Contains(req, "%252F") || strings.Contains(req, "%2520") {
			t.Errorf("name was encoded twice: %s", req)
		}
	}
	wantPrefixes := []string{
		"POST " + collectionPath + " ",
		"GET " + entityPath + " ",
		"POST " + entityPath + " ",
		"GET " + entityPath + " ",
		"DELETE " + entityPath + " ",
	}
	if len(requests) != len(wantPrefixes) {
		t.Fatalf("requests =\n%s", strings.Join(requests, "\n"))
	}
	for i, prefix := range wantPrefixes {
		if !strings.HasPrefix(requests[i], prefix) {
			t.Errorf("request %d = %q, want prefix %q", i, requests[i], prefix)
		}
	}
}

func TestAdminProxyssoGroupsReadDropsMissingGroup(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	tests := []struct {
		name   string
		status int
		body   string
	}{
		{
			name:   "400 role mapping message",
			status: http.StatusBadRequest,
			body:   `{"messages":[{"type":"ERROR","text":"Unable to find a role mapping for the given external role."}]}`,
		},
		{
			name:   "404",
			status: http.StatusNotFound,
			body:   `{"messages":[{"type":"ERROR","text":"Could not find object id=mygroup"}]}`,
		},
		{
			name:   "200 without the group",
			status: http.StatusOK,
			body:   `{"entry":[{"name":"other","content":{"roles":["user"]}}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer server.Close()

			d := schema.TestResourceDataRaw(t, adminProxyssoGroups().Schema, map[string]interface{}{
				"name":  "mygroup",
				"roles": []interface{}{"admin"},
			})
			d.SetId("mygroup")

			if err := adminProxyssoGroupsRead(d, proxyssoTestProvider(t, server)); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "" {
				t.Fatalf("id = %q, want empty so Terraform can recreate the group", d.Id())
			}
		})
	}
}

func TestAdminProxyssoGroupsReadKeepsStateOnOtherErrors(t *testing.T) {
	t.Setenv("HTTPScheme", "http")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"messages":[{"type":"ERROR","text":"boom"}]}`)
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, adminProxyssoGroups().Schema, map[string]interface{}{
		"name":  "mygroup",
		"roles": []interface{}{"admin"},
	})
	d.SetId("mygroup")

	err := adminProxyssoGroupsRead(d, proxyssoTestProvider(t, server))
	if err == nil {
		t.Fatal("expected read error")
	}
	if d.Id() != "mygroup" {
		t.Fatalf("id = %q, want mygroup", d.Id())
	}
}

func TestAccSplunkAdminProxyssoGroups(t *testing.T) {
	resourceName := "splunk_admin_proxysso_groups.proxy_sso_group"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkAdminProxyssoGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: adminProxyssoGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-proxysso-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					testCheckProxyssoRole(resourceName, "admin"),
					testCheckProxyssoRole(resourceName, "power"),
				),
			},
			{
				Config: updateAdminProxyssoGroupsInput,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "new-proxysso-group"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "3"),
					testCheckProxyssoRole(resourceName, "admin"),
					testCheckProxyssoRole(resourceName, "power"),
					testCheckProxyssoRole(resourceName, "user"),
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

func testAccSplunkAdminProxyssoGroupsDestroy(s *terraform.State) error {
	api, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "splunk_admin_proxysso_groups" {
			continue
		}
		endpoint := api.BuildSplunkURLWithEscapedPathPart(nil, rs.Primary.ID, "services", "admin", "ProxySSO-groups")
		resp, err := api.Get(endpoint)
		if resp != nil {
			resp.Body.Close()
		}
		if adminProxyssoGroupMissing(err) {
			continue
		}
		return fmt.Errorf("proxysso group %s still exists: %v", rs.Primary.ID, err)
	}
	return nil
}

func testCheckProxyssoRole(resourceName, role string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(
		resourceName,
		fmt.Sprintf("roles.%d", schema.HashString(role)),
		role,
	)
}

func proxyssoTestProvider(t *testing.T, server *httptest.Server) *SplunkProvider {
	t.Helper()
	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	api, err := client.NewSplunkdClient("", [2]string{}, backend.Host, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	return &SplunkProvider{Client: api}
}

func assertRoles(t *testing.T, d *schema.ResourceData, want ...string) {
	t.Helper()
	got := d.Get("roles").(*schema.Set)
	if got.Len() != len(want) {
		t.Fatalf("roles = %#v, want %v", got.List(), want)
	}
	for _, role := range want {
		if !got.Contains(role) {
			t.Fatalf("roles = %#v, missing %s", got.List(), role)
		}
	}
}

func sameRoles(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[string]int{}
	for _, role := range got {
		seen[role]++
	}
	for _, role := range want {
		seen[role]--
		if seen[role] < 0 {
			return false
		}
	}
	return true
}
