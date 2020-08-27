# **Terraform provider for Splunk**


### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.12
-	[Go](https://golang.org/doc/install) go1.14.4 (to build the provider plugin)

### Building The Provider

Clone the repository: https://git.splunk.com/projects/GSA/repos/terraform-provider-splunk

Create go src directory and setup $GOPATH

Build the provider: `make build`

### Developing The Provider
* Use the Splunk REST API manual: https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTprolog to design resource schemas for the provider.
* Add a resource_x_test.go file to test the new resources' CRUD operations
* Before merging your changes lint your code by running `make fmt`
* Test the provider with the existing suite of provider tests before merging your changes
* Build the provider and test the new resources' CRUD and import operations before merging your changes

### Testing The Provider
* To run unit tests: `make test`
* To run acceptance tests: `make testacc`
  * Set the following variables to run acceptance tests `SPLUNK_HOME`, `SPLUNK_USERNAME`, `SPLUNK_URL`, `SPLUNK_PASSWORD`

### Using the provider

* Install Terraform
* Build the binary by `make build`
* Initialize terraform by `terraform init`
* Run `terraform plan` and `terraform apply` to apply configurations
* To update run `terraform plan` first to check config diff
* For importing existing resources use `terraform import`
* To remove all terraform managed resources use `terraform destroy`

#### Examples
* Use the `example.tf` provided in the repo to run `terraform plan` and `terraform apply` to apply configuration
  * Modify `provider "splunk"` resource block with proper instance details
* To update values modify the `example.tf` file and execute `terraform plan` and `terraform apply`
* Resource examples are also available in their respective resource_x_test.go files
* Examples to import existing configuration:
  * `terraform import splunk_inputs_http_event_collector.foo <hec-token-name>`
  * `terraform import splunk_inputs_script.bar "\$SPLUNK_HOME/etc/apps/splunk_instrumentation/bin/instrumentation.py"`
    * NOTE: Create a resource block first before importing resources (USAGE: https://www.terraform.io/docs/import/usage.html)
    * Example: `resource "splunk_inputs_http_event_collector" "foo" { }`
    `resource "splunk_inputs_scripts" "bar" { }`

#### Notes and Troubleshooting
* When importing an existing conf file, Splunk will respond with all default values for the conf file stanza (even if they do not appear explicitly in the stanza itself). These can be added to the associated `configs_conf` Terraform resource in your `.tf` file, otherwise they will show up as removed in the `terraform plan` diff. Although the plan will show them being removed, these default fields will <b>not</b> actually be modified or removed by Splunk.
* Testing errors mentioning to `too many open files` may be related to ulimits on your machine. Check current and increase the maximum number of open files `1024` using `ulimits -n 1024`
* If conf files are edited or deleted <b>manually</b>, restart Splunk to ensure state consistency before applying ore reapplying a template.
* Splunk environments with a large number of indexes, saved searches, knowledge objects, etc. may cause issues with the provided tests. To avoid these errors, use a fresh or lightly configured Splunk environment.
