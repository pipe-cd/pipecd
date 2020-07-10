module github.com/pipe-cd/pipe

go 1.14

require (
	cloud.google.com/go v0.56.0
	cloud.google.com/go/firestore v1.2.0
	cloud.google.com/go/storage v1.6.0
	github.com/NYTimes/gziphandler v0.0.0-20170623195520-56545f4a5d46
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.4.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.4.0
	github.com/google/uuid v1.1.1
	github.com/hashicorp/golang-lru v0.5.1
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.uber.org/atomic v1.4.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.1-0.20190709142728-9a9fa7d4b5f0
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	google.golang.org/api v0.20.0
	google.golang.org/grpc v1.28.1
	google.golang.org/protobuf v1.21.0
	istio.io/api v0.0.0-20200710191538-00b73d23c685
	k8s.io/api v0.0.0-20200410021914-5778e4f3d00d
	//k8s.io/api v0.0.0-20200410021914-5778e4f3d00d
	k8s.io/apimachinery v0.18.1
	k8s.io/client-go v0.0.0-20200410022504-7b0589a2468d
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20200410021914-5778e4f3d00d
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200410021338-ff54c5b023af
	k8s.io/client-go => k8s.io/client-go v0.0.0-20200410022504-7b0589a2468d
)
