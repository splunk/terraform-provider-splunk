package splunk

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Live checks for object names that contain a literal "+".
// splunkd treats a bare "+" in a path segment as a space, so create, update,
// refresh, import, and delete have to send that character as %2B.

func TestAccSplunkSavedSearchPlusName(t *testing.T) {
	const (
		resourceName = "splunk_saved_searches.plus"
		name         = "Test 30+ days"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c, err := newTestClient()
					if err != nil {
						t.Fatal(err)
					}
					for _, owner := range []string{"nobody", "admin"} {
						resp, _ := c.DeleteSavedSearches(name, owner, "search")
						if resp != nil {
							resp.Body.Close()
						}
					}
				},
				Config: savedSearchPlusConfig("index=main"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "search", "index=main"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "app"),
				),
			},
			{
				Config: savedSearchPlusConfig("index=_internal"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "search", "index=_internal"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "/servicesNS/nobody/search/saved/searches/Test%2030%2B%20days",
				ImportStateVerify: true,
			},
		},
	})
}

func savedSearchPlusConfig(search string) string {
	return fmt.Sprintf(`
resource "splunk_saved_searches" "plus" {
  name   = "Test 30+ days"
  search = %q
  acl {
    owner   = "admin"
    app     = "search"
    sharing = "app"
  }
}
`, search)
}

func TestAccSplunkDashboardPlusName(t *testing.T) {
	// data/ui/views names may contain only alphanumeric ASCII, "-", ".", and "_".
	// A literal "+" is rejected before the path encoding matters.
	t.Skip("splunkd rejects '+' in dashboard view names")

	const (
		resourceName = "splunk_data_ui_views.plus"
		name         = "Test 30+ days"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c, err := newTestClient()
					if err != nil {
						t.Fatal(err)
					}
					for _, owner := range []string{"admin", "nobody"} {
						resp, _ := c.DeleteDashboardObject(owner, "search", name)
						if resp != nil {
							resp.Body.Close()
						}
					}
				},
				Config: dashboardPlusConfig("<dashboard><label>Test 30+ days</label></dashboard>"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Test 30+ days</label></dashboard>"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.owner", "admin"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.app", "search"),
					resource.TestCheckResourceAttr(resourceName, "acl.0.sharing", "app"),
				),
			},
			{
				Config: dashboardPlusConfig("<dashboard><label>Test 30+ days updated</label></dashboard>"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "eai_data", "<dashboard><label>Test 30+ days updated</label></dashboard>"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "/servicesNS/nobody/search/data/ui/views/Test%2030%2B%20days",
				ImportStateVerify: true,
			},
		},
	})
}

func dashboardPlusConfig(eaiData string) string {
	return fmt.Sprintf(`
resource "splunk_data_ui_views" "plus" {
  name     = "Test 30+ days"
  eai_data = %q
  acl {
    owner   = "admin"
    app     = "search"
    sharing = "app"
  }
}
`, eaiData)
}

