# Resource: splunk_apps_local
Create, install and manage apps on your Splunk instance

## Example Usage
```
resource "splunk_apps_local" "amazon_connect_app" {
  filename = true
  name = "/usr/home/amazon_connect_app_for_splunk-0.0.1.tar.gz"
  explicit_appname = "amazon_connect_app_for_splunk" // Original app name is required when filename is set to true
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.0/RESTREF/RESTapps#apps.2Flocal

This resource block supports the following arguments:
* `name` - (Required) Literal app name or path for the file to install, depending on the value of filename.
                      <br>filename = false indicates that name is the literal app name and that the app is created from Splunkbase using a template.
                      <br>filename = true indicates that name is the URL or path to the local .tar, .tgz or .spl file. If name is the Splunkbase URL, set auth or session to authenticate the request.
                      The app folder name cannot include spaces or special characters.
* `auth` - (Optional) Splunkbase session token for operations like install and update that require login. Use auth or session when installing or updating an app through Splunkbase.
* `author` - (Optional) For apps posted to Splunkbase, use your Splunk account username. For internal apps, include your name and contact information.
* `configured` - (Optional) Custom setup complete indication:
                            <br>true = Custom app setup complete.
                            <br>false = Custom app setup not complete.
* `description` - (Optional) Short app description also displayed below the app title in Splunk Web Launcher.
* `explicit_appname` - (Optional) Custom app name. Overrides name when installing an app from a file where filename is set to true. See also filename.
* `filename` - (Optional) Indicates whether to use the name value as the app source location.
                          <br>true indicates that name is a path to a file to install.
                          <br>false indicates that name is the literal app name and that the app is created from Splunkbase using a template.
* `label` - (Optional) App name displayed in Splunk Web, from five to eighty characters excluding the prefix "Splunk for".
* `session` - (Optional) Login session token for installing or updating an app on Splunkbase. Alternatively, use auth.
* `update` - (Optional) File-based update indication:
                         <br>true specifies that filename should be used to update an existing app. If not specified, update defaults to
                         <br>false, which indicates that filename should not be used to update an existing app.
* `version` - (Optional) App version.
* `visible` - (Optional) Indicates whether the app is visible and navigable from Splunk Web.
                         <br>true = App is visible and navigable.
                         <br>false = App is not visible or navigable.
* `acl` - (Optional) The app/user context that is the namespace for the resource

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the resource
