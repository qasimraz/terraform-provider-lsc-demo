package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/qasimraz/terraform-provider-lsc-demo/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