func TestAccSplunkAdminSAMLGroupsPlusName(t *testing.T) {
	const (
		resourceName = "splunk_admin_saml_groups.plus"
		name         = "tf-acc+saml"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c, err := newTestClient()
					if err != nil {
						t.Fatal(err)
					}
					resp, _ := c.DeleteAdminSAMLGroups(name)
					if resp != nil {
						resp.Body.Close()
					}
				},
				Config: samlPlusConfig(`["admin", "power"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
				),
			},
			{
				Config: samlPlusConfig(`["admin", "power", "user"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "3"),
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

func samlPlusConfig(roles string) string {
	return fmt.Sprintf(`
resource "splunk_admin_saml_groups" "plus" {
  name  = "tf-acc+saml"
  roles = %s
}
`, roles)
}

func TestAccSplunkAdminProxyssoGroupsPlusName(t *testing.T) {
	const (
		resourceName = "splunk_admin_proxysso_groups.plus"
		name         = "tf-acc+proxysso"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c, err := newTestClient()
					if err != nil {
						t.Fatal(err)
					}
					resp, _ := c.DeleteAdminProxyssoGroups(name)
					if resp != nil {
						resp.Body.Close()
					}
				},
				Config: proxyssoPlusConfig(`["admin", "power"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
				),
			},
			{
				Config: proxyssoPlusConfig(`["admin", "power", "user"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "3"),
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

func proxyssoPlusConfig(roles string) string {
	return fmt.Sprintf(`
resource "splunk_admin_proxysso_groups" "plus" {
  name  = "tf-acc+proxysso"
  roles = %s
}
`, roles)
}

func TestAccSplunkTCPTokenPlusName(t *testing.T) {
	const (
		resourceName = "splunk_inputs_tcp_splunk_tcp_token.plus"
		name         = "tf-acc+token"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c, err := newTestClient()
					if err != nil {
						t.Fatal(err)
					}
					resp, _ := c.DeleteSplunkTCPTokenInput(name, "nobody", "search")
					if resp != nil {
						resp.Body.Close()
					}
				},
				Config: tcpTokenPlusConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "token", "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee"),
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

const tcpTokenPlusConfig = `
resource "splunk_inputs_tcp_splunk_tcp_token" "plus" {
  name  = "tf-acc+token"
  token = "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee"
}
`

func TestAccSplunkInputsMonitorPlusName(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance tests skipped unless TF_ACC is set")
	}
	path := monitorPlusPath(t)

	resourceName := "splunk_inputs_monitor.plus"
	config := fmt.Sprintf(`
resource "splunk_inputs_monitor" "plus" {
  name       = %q
  sourcetype = "text"
}
`, path)
	updated := fmt.Sprintf(`
resource "splunk_inputs_monitor" "plus" {
  name       = %q
  sourcetype = "text"
  recursive  = true
}
`, path)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", path),
					resource.TestCheckResourceAttr(resourceName, "sourcetype", "text"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", path),
					resource.TestCheckResourceAttr(resourceName, "recursive", "true"),
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

// monitorPlusPath returns a file path Splunk can monitor. A local Splunk
// install can read a file written under SPLUNK_HOME. GitHub Actions runs the
// tests on the host and Splunk in a service container, so SPLUNK_HOME is not
// on the host. That workflow mounts the workspace at /workspace inside the
// container, which is where the file has to live.
func monitorPlusPath(t *testing.T) string {
	t.Helper()
	const fileName = "tfacc-plus+name.log"
	home := os.Getenv("SPLUNK_HOME")
	if home == "" {
		t.Fatal("SPLUNK_HOME must be set for acceptance tests")
	}
	dir := filepath.Join(home, "var", "log", "splunk")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		path := filepath.Join(dir, fileName)
		writeMonitorPlusFile(t, path)
		return path
	}
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		t.Skipf("monitor directory %s is not available on this host", dir)
	}
	writeMonitorPlusFile(t, filepath.Join(workspace, fileName))
	return "/workspace/" + fileName
}

func writeMonitorPlusFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("tfacc\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })
}

func TestAccSplunkInputsScriptPlusName(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance tests skipped unless TF_ACC is set")
	}
	home := os.Getenv("SPLUNK_HOME")
	if home == "" {
		t.Fatal("SPLUNK_HOME must be set for acceptance tests")
	}
	dir := filepath.Join(home, "bin", "scripts")
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("scripted input directory %s is not available: %s", dir, err)
	}
	path := filepath.Join(dir, "tfacc+plus.py")
	if err := os.WriteFile(path, []byte("#!/usr/bin/env python3\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })

	resourceName := "splunk_inputs_script.plus"
	config := fmt.Sprintf(`
resource "splunk_inputs_script" "plus" {
  name     = %q
  interval = 60
}
`, path)
	updated := fmt.Sprintf(`
resource "splunk_inputs_script" "plus" {
  name     = %q
  interval = 120
}
`, path)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPlusNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", path),
					resource.TestCheckResourceAttr(resourceName, "interval", "60"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", path),
					resource.TestCheckResourceAttr(resourceName, "interval", "120"),
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

func testAccPlusNameDestroy(s *terraform.State) error {
	c, err := newTestClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		var resp *http.Response
		switch rs.Type {
		case "splunk_saved_searches":
			resp, err = c.ReadSavedSearches(rs.Primary.ID, "nobody", "search")
		case "splunk_data_ui_views":
			resp, err = c.ReadDashboardObject(rs.Primary.ID, "nobody", "search")
		case "splunk_admin_saml_groups":
			resp, err = c.ReadAdminSAMLGroups(rs.Primary.ID)
		case "splunk_admin_proxysso_groups":
			resp, err = c.ReadAdminProxyssoGroups(rs.Primary.ID)
		case "splunk_inputs_tcp_splunk_tcp_token":
			resp, err = c.ReadSplunkTCPTokenInput(rs.Primary.ID, "nobody", "search")
		case "splunk_inputs_monitor":
			resp, err = c.ReadMonitorInput(rs.Primary.ID, "nobody", "search")
		case "splunk_inputs_script":
			resp, err = c.ReadScriptedInput(rs.Primary.ID, "nobody", "search")
		default:
			continue
		}
		if resp != nil {
			resp.Body.Close()
		}
		if plusNameGone(rs.Type, resp, err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("%s %s: %w", rs.Type, rs.Primary.ID, err)
		}
		return fmt.Errorf("%s %q still exists after destroy", rs.Type, rs.Primary.ID)
	}
	return nil
}

// plusNameGone reports a read after delete. Input readers return a 404 as an
// error and a nil response. SAML and ProxySSO groups answer a missing mapping
// with 400 "Unable to find a role mapping".
func plusNameGone(resourceType string, resp *http.Response, err error) bool {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "404 ") {
		return true
	}
	switch resourceType {
	case "splunk_admin_saml_groups", "splunk_admin_proxysso_groups":
		return strings.Contains(msg, "Unable to find a role mapping")
	default:
		return false
	}
}
