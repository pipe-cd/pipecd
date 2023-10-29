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

package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/pipe-cd/pipecd/pkg/log"
	"github.com/pipe-cd/pipecd/pkg/version"
)

type App struct {
	rootCmd        *cobra.Command
	telemetryFlags TelemetryFlags
}

var ErrFlagParse = errors.New("FlagParseErr")

func NewApp(name, desc string) *App {
	a := &App{
		rootCmd: &cobra.Command{
			Use:           name,
			Short:         desc,
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		telemetryFlags: defaultTelemetryFlags,
	}
	a.rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return ErrFlagParse
	})
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the information of current binary.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Get())
		},
	}
	a.rootCmd.AddCommand(versionCmd)
	a.setGlobalFlags()
	return a
}

func (a *App) AddCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		a.rootCmd.AddCommand(cmd)
	}
}

func (a *App) Run() error {
	return a.rootCmd.Execute()
}

type TelemetryFlags struct {
	LogLevel                string
	LogEncoding             string
	Profile                 bool
	ProfileDebugLogging     bool
	ProfilerCredentialsFile string
	Metrics                 bool
}

var defaultTelemetryFlags = TelemetryFlags{
	LogLevel:    string(log.DefaultLevel),
	LogEncoding: string(log.DefaultEncoding),
	Metrics:     true,
}

func (a *App) setGlobalFlags() {
	a.rootCmd.PersistentFlags().StringVar(
		&a.telemetryFlags.LogLevel,
		"log-level",
		a.telemetryFlags.LogLevel,
		"The minimum enabled logging level.",
	)
	a.rootCmd.PersistentFlags().StringVar(
		&a.telemetryFlags.LogEncoding,
		"log-encoding",
		a.telemetryFlags.LogEncoding,
		"The encoding type for logger [json|console|humanize].",
	)
	a.rootCmd.PersistentFlags().BoolVar(
		&a.telemetryFlags.Profile,
		"profile",
		a.telemetryFlags.Profile,
		"If true enables uploading the profiles to Stackdriver.",
	)
	a.rootCmd.PersistentFlags().BoolVar(
		&a.telemetryFlags.ProfileDebugLogging,
		"profile-debug-logging",
		a.telemetryFlags.ProfileDebugLogging,
		"If true enables logging debug information of profiler.",
	)
	a.rootCmd.PersistentFlags().StringVar(
		&a.telemetryFlags.ProfilerCredentialsFile,
		"profiler-credentials-file",
		a.telemetryFlags.ProfilerCredentialsFile,
		"The path to the credentials file using while sending profiles to Stackdriver.",
	)
	a.rootCmd.PersistentFlags().BoolVar(
		&a.telemetryFlags.Metrics,
		"metrics",
		a.telemetryFlags.Metrics,
		"Whether metrics is enabled or not.",
	)
}

func parseTelemetryFlags(fs *pflag.FlagSet) (TelemetryFlags, error) {
	flags := defaultTelemetryFlags

	// Extract log-level.
	if fs.Lookup("log-level") != nil {
		s, err := fs.GetString("log-level")
		if err != nil {
			return flags, err
		}
		flags.LogLevel = s
	}

	// Extract log-encoding.
	if fs.Lookup("log-encoding") != nil {
		s, err := fs.GetString("log-encoding")
		if err != nil {
			return flags, err
		}
		flags.LogEncoding = s
	}

	// Extract profile.
	if fs.Lookup("profile") != nil {
		b, err := fs.GetBool("profile")
		if err != nil {
			return flags, err
		}
		flags.Profile = b
	}

	// Extract profile-debug-logging.
	if fs.Lookup("profile-debug-logging") != nil {
		b, err := fs.GetBool("profile-debug-logging")
		if err != nil {
			return flags, err
		}
		flags.ProfileDebugLogging = b
	}

	// Extract profiler-credentials-file.
	if fs.Lookup("profiler-credentials-file") != nil {
		s, err := fs.GetString("profiler-credentials-file")
		if err != nil {
			return flags, err
		}
		flags.ProfilerCredentialsFile = s
	}

	// Extract metrics.
	if fs.Lookup("metrics") != nil {
		b, err := fs.GetBool("metrics")
		if err != nil {
			return flags, err
		}
		flags.Metrics = b
	}

	return flags, nil
}
