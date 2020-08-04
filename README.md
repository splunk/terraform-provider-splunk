# **Terraform provider for Splunk**


### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.11.8
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

* Install `terraform`
* `make build`
* `terraform init`
* Use the `example.tf` to run `terraform plan` and `terraform apply` to apply configuration
* For importing existing resources use `terraform import`
  * Example: `terraform import splunk_inputs_http_event_collector.foo <hec-token>`
  * NOTE: Create a resource block first before importing resources (USAGE: https://www.terraform.io/docs/import/usage.html)
  
