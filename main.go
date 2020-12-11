package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/jaware-splunk/terraform-provider-splunk/splunk"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: splunk.Provider,
	}
	plugin.Serve(&opts)
}
