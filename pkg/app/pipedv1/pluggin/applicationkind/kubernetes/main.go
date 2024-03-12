package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"pipecd-plugin-kubernetes",
		"Plugin component to deploy Kubernetes Application.",
	)
	app.AddCommands(
		NewServerCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
