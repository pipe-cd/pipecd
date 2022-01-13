module github.com/pipe-cd/pipecd

go 1.16

require (
	cloud.google.com/go v0.65.0
	cloud.google.com/go/firestore v1.2.0
	cloud.google.com/go/storage v1.11.0
	github.com/DataDog/datadog-api-client-go v1.0.0-beta.16
	github.com/NYTimes/gziphandler v0.0.0-20170623195520-56545f4a5d46
	github.com/aws/aws-sdk-go-v2 v1.6.0
	github.com/aws/aws-sdk-go-v2/config v1.1.1
	github.com/aws/aws-sdk-go-v2/credentials v1.1.1
	github.com/aws/aws-sdk-go-v2/internal/ini v1.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.1.1
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.3.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.1.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.2.0
	github.com/creasty/defaults v1.5.2
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fsouza/fake-gcs-server v1.21.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/goccy/go-yaml v1.9.3
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-github/v29 v29.0.3
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/golang-lru v0.5.3
	github.com/minio/minio-go/v7 v7.0.5
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.7.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.1-0.20190709142728-9a9fa7d4b5f0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/api v0.31.0
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.27.1
	istio.io/api v0.0.0-20200710191538-00b73d23c685
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.3
	sigs.k8s.io/yaml v1.2.0
)
