package plugin

import (
	"os/exec"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/platform"
)

func GetPlanSourceCloner(input *platform.PlanPluginInput) (deploysource.SourceCloner, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	cloner := deploysource.NewLocalSourceCloner(
		git.NewRepo(input.GetSourceRemoteUrl(), gitPath, input.GetSourceRemoteUrl(), input.GetDeployment().GetGitPath().GetRepo().GetBranch(), nil),
		"target",
		input.GetDeployment().GetGitPath().GetRepo().GetBranch(),
	)

	return cloner, nil
}
