# **Terraform provider for Splunk**


### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.11.8
-	[Go](https://golang.org/doc/install) go1.14.4 (to build the provider plugin)

### Building The Provider

Clone the repository: https://git.splunk.com/projects/GSA/repos/terraform-provider-splunk

Create go src directory and setup $GOPATH

Build the provider: `go build -o terraform-provider-splunk .`

### Using the provider

* Install `terraform`
* Run `terraform plan`
* Use the `example.tf` to run `terraform plan` and `terraform apply`