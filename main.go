package main

import (
	"github.com/hashicorp/terraform/plugin v0.11.3"
	"github.com/hashicorp/terraform/terraform v0.11.3"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}
