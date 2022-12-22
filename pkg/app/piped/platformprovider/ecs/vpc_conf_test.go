package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestParseVpcConf(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    types.AwsVpcConfiguration
		expectedErr bool
	}{
		{
			name: "yaml format input",
			input: `
assignPublicIp: ENABLED
securityGroups:
  - sg-YYYY
subnets:
  - subnet-XXXX
  - subnet-YYYY
`,
			expected: types.AwsVpcConfiguration{
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{
					"sg-YYYY",
				},
				Subnets: []string{
					"subnet-XXXX",
					"subnet-YYYY",
				},
			},
		},
		{
			name: "json format input",
			input: `
{
  "assignPublicIp": "ENABLED",
  "securityGroups": [
    "sg-YYYY"
  ],
  "subnets": [
    "subnet-XXXX",
    "subnet-YYYY"
  ]
}
`,
			expected: types.AwsVpcConfiguration{
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{
					"sg-YYYY",
				},
				Subnets: []string{
					"subnet-XXXX",
					"subnet-YYYY",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseVpcConfig([]byte(tc.input))
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
