# Resource: splunk_configs_conf
Create and manage configuration file stanzas.

## Example Usage
```
resource "splunk_configs_conf" "new-conf-stanza" {
  name = "custom-conf/custom"
  variables = {
    "disabled" : "false"
    "custom_key" : "value"
  }
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTconf#configs.2Fconf-.7Bfile.7D

This resource block supports the following arguments:
* `name` - (Required) A '/' separated string consisting of {conf_file_name}/{stanza_name} ex. props/custom_stanza
* `variables` - (Optional) A map of key value pairs for a stanza.
* `acl` - (Optional) The app/user context that is the namespace for the resource

**NOTE:** When importing an existing conf file, Splunk will respond with all default values for the conf file stanza (even if they do not appear explicitly in the stanza itself). These can be added to the associated `configs_conf` Terraform resource in your `.tf` file, otherwise they will show up as removed in the `terraform plan` diff. <b>Although the plan will show them being removed, these default fields will <b>not</b> actually be modified or removed by Splunk.</b>

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the resource
