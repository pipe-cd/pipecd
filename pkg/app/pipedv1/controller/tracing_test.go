package controller

import (
	"testing"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentTraceID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		d1     *model.Deployment
		d2     *model.Deployment
		assert assert.ComparisonAssertionFunc
	}{
		{
			name: "same deployment id",
			d1: &model.Deployment{
				Id: "example-deployment-id",
			},
			d2: &model.Deployment{
				Id: "example-deployment-id",
			},
			assert: assert.Equal,
		},
		{
			name: "different deployment id",
			d1: &model.Deployment{
				Id: "example-deployment-id",
			},
			d2: &model.Deployment{
				Id: "example-deployment-id-other",
			},
			assert: assert.NotEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id1 := deploymentTraceID(tt.d1)
			id2 := deploymentTraceID(tt.d2)

			tt.assert(t, id1, id2)
		})
	}
}
