# **Terraform provider for Splunk**


### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.12
-	[Go](https://golang.org/doc/install) go1.14.4 (to build the provider plugin)

### Building The Provider

Clone the [repository](https://github.com/splunk/terraform-provider-splunk/)

Create go src directory and setup $GOPATH

Build the provider: `make build`

### Developing The Provider
* Use the [Splunk REST API manual](https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTprolog) to design resource schemas for the provider.
* Add a resource_x_test.go file to test the new resources' CRUD operations
* Before merging your changes lint your code by running `make fmt`
* Test the provider with the existing suite of provider tests before merging your changes
* Build the provider and test the new resources' CRUD and import operations before merging your changes
* Add all necessary documentation in the docs folder
* As per [best practices](https://www.terraform.io/docs/extend/best-practices/versioning.html), update changelog.md and version as required

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
* Use `terraform refresh` for state migrations to be applied

**NOTE:** When developing and testing local provider builds, if terraform version `>= 0.13` you would have to replace the provider binaries in the `.terraform` folder with your local build. [Follow these guidelines](https://github.com/hashicorp/terraform/blob/master/website/upgrade-guides/0-13.html.markdown)

### Contributions
We are open to contributions!
<p>Please follow development guidelines and feel free to open a PR against the `master` branch with your changes. The PR should trigger the GitHub actions to run
both the unit and acceptance tests. Once all tests have passed, our team will review, make suggestions, approve, and merge the PR.
After merging, our team will update the changelog.MD file and create a version tag that should automatically create a new release.</p>

#### Examples
* The examples folder contains a few basic splunk provider examples, run `terraform init` and `terraform apply` to apply these example configuration.
* Resource examples are also available in their respective docs/resources folder

**NOTE:** Create a resource block first before importing resources. Docs on the import [usage](https://www.terraform.io/docs/import/usage.html)

#### Notes and Troubleshooting
* When conflicts arise during resource creation, import the resource first using `terraform import` command and make modifications to the resource.
* The error `too many open files` may be due to `ulimit` settings on your machine. Check current and increase the maximum number of open files `1024` using `ulimit -n 1024`
* When deleting or editing conf files <b>manually</b>, restart Splunk to ensure state consistency before applying or reapplying a template.
* Splunk environment with numerous indexes, saved searches, knowledge objects, etc. may cause issues with the provided tests. To avoid these errors, use a fresh or lightly configured Splunk environment.

### Support
Use the [GitHub issue tracker](https://github.com/splunk/terraform-provider-splunk/issues) to submit bugs or request features.
* Please add the Terraform and provider version, and the version of Splunk Enterprise used.

[Splunk Ideas](https://ideas.splunk.com/) is another place for your suggestions and [Splunk Answers](https://community.splunk.com/t5/Community/ct-p/en-us) for questions.
