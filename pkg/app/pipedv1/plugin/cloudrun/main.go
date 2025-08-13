package main

import (
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrun/deployment"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrun/livestate"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"log"
)

func main() {
	plugin, err := sdk.NewPlugin(
		"0.0.1",
		sdk.WithDeploymentPlugin(&deployment.Plugin{}),
		sdk.WithLivestatePlugin(&livestate.Plugin{}),
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
