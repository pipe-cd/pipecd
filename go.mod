module github.com/kapetaniosci/pipe

go 1.14

require (
	cloud.google.com/go v0.56.0
	cloud.google.com/go/storage v1.6.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	contrib.go.opencensus.io/exporter/stackdriver v0.12.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.3.5
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/hashicorp/golang-lru v0.5.1
	github.com/prometheus/client_golang v0.9.3
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.opencensus.io v0.22.3
	go.uber.org/atomic v1.4.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.1-0.20190709142728-9a9fa7d4b5f0
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	google.golang.org/api v0.20.0
	google.golang.org/grpc v1.28.0
	//k8s.io/api v0.0.0-20200410021914-5778e4f3d00d
	k8s.io/apimachinery v0.0.0-20200410021338-ff54c5b023af
	k8s.io/client-go v0.0.0-20200410022504-7b0589a2468d
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20200410021914-5778e4f3d00d
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200410021338-ff54c5b023af
	k8s.io/client-go => k8s.io/client-go v0.0.0-20200410022504-7b0589a2468d
)
