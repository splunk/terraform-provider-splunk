# Resource: splunk_data_ui_views
Create and manage splunk dashboards/views.
## Example Usage
```
resource "splunk_data_ui_views" "dashboard" {
  name     = "Terraform_Test_Dashboard"
  eai_data = "<dashboard><label>Terraform Test Dashboard</label></dashboard>"
  acl {
    owner = "admin"
    app = "search"
  }
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.1/RESTREF/RESTknowledge#data.2Fui.2Fviews

This resource block supports the following arguments:
* `name` - (Required) Dashboard name.
* `eai:data` - (Required) Dashboard XML definition.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the dashboard
