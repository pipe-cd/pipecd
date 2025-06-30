package main

import (
	"log"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func main() {
	plugin, err := sdk.NewPlugin("0.0.1", sdk.WithStagePlugin(&plugin{}))
	if err != nil {
		log.Fatalln(err)
	}

	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
