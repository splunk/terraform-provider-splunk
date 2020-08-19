# **Terraform provider for Splunk**


### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.12
-	[Go](https://golang.org/doc/install) go1.14.4 (to build the provider plugin)

### Building The Provider

Clone the repository: https://git.splunk.com/projects/GSA/repos/terraform-provider-splunk

Create go src directory and setup $GOPATH

Build the provider: `make build`

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
* Examples to import existing configuration:
  * `terraform import splunk_inputs_http_event_collector.foo <hec-token-name>`
  * `terraform import splunk_inputs_script.bar "\$SPLUNK_HOME/etc/apps/splunk_instrumentation/bin/instrumentation.py"`
    * NOTE: Create a resource block first before importing resources (USAGE: https://www.terraform.io/docs/import/usage.html)
    * Example: `resource "splunk_inputs_http_event_collector" "foo" { }`
    `resource "splunk_inputs_scripts" "bar" { }`
