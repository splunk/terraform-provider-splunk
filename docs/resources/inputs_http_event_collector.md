# Resource: splunk_inputs_http_event_collector
Create or update HTTP Event Collector input configuration tokens.

## Example Usage
```
resource "splunk_inputs_http_event_collector" "hec-token-01" {
  name       = "hec-token-01"
  index      = "main"
  indexes    = ["main", "history", "summary"]
  source     = "new:source"
  sourcetype = "new:sourcetype"
  disabled   = false
  use_ack    = false
  acl {
    owner   = "user01"
    sharing = "global"
    read    = ["admin"]
    write   = ["admin"]
  }
}
```

```
terraform import splunk_inputs_http_event_collector.token01 <hec-token-name>
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTinput#data.2Finputs.2Fhttp

This resource block supports the following arguments:
* `name` - (Required) Token name (inputs.conf key)
* `token` - (Optional) Token value for sending data to collector/event endpoint
* `index` - (Optional) Index to store generated events
* `indexes` - (Optional) Set of indexes allowed for events with this token
* `host` - (Optional) Default host value for events with this token
* `source` - (Optional) Default source for events with this token
* `sourcetype` - (Optional) Default source type for events with this token
* `disabled` - (Optional) Input disabled indicator
* `use_ack` - (Optional) Indexer acknowledgement for this token
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the http event collector resource
