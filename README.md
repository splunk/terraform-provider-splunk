# **Terraform provider for Splunk**


### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.11.8
-	[Go](https://golang.org/doc/install) go1.14.4 (to build the provider plugin)

### Building The Provider

Clone the repository: https://git.splunk.com/projects/GSA/repos/terraform-provider-splunk

Create go src directory and setup $GOPATH

Build the provider: `make build`

### Testing The Provider
To run unit tests: `make test`
To run acceptance tests: `set TF_ACC=1`

### Using the provider

* Install `terraform`
* `make default` or run `terraform plan` (Requires correct credentials to splunk instance)
* Use the `example.tf` to run `terraform plan` and `terraform apply`