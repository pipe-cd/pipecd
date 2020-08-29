// Copyright 2020 The PipeCD Authors.
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
	"log"

	"github.com/pipe-cd/pipe/pkg/app/operator/cmd/server"
	"github.com/pipe-cd/pipe/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"operator",
		"A single component for operating owner tasks such as adding new project, deleting old data.",
	)
	app.AddCommands(
		server.NewCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
