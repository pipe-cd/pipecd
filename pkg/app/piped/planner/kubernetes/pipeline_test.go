package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestBuildQuickSyncPipeline(t *testing.T) {
	tests := []struct {
		name             string
		wantAutoRollback bool
	}{
		{
			name:             "want auto rollback stage",
			wantAutoRollback: true,
		},
		{
			name:             "don't want auto rollback stage",
			wantAutoRollback: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStages := buildQuickSyncPipeline(tc.wantAutoRollback, time.Now())
			var gotAutoRollback bool
			for _, stage := range gotStages {
				if stage.Name == string(model.StageRollback) {
					gotAutoRollback = true
				}
			}
			assert.Equal(t, tc.wantAutoRollback, gotAutoRollback)
		})
	}
}

func TestBuildProgressivePipeline(t *testing.T) {
	tests := []struct {
		name             string
		wantAutoRollback bool
	}{
		{
			name:             "want auto rollback stage",
			wantAutoRollback: true,
		},
		{
			name:             "don't want auto rollback stage",
			wantAutoRollback: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStages := buildProgressivePipeline(&config.DeploymentPipeline{}, tc.wantAutoRollback, time.Now())
			var gotAutoRollback bool
			for _, stage := range gotStages {
				if stage.Name == string(model.StageRollback) {
					gotAutoRollback = true
				}
			}
			assert.Equal(t, tc.wantAutoRollback, gotAutoRollback)
		})
	}
}
