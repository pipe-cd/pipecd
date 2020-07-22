package api

import (
	"reflect"
	"testing"

	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
)

func Test_filterDeploymentConfigTemplates(t *testing.T) {
	type args struct {
		labels    []webservice.DeploymentConfigTemplateLabel
		templates []*webservice.DeploymentConfigTemplate
	}
	tests := []struct {
		name string
		args args
		want []*webservice.DeploymentConfigTemplate
	}{
		{
			name: "Specify just one label",
			args: args{
				labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
				templates: []*webservice.DeploymentConfigTemplate{
					{
						Name:   "Canary",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
					},
					{
						Name:   "Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
				},
			},
			want: []*webservice.DeploymentConfigTemplate{
				{
					Name:   "Canary",
					Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
				},
			},
		},
		{
			name: "Two labels specified, non-existent",
			args: args{
				labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
				templates: []*webservice.DeploymentConfigTemplate{
					{
						Name:   "Canary",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
					},
					{
						Name:   "Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
				},
			},
			want: []*webservice.DeploymentConfigTemplate{},
		},
		{
			name: "Two labels specified, existent",
			args: args{
				labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
				templates: []*webservice.DeploymentConfigTemplate{
					{
						Name:   "Canary Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
					{
						Name:   "Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
				},
			},
			want: []*webservice.DeploymentConfigTemplate{
				{
					Name:   "Canary Blue/Green",
					Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterDeploymentConfigTemplates(tt.args.templates, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterDeploymentConfigTemplates() = %v, want %v", got, tt.want)
			}
		})
	}
}
