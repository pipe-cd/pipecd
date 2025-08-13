package livestate

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrun/config"
	"github.com/pipe-cd/piped-plugin-sdk-go"
)

type Plugin struct{}

func (p Plugin) GetLivestate(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.CloudRunDeployTargetConfig], input *sdk.GetLivestateInput[config.CloudRunApplicationSpec]) (*sdk.GetLivestateResponse, error) {
	//TODO implement me
	panic("implement me")
}
