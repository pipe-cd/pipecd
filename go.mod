module github.com/pipe-cd/pipecd

go 1.26.2

require (
	cloud.google.com/go/firestore v1.22.0
	cloud.google.com/go/profiler v0.3.1
	cloud.google.com/go/secretmanager v1.20.0
	cloud.google.com/go/storage v1.65.0
	github.com/DataDog/datadog-api-client-go v1.16.0
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/NYTimes/gziphandler v1.1.1
	github.com/aws/aws-sdk-go-v2 v1.43.4
	github.com/aws/aws-sdk-go-v2/config v1.27.38
	github.com/aws/aws-sdk-go-v2/credentials v1.17.36
	github.com/aws/aws-sdk-go-v2/service/ecs v1.46.2
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.38.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.88.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.107.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.33.2
	github.com/aws/aws-sdk-go-v2/service/ssm v1.54.3
	github.com/coreos/go-oidc/v3 v3.21.0
	github.com/creasty/defaults v1.8.0
	github.com/envoyproxy/go-control-plane/envoy v1.39.0
	github.com/envoyproxy/protoc-gen-validate v1.3.3
	github.com/fsouza/fake-gcs-server v1.56.1
	github.com/go-logr/logr v1.4.4
	github.com/go-sql-driver/mysql v1.10.1
	github.com/goccy/go-yaml v1.19.2
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.4
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.7.0
	github.com/google/go-github/v29 v29.0.3
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/golang-lru v0.5.3
	github.com/hashicorp/hcl/v2 v2.24.0
	github.com/minio/minio-go/v7 v7.3.0
	github.com/opencontainers/image-spec v1.1.1
	github.com/ory/dockertest/v3 v3.12.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.24.1
	github.com/prometheus/common v0.70.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/slack-go/slack v0.29.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
	go.uber.org/atomic v1.11.0
	go.uber.org/mock v0.6.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.55.0
	golang.org/x/mod v0.38.0
	golang.org/x/net v0.58.0
	golang.org/x/oauth2 v0.37.0
	golang.org/x/sync v0.23.0
	google.golang.org/api v0.293.0
	google.golang.org/grpc v1.83.1
	google.golang.org/protobuf v1.36.12
	istio.io/api v0.0.0-20200710191538-00b73d23c685
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v0.24.3
	oras.land/oras-go/v2 v2.6.2
	sigs.k8s.io/controller-runtime v0.12.3
	sigs.k8s.io/yaml v1.5.0
)

require (
	cel.dev/expr v0.25.2 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.23.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.11.0 // indirect
	cloud.google.com/go/longrunning v1.2.0 // indirect
	cloud.google.com/go/monitoring v1.29.0 // indirect
	cloud.google.com/go/pubsub/v2 v2.6.2 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.33.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.57.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.57.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.23.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.27.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.31.2 // indirect
	github.com/aws/smithy-go v1.27.6 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/containerd/continuity v0.4.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/cli v27.4.1+incompatible // indirect
	github.com/docker/docker v28.0.0+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/emicklei/go-restful v2.16.0+incompatible // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v1.0.0 // indirect
	github.com/go-openapi/jsonreference v1.0.0 // indirect
	github.com/go-openapi/swag v0.28.0 // indirect
	github.com/go-openapi/swag/cmdutils v0.28.0 // indirect
	github.com/go-openapi/swag/conv v0.28.0 // indirect
	github.com/go-openapi/swag/fileutils v0.28.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.28.0 // indirect
	github.com/go-openapi/swag/loading v0.28.0 // indirect
	github.com/go-openapi/swag/mangling v0.28.0 // indirect
	github.com/go-openapi/swag/netutils v0.28.0 // indirect
	github.com/go-openapi/swag/pools v0.28.0 // indirect
	github.com/go-openapi/swag/stringutils v0.28.0 // indirect
	github.com/go-openapi/swag/typeutils v0.28.0 // indirect
	github.com/go-openapi/swag/yamlutils v0.28.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20221103000818-d260c55eee4c // indirect
	github.com/google/renameio/v2 v2.0.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.20 // indirect
	github.com/googleapis/gax-go/v2 v2.23.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runc v1.3.6 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/xattr v0.4.12 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spiffe/go-spiffe/v2 v2.7.0 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/zclconf/go-cty v1.16.3 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.44.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.70.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.70.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.46.0 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.46.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
	google.golang.org/genproto v0.0.0-20260519071638-aa98bba5eb94 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260819154853-08b0e4226688 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260819154853-08b0e4226688 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.3 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	istio.io/gogo-genproto v0.0.0-20190930162913-45029607206a // indirect
	k8s.io/apiextensions-apiserver v0.24.2 // indirect
	k8s.io/klog/v2 v2.60.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)
