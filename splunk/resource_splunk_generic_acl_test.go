package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const newGenericAcl = `
// simulate an externally (to terraform) app deployment, so not defining ACL for my_app or my_view
resource "splunk_apps_local" "my_app" {
  name = "my_app"
}

resource "splunk_data_ui_views" "my_view" {
  name = "my_view"
  acl {
    owner = "admin"
    app   = "my_app"
  }
  eai_data = <<EOF
<dashboard>
 <label>my_dashboard</label>
</dashboard>
EOF

  depends_on = [splunk_apps_local.my_app]
}

resource "splunk_generic_acl" "my_app" {
  path = "apps/local/my_app"
  acl {
    owner = "nobody"
    app   = "system"
    read  = ["admin", "power"]
    write = ["admin", "power"]
  }

  depends_on = [splunk_apps_local.my_app]
}

resource "splunk_generic_acl" "my_view" {
  path = "data/ui/views/my_view"
  acl {
    owner = "admin"
    app   = "my_app"
    read  = ["admin", "power"]
    write = ["admin", "power"]
  }

  depends_on = [splunk_data_ui_views.my_view]
}
`

const updateGenericAcl = `
// simulate an externally (to terraform) app deployment, so not defining ACL for my_app or my_view
resource "splunk_apps_local" "my_app" {
  name = "my_app"
}

resource "splunk_data_ui_views" "my_view" {
  name = "my_view"
  acl {
    owner = "admin"
    app   = "my_app"
  }
  eai_data = <<EOF
<dashboard>
 <label>my_dashboard</label>
</dashboard>
EOF

  depends_on = [splunk_apps_local.my_app]
}

resource "splunk_generic_acl" "my_app" {
  path = "apps/local/my_app"
  acl {
    owner = "nobody"
    app   = "system"
    read  = ["admin"]
    write = ["admin"]
  }

  depends_on = [splunk_apps_local.my_app]
}

resource "splunk_generic_acl" "my_view" {
  path = "data/ui/views/my_view"
  acl {
    owner = "admin"
    app   = "my_app"
    read  = ["admin"]
    write = ["admin"]
  }

  depends_on = [splunk_data_ui_views.my_view]
}
`

func TestAccSplunkGenericAcl(t *testing.T) {
	appAclResourceName := "splunk_generic_acl.my_app"
	viewAclResourceName := "splunk_generic_acl.my_view"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newGenericAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appAclResourceName, "id", "apps/local/my_app"),
					resource.TestCheckResourceAttr(appAclResourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(appAclResourceName, "acl.0.read.1", "power"),
					resource.TestCheckNoResourceAttr(appAclResourceName, "acl.0.read.2"),
					resource.TestCheckResourceAttr(appAclResourceName, "acl.0.write.0", "admin"),
					resource.TestCheckResourceAttr(appAclResourceName, "acl.0.write.1", "power"),
					resource.TestCheckNoResourceAttr(appAclResourceName, "acl.0.write.2"),
					resource.TestCheckResourceAttr(viewAclResourceName, "id", "data/ui/views/my_view"),
					resource.TestCheckResourceAttr(viewAclResourceName, "acl.0.read.0", "admin"),
					resource.TestCheckResourceAttr(viewAclResourceName, "acl.0.read.1", "power"),
					resource.TestCheckNoResourceAttr(viewAclResourceName, "acl.0.read.2"),
					resource.TestCheckResourceAttr(viewAclResourceName, "acl.0.write.0", "admin"),
					resource.TestCheckResourceAttr(viewAclResourceName, "acl.0.write.1", "power"),
					resource.TestCheckNoResourceAttr(viewAclResourceName, "acl.0.write.2"),
				),
			},
			{
				Config: updateGenericAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appAclResourceName, "id", "apps/local/my_app"),
					resource.TestCheckResourceAttr(appAclResourceName, "acl.0.read.0", "admin"),
					resource.TestCheckNoResourceAttr(appAclResourceName, "acl.0.read.1"),
					resource.TestCheckResourceAttr(appAclResourceName, "acl.0.write.0", "admin"),
					resource.TestCheckNoResourceAttr(appAclResourceName, "acl.0.write.1"),
					resource.TestCheckResourceAttr(viewAclResourceName, "id", "data/ui/views/my_view"),
					resource.TestCheckResourceAttr(viewAclResourceName, "acl.0.read.0", "admin"),
					resource.TestCheckNoResourceAttr(viewAclResourceName, "acl.0.read.1"),
					resource.TestCheckResourceAttr(viewAclResourceName, "acl.0.write.0", "admin"),
					resource.TestCheckNoResourceAttr(viewAclResourceName, "acl.0.write.1"),
				),
			},
		},
	})
}
