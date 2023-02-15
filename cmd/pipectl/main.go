// Copyright 2023 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/application"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/deployment"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/encrypt"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/event"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/piped"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/planpreview"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/quickstart"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"pipectl",
		"The command line tool for PipeCD.",
	)

	app.AddCommands(
		application.NewCommand(),
		deployment.NewCommand(),
		event.NewCommand(),
		planpreview.NewCommand(),
		piped.NewCommand(),
		encrypt.NewCommand(),
		quickstart.NewCommand(),
	)

	if err := app.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
