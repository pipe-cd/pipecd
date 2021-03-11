module github.com/pipe-cd/pipe

go 1.16

require (
	cloud.google.com/go v0.65.0
	cloud.google.com/go/firestore v1.2.0
	cloud.google.com/go/storage v1.11.0
	github.com/Azure/go-autorest/autorest v0.10.2 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.5 // indirect
	github.com/DataDog/datadog-api-client-go v1.0.0-beta.16
	github.com/NYTimes/gziphandler v0.0.0-20170623195520-56545f4a5d46
	github.com/aws/aws-sdk-go v1.36.21 // indirect
	github.com/aws/aws-sdk-go-v2 v1.2.0
	github.com/aws/aws-sdk-go-v2/config v1.1.1
	github.com/aws/aws-sdk-go-v2/credentials v1.1.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.1.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.2.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fsouza/fake-gcs-server v1.21.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/goccy/go-yaml v1.8.8
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-github/v29 v29.0.3
	github.com/google/uuid v1.2.0
	github.com/googleapis/gnostic v0.2.2 // indirect
	github.com/hashicorp/golang-lru v0.5.3
	github.com/klauspost/compress v1.10.11 // indirect
	github.com/minio/minio-go/v7 v7.0.5
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.4.0
	go.uber.org/atomic v1.7.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.1-0.20190709142728-9a9fa7d4b5f0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/tools v0.0.0-20200916195026-c9a70fc28ce3 // indirect
	google.golang.org/api v0.31.0
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	istio.io/api v0.0.0-20200710191538-00b73d23c685
	k8s.io/api v0.18.9
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
	sigs.k8s.io/yaml v1.2.0
)
