package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-splunk/splunk"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: splunk.Provider,
	}
	plugin.Serve(&opts)
}
