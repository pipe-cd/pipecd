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

package log

import (
	golog "log"

	"github.com/spf13/cobra"
)

type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type Options struct {
	Level   string
	Disable bool
}

func (o *Options) RegisterPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.Level, "log-level", o.Level, "The minimum enabled logging level.")
	cmd.PersistentFlags().BoolVar(&o.Disable, "log-disable", o.Disable, "Whether the logger is disabled or not.")
}

func (o *Options) NewLogger() Logger {
	var level = Info
	if o.Level == "error" {
		level = Error
	}
	return &logger{
		level:   level,
		disable: o.Disable,
	}
}

const (
	Info = iota
	Error
)

type logger struct {
	level   int
	disable bool
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.disable {
		return
	}
	if l.level > Info {
		return
	}
	golog.Printf(format, v...)
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.disable {
		return
	}
	golog.Printf(format, v...)
}

func (l *logger) Fatal(format string, v ...interface{}) {
	if l.disable {
		return
	}
	golog.Fatalf(format, v...)
}
