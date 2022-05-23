package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/myoung34/terraform-provider-uptycs/uptycs"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	err := providerserver.Serve(
		context.Background(),
		uptycs.New,
		providerserver.ServeOpts{
			Address: "snooguts.net/reddit/uptycs",
		},
	)

	if err != nil {
		panic(fmt.Sprintf("Error serving provider: %v", err))
	}
}
