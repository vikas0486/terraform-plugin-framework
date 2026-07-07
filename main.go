package main

import (
	"context"
	"log"

	"terraform-provider-thales/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {

	err := providerserver.Serve(
		context.Background(),
		provider.New,
		providerserver.ServeOpts{
			Address: "registry.terraform.io/provider/thales",
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
