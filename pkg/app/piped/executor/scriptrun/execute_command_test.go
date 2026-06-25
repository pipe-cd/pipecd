package scriptrun

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeLogPersister struct{}

func (l *fakeLogPersister) Write(b []byte) (int, error)         { fmt.Println(string(b)); return len(b), nil }
func (l *fakeLogPersister) Info(s string)                       { fmt.Println("INFO:", s) }
func (l *fakeLogPersister) Infof(s string, args ...interface{}) { fmt.Printf("INFO: "+s+"\n", args...) }
func (l *fakeLogPersister) Success(s string)                    { fmt.Println("SUCCESS:", s) }
func (l *fakeLogPersister) Successf(s string, args ...interface{}) { fmt.Printf("SUCCESS: "+s+"\n", args...) }
func (l *fakeLogPersister) Error(s string)                      { fmt.Println("ERROR:", s) }
func (l *fakeLogPersister) Errorf(s string, args ...interface{})   { fmt.Printf("ERROR: "+s+"\n", args...) }

func TestExecuteCommandCancellation(t *testing.T) {
	// Ensure that executeCommand correctly stops when its context is canceled.
	e := &Executor{
		Input: executor.Input{
			Deployment:   &model.Deployment{
				Id: "deploy-1",
				Trigger: &model.DeploymentTrigger{
					Commit: &model.Commit{
						Hash: "hash-1",
					},
				},
				GitPath: &model.ApplicationGitPath{
					Repo: &model.ApplicationGitRepository{
						Id:     "repo-1",
						Remote: "repo-url",
					},
				},
			},
			LogPersister: &fakeLogPersister{},
			StageConfig: config.PipelineStage{
				ScriptRunStageOptions: &config.ScriptRunStageOptions{
					Run: "exec sleep 10",
				},
			},
		},
		appDir: os.TempDir(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	// Run command in a goroutine
	done := make(chan model.StageStatus)
	go func() {
		done <- e.executeCommand(ctx)
	}()

	// Cancel almost immediately
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case status := <-done:
		assert.Equal(t, model.StageStatus_STAGE_FAILURE, status, "Expected command to fail on cancellation")
	case <-time.After(2 * time.Second):
		t.Fatal("executeCommand did not return after context cancellation")
	}
}
