# Resource: splunk_saved_eventtypes
Create and manage saved searches.

## Example Usage
```
resource "splunk_saved_event_types" "test" {
    name        = "test"
    description = "Test New event description"
    disabled 	= "0"
    priority 	= 1
    search 		= "index=main"
    color		= "blue"
    tags 		= "tag"
    acl {
      owner = "admin"
      sharing = "app"
      app = "launcher"
    }
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTknowledge#saved.2Feventtypes

This resource block supports the following arguments:
* `name` - (Required) A name for the event type.
* `description` - (Optional) Human-readable description of this event type.
* `search` - (Required) Event type search string.
* `color`- (Optional) Color for this event type. The supported colors are: none, et_blue, et_green, et_magenta, et_orange, et_purple, et_red, et_sky, et_teal, et_yellow.
* `disabled` - (Optional) If True, disables the event type.
* `priority` - (Optional) Specify an integer from 1 to 10 for the value used to determine the order in which the matching event types of an event are displayed. 1 is the highest priority.
* `tags` - (Optional) [Deprecated] Use tags.conf.spec file to assign tags to groups of events with related field values.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the saved search event type
