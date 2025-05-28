package deployment

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func Test_findConfigMapManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name: "found ConfigMap",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
		},
		{
			name: "no match",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
`)),
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := findConfigMapManifests(tt.manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findSecretManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name: "found Secret",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
data:
  password: dGVzdA==
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
data:
  password: dGVzdA==
`)),
		},
		{
			name: "no match",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := findSecretManifests(tt.manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}
