package deployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestPlugin_executeApplyStage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *sdk.ExecuteStageInput[config.ApplicationConfigSpec]
		dts     []*sdk.DeployTarget[config.DeployTargetConfig]
		want    sdk.StageStatus
	}{
		{
			name: "failure when StageLogPersister is unavailable",
			input: &sdk.ExecuteStageInput[config.ApplicationConfigSpec]{
				Client: sdk.NewClient(nil, "terraform", "app1", "stage1", nil, nil),
				Logger: zap.NewNop(),
			},
			dts:  []*sdk.DeployTarget[config.DeployTargetConfig]{{}},
			want: sdk.StageStatusFailure,
		},
	}

	p := &Plugin{}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := p.executeApplyStage(context.Background(), tc.input, tc.dts)
			assert.Equal(t, tc.want, got)
		})
	}
}
