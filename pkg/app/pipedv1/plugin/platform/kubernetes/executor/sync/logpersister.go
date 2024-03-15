package sync

import (
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/platform"
	"go.uber.org/zap"
)

type logger struct {
	stageLogger *zap.Logger
	stream      platform.ExecutorService_ExecuteStageServer
}

type LogPersister interface {
	Write(log []byte) (int, error)
	Info(log string)
	Infof(format string, a ...interface{})
	Success(log string)
	Successf(format string, a ...interface{})
	Error(log string)
	Errorf(format string, a ...interface{})
}

func NewLogPersister(stageLogger *zap.Logger, stream platform.ExecutorService_ExecuteStageServer) LogPersister {
	return &logger{
		stageLogger: stageLogger,
		stream:      stream,
	}
}

func (l *logger) Write(log []byte) (int, error) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: string(log),
	})
	l.stageLogger.Info(string(log))
	return len(log), nil
}

func (l *logger) Info(log string) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: log,
	})
	l.stageLogger.Info(log)
}

func (l *logger) Infof(format string, a ...interface{}) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: fmt.Sprintf(format, a...),
	})
	l.stageLogger.Info(fmt.Sprintf(format, a...))
}

func (l *logger) Success(log string) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: log,
	})
	l.stageLogger.Info(log)
}

func (l *logger) Successf(format string, a ...interface{}) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: fmt.Sprintf(format, a...),
	})
	l.stageLogger.Info(fmt.Sprintf(format, a...))
}

func (l *logger) Error(log string) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: log,
	})
	l.stageLogger.Error(log)
}

func (l *logger) Errorf(format string, a ...interface{}) {
	l.stream.Send(&platform.ExecuteStageResponse{
		Log: fmt.Sprintf(format, a...),
	})
	l.stageLogger.Error(fmt.Sprintf(format, a...))
}
