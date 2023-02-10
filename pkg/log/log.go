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

package log

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type EncodingType string

const (
	JSONEncoding     EncodingType = "json"
	ConsoleEncoding  EncodingType = "console"
	HumanizeEncoding EncodingType = "humanize"
)

const (
	DefaultLevel    = "info"
	DefaultEncoding = HumanizeEncoding
)

var (
	DefaultConfigs = Configs{
		Level:    DefaultLevel,
		Encoding: DefaultEncoding,
	}
)

type Configs struct {
	Level          string
	Encoding       EncodingType
	ServiceContext *ServiceContext
}

func NewLogger(c Configs) (*zap.Logger, error) {
	level := new(zapcore.Level)
	if err := level.Set(c.Level); err != nil {
		return nil, err
	}
	var options []zap.Option
	if c.ServiceContext != nil && c.Encoding != HumanizeEncoding {
		options = []zap.Option{
			zap.Fields(zap.Object("serviceContext", c.ServiceContext)),
		}
	}
	logger, err := newConfig(*level, c.Encoding).Build(options...)
	if err != nil {
		return nil, err
	}
	if c.ServiceContext != nil {
		return logger.Named(c.ServiceContext.Service), nil
	}
	return logger, nil
}

func newConfig(level zapcore.Level, encoding EncodingType) zap.Config {
	c := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig:    newEncoderConfig(encoding),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if encoding == HumanizeEncoding {
		c.Encoding = string(ConsoleEncoding)
		c.DisableCaller = true
	} else {
		c.Encoding = string(encoding)
	}
	return c
}

func newEncoderConfig(encoding EncodingType) zapcore.EncoderConfig {
	if encoding == HumanizeEncoding {
		return zapcore.EncoderConfig{
			TimeKey:        "eventTime",
			LevelKey:       "",
			NameKey:        "",
			CallerKey:      "",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	return zapcore.EncoderConfig{
		TimeKey:        "eventTime",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARNING")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("CRITICAL")
	case zapcore.PanicLevel:
		enc.AppendString("ALERT")
	case zapcore.FatalLevel:
		enc.AppendString("EMERGENCY")
	}
}

type ServiceContext struct {
	Service string
	Version string
}

func (sc ServiceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if sc.Service == "" {
		return errors.New("service name is mandatory")
	}
	enc.AddString("service", sc.Service)
	enc.AddString("version", sc.Version)
	return nil
}
