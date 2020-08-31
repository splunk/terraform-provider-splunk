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
* Add all necessary documentation in the docs folder

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

**NOTE:** When developing or testing with terraform `>= 0.13` you must replace the provider location from remote (registry.terraform.io) to local build.

Follow guidelines: https://github.com/hashicorp/terraform/blob/master/website/upgrade-guides/0-13.html.markdown

#### Examples
* Use the `example.tf` provided in the repo to run `terraform plan` and `terraform apply` to apply configuration
  * Modify `provider "splunk"` resource block with proper instance details
* To update values modify the `example.tf` file and execute `terraform plan` and `terraform apply`
* Resource examples are also available in their respective docs/resources folder

**NOTE:** Create a resource block first before importing resources (USAGE: https://www.terraform.io/docs/import/usage.html)

#### Notes and Troubleshooting
* When conflicts arise during resource creation, import the resource first using `terraform import` command and make modifications to the resource.
* Testing errors mentioning `too many open files` may be related to `ulimits` on your machine. Check current and increase the maximum number of open files `1024` using `ulimits -n 1024`
* If conf files are edited or deleted <b>manually</b>, restart Splunk to ensure state consistency before applying or reapplying a template.
* Splunk environments with numerous indexes, saved searches, knowledge objects, etc. may cause issues with the provided tests. To avoid these errors, use a fresh or lightly configured Splunk environment.
