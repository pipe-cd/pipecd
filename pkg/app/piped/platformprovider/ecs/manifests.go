package ecs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pipe-cd/pipecd/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type Manifest struct {
	u *unstructured.Unstructured
}

func LoadPlainYAMLManifests(dir string, names []string, configFileName string) ([]Manifest, error) {
	if len(names) == 0 {
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == dir {
				return nil
			}
			if f.IsDir() {
				return filepath.SkipDir
			}
			ext := filepath.Ext(f.Name())
			if ext != ".yaml" && ext != ".yml" && ext != ".json" {
				return nil
			}
			if model.IsApplicationConfigFile(f.Name()) {
				return nil
			}
			if f.Name() == configFileName {
				return nil
			}
			names = append(names, f.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	manifests := make([]Manifest, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		ms, err := LoadManifestFromYAMLFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load manifest at %s (%w)", path, err)
		}
		manifests = append(manifests, ms)
	}

	return manifests, nil
}

func LoadManifestFromYAMLFile(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	return parseManifest(data)
}

func parseManifest(data []byte) (Manifest, error) {
	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return Manifest{}, err
	}

	return Manifest{u: &obj}, nil
}
