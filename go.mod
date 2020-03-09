module github.com/nghialv/dianomi

go 1.14

require (
	cloud.google.com/go v0.43.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	contrib.go.opencensus.io/exporter/stackdriver v0.12.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/prometheus/client_golang v0.9.3
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.opencensus.io v0.22.0
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.1-0.20190709142728-9a9fa7d4b5f0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
	google.golang.org/api v0.7.0
	google.golang.org/grpc v1.22.0
	k8s.io/code-generator v0.0.0-20191220033320-6b257a9d6f46
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191221033533-72223a9f9901
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191221033353-3253b0a30d67
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191222113738-1b1a35e41a57
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191220033320-6b257a9d6f46
)
