package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/rsrdesarrollo/terraform-provider-splunk/splunk"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: splunk.Provider,
	}
	plugin.Serve(&opts)
}
