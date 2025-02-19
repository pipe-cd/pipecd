package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func main() {
	sdk.RegisterPipelineSyncPlugin(&plugin{})

	if err := sdk.Run(); err != nil {
		log.Fatalln(err)
	}
}
