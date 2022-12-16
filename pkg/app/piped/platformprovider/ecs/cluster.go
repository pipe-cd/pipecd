package ecs

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"sigs.k8s.io/yaml"
)

func loadClusterDefinition(path string) (types.Cluster, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.Cluster{}, err
	}

	return parseClusterDefinition(data)
}

func parseClusterDefinition(data []byte) (types.Cluster, error) {
	var obj types.Cluster
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return types.Cluster{}, err
	}

	return obj, nil
}
