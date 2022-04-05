package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-uptycs/uptycs"
)

func main() {
	tfsdk.Serve(context.Background(), uptycs.New, tfsdk.ServeOpts{
		Name: "uptycs",
	})
}
