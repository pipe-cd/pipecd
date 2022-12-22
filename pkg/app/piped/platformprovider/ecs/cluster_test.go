package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestParseClusterDefinition(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    types.Cluster
		expectedErr bool
	}{
		{
			name: "yaml format input",
			input: `
clusterArn: arn:aws:ecs:ap-northeast-1:XXXX:cluster/test-cluster
`,
			expected: types.Cluster{
				ClusterArn: aws.String("arn:aws:ecs:ap-northeast-1:XXXX:cluster/test-cluster"),
			},
		},
		{
			name: "json format input",
			input: `
{
	"clusterArn": "arn:aws:ecs:ap-northeast-1:XXXX:cluster/test-cluster"
}
`,
			expected: types.Cluster{
				ClusterArn: aws.String("arn:aws:ecs:ap-northeast-1:XXXX:cluster/test-cluster"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseClusterDefinition([]byte(tc.input))
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
