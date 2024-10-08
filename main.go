package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/dmalch/terraform-provider-myheritage/internal"
)

func main() {
	err := providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "github.com/dmalch/myheritage",
	})
	if err != nil {
		log.Fatal(err)
	}
}
