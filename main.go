package main

import (
	"context"

	"github.com/DavidKrau/terraform-provider-simplemdm/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name simplemdm

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "freenow.com/terraform-provider/simplemdm",
	})
	if err != nil {
		log.Panic().Msgf("starting server failed")
	}
}
