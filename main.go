package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/myoung34/terraform-provider-uptycs/uptycs"
)

func main() {
	err := providerserver.Serve(
		context.Background(),
		uptycs.New,
		providerserver.ServeOpts{},
	)

	if err != nil {
		panic(fmt.Sprintf("Error serving provider: %v", err))
	}
}
