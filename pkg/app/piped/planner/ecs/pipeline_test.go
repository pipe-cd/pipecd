package ecs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
			stages := buildQuickSyncPipeline(tc.wantAutoRollback, time.Now())
			var autoRollback bool
			for _, stage := range stages {
				if stage.Name == string(model.StageRollback) {
					autoRollback = true
				}
			}
			assert.Equal(t, tc.wantAutoRollback, autoRollback)
		})
	}
}
