package splunk

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const splunkDashboardsObjectWithAcl = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
  acl {
    owner = "admin"
    app = "search"
    sharing = "global"
    read = ["*"]
    write = ["admin"]
  }
}
`

const updateSplunkDashboardsObjectWithAcl = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
}
`

const splunkDashboardsObject = `
resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
}
`

const updateSplunkDashboardsObject = `

resource "splunk_data_ui_views" "dashboard" {
  name     = "terraform"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
}
`

func TestAccSplunkDashboards(t *testing.T) {
	resourceName := "splunk_data_ui_views.dashboard"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkDashboardDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: splunkDashboardsObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
				),
			},
			{
				Config: updateSplunkDashboardsObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "user"),
				),
			},
			{
				ResourceName:      "splunk_data_ui_views.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSplunkDashboardsWithAcl(t *testing.T) {
	resourceName := "splunk_data_ui_views.dashboard"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSplunkDashboardDestroyResources,
		Steps: []resource.TestStep{
			{
				Config: splunkDashboardsObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				Config: updateSplunkDashboardsObjectWithAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Terraform Test Dashboard</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "global"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.read.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.write.0", "admin"),
				),
			},
			{
				ResourceName:      "splunk_data_ui_views.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSplunkDashboardDestroyResources(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "splunk_data_ui_views":
			endpoint := client.BuildSplunkURL(nil, "servicesNS", "admin", "search", "data", "ui", "views", rs.Primary.ID)
			resp, err := client.Get(endpoint)
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error: %s: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}

const dashboardEntry = `{"entry":[{"name":"terraform","acl":{"app":"search","owner":"admin","sharing":"app"},"content":{"eai:data":"<dashboard/>"}}]}`

func TestSplunkDashboardsCreateWaitsForReplication(t *testing.T) {
	var creates, reads int
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/servicesNS/admin/search/data/ui/views":
			creates++
			if creates == 1 {
				return jsonResponse(http.StatusConflict, `{"messages":[{"type":"ERROR","text":"An object with name=terraform already exists"}]}`), nil
			}
			return jsonResponse(http.StatusCreated, "{}"), nil
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/acl"):
			return jsonResponse(http.StatusOK, "{}"), nil
		case r.Method == http.MethodGet && r.URL.Path == "/servicesNS/nobody/search/data/ui/views/terraform":
			reads++
			if reads == 1 {
				return notFoundResponse(), nil
			}
			return jsonResponse(http.StatusOK, dashboardEntry), nil
		}
		t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		return nil, nil
	})

	dashboard := splunkDashboards()
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":     "terraform",
		"eai_data": "<dashboard/>",
		"acl":      []interface{}{map[string]interface{}{"owner": "admin", "app": "search", "sharing": "app"}},
	})
	diff, err := dashboard.Diff(nil, config, meta)
	if err != nil {
		t.Fatal(err)
	}
	state, err := dashboard.Apply(nil, diff, meta)
	if err != nil {
		t.Fatal(err)
	}
	if state.ID != "terraform" || state.Attributes["eai_data"] != "<dashboard/>" {
		t.Errorf("state = %v", state)
	}
	if creates != 2 || reads != 2 {
		t.Errorf("creates = %d, reads = %d, want 2 and 2", creates, reads)
	}
}

func TestSplunkDashboardsReadRemovesDeletedDashboard(t *testing.T) {
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		return notFoundResponse(), nil
	})

	d := schema.TestResourceDataRaw(t, splunkDashboards().Schema, map[string]interface{}{
		"name":     "terraform",
		"eai_data": "<dashboard/>",
		"acl":      []interface{}{map[string]interface{}{"owner": "admin", "app": "search", "sharing": "app"}},
	})
	d.SetId("terraform")

	if err := splunkDashboardsRead(d, meta); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "" {
		t.Errorf("id = %q, want the resource removed from state", d.Id())
	}
}

func TestSplunkDashboardsDeleteIgnoresDeletedDashboard(t *testing.T) {
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		return notFoundResponse(), nil
	})

	d := schema.TestResourceDataRaw(t, splunkDashboards().Schema, map[string]interface{}{
		"name":     "terraform",
		"eai_data": "<dashboard/>",
		"acl":      []interface{}{map[string]interface{}{"owner": "admin", "app": "search", "sharing": "app"}},
	})
	d.SetId("terraform")

	if err := splunkDashboardsDelete(d, meta); err != nil {
		t.Fatal(err)
	}
}

func TestSplunkDashboardsTimeoutsKeepExistingStateUnchanged(t *testing.T) {
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodGet {
			return jsonResponse(http.StatusOK, dashboardEntry), nil
		}
		return jsonResponse(http.StatusCreated, "{}"), nil
	})
	dashboard := splunkDashboards()
	if got := *dashboard.Timeouts.Create; got != defaultReplicationTimeout {
		t.Errorf("default create timeout = %s, want %s", got, defaultReplicationTimeout)
	}

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":     "terraform",
		"eai_data": "<dashboard/>",
		"acl":      []interface{}{map[string]interface{}{"owner": "admin", "app": "search", "sharing": "app"}},
	})
	diff, err := dashboard.Diff(nil, config, meta)
	if err != nil {
		t.Fatal(err)
	}
	state, err := dashboard.Apply(nil, diff, meta)
	if err != nil {
		t.Fatal(err)
	}
	// state from a version without timeouts
	delete(state.Meta, schema.TimeoutKey)

	diff, err = dashboard.Diff(state, config, meta)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil && !diff.Empty() {
		t.Errorf("diff = %v, want no changes", diff.Attributes)
	}
}

func TestSplunkDashboardsUpdateSendsACLOnlyWhenChanged(t *testing.T) {
	eaiData, write := "<dashboard/>", "admin"
	var aclRequests int
	meta := newRoundTripProvider(t, func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			form, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasSuffix(r.URL.Path, "/acl") {
				aclRequests++
				write = form.Get("perms.write")
			} else if form.Has("eai:data") {
				eaiData = form.Get("eai:data")
			}
			return jsonResponse(http.StatusOK, "{}"), nil
		}
		return jsonResponse(http.StatusOK, fmt.Sprintf(`{"entry":[{"name":"terraform","acl":{"app":"search","owner":"jdoe","sharing":"app","perms":{"read":["*"],"write":[%q]}},"content":{"eai:data":%q}}]}`, write, eaiData)), nil
	})

	dashboard := splunkDashboards()
	apply := func(state *terraform.InstanceState, eaiData, write string) *terraform.InstanceState {
		t.Helper()
		config := terraform.NewResourceConfigRaw(map[string]interface{}{
			"name":     "terraform",
			"eai_data": eaiData,
			"acl":      []interface{}{map[string]interface{}{"owner": "jdoe", "app": "search", "sharing": "app", "read": []interface{}{"*"}, "write": []interface{}{write}}},
		})
		diff, err := dashboard.Diff(state, config, meta)
		if err != nil {
			t.Fatal(err)
		}
		state, err = dashboard.Apply(state, diff, meta)
		if err != nil {
			t.Fatal(err)
		}
		return state
	}

	state := apply(nil, "<dashboard/>", "admin")
	aclRequests = 0

	state = apply(state, "<form/>", "admin")
	if aclRequests != 0 {
		t.Errorf("acl requests = %d after an eai_data change, want 0", aclRequests)
	}

	apply(state, "<form/>", "power")
	if aclRequests != 1 {
		t.Errorf("acl requests = %d after an acl change, want 1", aclRequests)
	}
}
