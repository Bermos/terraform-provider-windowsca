package main

import (
	"context"
	"github.com/bermos/terraform-provider-windowsca/windowsca"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), windowsca.New, providerserver.ServeOpts{
		Address: "hashicorp.com/edu/windowsca", // TODO: change this to your provider's address
	})
}
