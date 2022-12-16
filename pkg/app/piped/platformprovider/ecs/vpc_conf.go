package ecs

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"sigs.k8s.io/yaml"
)

func loadVpcConfig(path string) (types.AwsVpcConfiguration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.AwsVpcConfiguration{}, err
	}

	return parseVpcConfig(data)
}

func parseVpcConfig(data []byte) (types.AwsVpcConfiguration, error) {
	var obj types.AwsVpcConfiguration
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return types.AwsVpcConfiguration{}, err
	}

	return obj, nil
}
