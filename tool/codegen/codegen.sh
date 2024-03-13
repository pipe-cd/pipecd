#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

while [[ -z "$1" ]]
do
    echo "Missing path to the root of repository"
    exit 1
done

cd $1

# Generate Go files from proto files.
goProtoDirs=(
  "pkg/model"
  "pkg/app/server/service/apiservice"
  "pkg/app/server/service/pipedservice"
  "pkg/app/server/service/webservice"
  "pkg/app/helloworld/service"
  "pkg/plugin/platform/api"
)

for dir in ${goProtoDirs[*]}; do
  echo ""
  echo "- ${dir}"
  echo "deleting previously generated Go files..."
  find ${dir} -name "*.pb.go" -o -name "*.pb.*.go" -type f -delete
  echo "successfully deleted"

  echo "generating new Go files..."
  protoc \
    -I . \
    -I /go/src/github.com/envoyproxy/protoc-gen-validate \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --validate_out="lang=go:." \
    --validate_opt=paths=source_relative \
    --auth_out=. \
    --auth_opt=paths=source_relative \
    ${dir}/*.proto
  echo "successfully generated"
done

# Generate JS files from proto files.
jsProtoDirs=(
  "pkg/model"
  "pkg/app/server/service/webservice"
)
jsOutDirs=(
  "web/model"
  "web/api_client"
)

i=0
while [ $i -lt ${#jsProtoDirs[*]} ]; do
  inDir=${jsProtoDirs[$i]}
  outDir=${jsOutDirs[$i]}
  i=$(( $i + 1))

  echo ""
  echo "- ${inDir}"
  echo "deleting previously generated JS files..."
  rm -rf ${outDir}
  mkdir -p ${outDir}
  echo "successfully deleted"

  echo "generating new JS files..."
  protoc \
    -I . \
    -I /go/src/github.com/envoyproxy/protoc-gen-validate \
    --js_out=import_style=commonjs:. \
    --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:. \
    ${inDir}/*.proto
  mv ${inDir}/*.js ${outDir}
  mv ${inDir}/*.ts ${outDir}

  find ${outDir} -type f -exec sed -i 's:.*validate_pb.*::g' {} \;
  find ${outDir} -type f -exec sed -i "s:'.*pkg:'pipecd\/web:g;" {} \;
  find ${outDir} -type f -exec sed -i "s:'.*\/model\/:'pipecd\/web\/model\/:g;" {} \;
  echo "successfully generated"
done


# Generate Go mock.
# TODO: This is just a temporary solution. We will move to have .codegen.yaml config file instead of hard coding like this.
mockPackageNames=(
  "redistest"
  "datastoretest"
  "filestoretest"
  "kubernetestest"
  "cachetest"
  "gittest"
  "jwttest"
  "insightstoretest"
)
mockDestinations=(
  "pkg/redis/redistest/redis.mock.go"
  "pkg/datastore/datastoretest/datastore.mock.go"
  "pkg/filestore/filestoretest/filestore.mock.go"
  "pkg/app/piped/platformprovider/kubernetes/kubernetestest/kubernetes.mock.go"
  "pkg/cache/cachetest/cache.mock.go"
  "pkg/git/gittest/git.mock.go"
  "pkg/jwt/jwttest/jwt.mock.go"
  "pkg/insight/insighttest/insight.mock.go"
)
mockSources=(
  "github.com/pipe-cd/pipecd/pkg/redis"
  "github.com/pipe-cd/pipecd/pkg/datastore"
  "github.com/pipe-cd/pipecd/pkg/filestore"
  "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
  "github.com/pipe-cd/pipecd/pkg/cache"
  "github.com/pipe-cd/pipecd/pkg/git"
  "github.com/pipe-cd/pipecd/pkg/jwt"
  "github.com/pipe-cd/pipecd/pkg/insight"
)
mockInterfaces=(
  "Redis"
  "ProjectStore,PipedStore,ApplicationStore,DeploymentStore,CommandStore"
  "Store"
  "Applier,Loader"
  "Getter,Putter,Deleter,Cache"
  "Repo"
  "Signer,Verifier"
  "Store"
)

i=0
while [ $i -lt ${#mockPackageNames[*]} ]; do
  package=${mockPackageNames[$i]}
  destination=${mockDestinations[$i]}
  source=${mockSources[$i]}
  interfaces=${mockInterfaces[$i]}
  i=$(( $i + 1))

  echo ""
  echo "- ${destination}"
  echo "generating mock..."
  mockgen --build_flags=--mod=mod --package=${package} --destination=${destination} ${source} ${interfaces}
  echo "successfully generated"
done

echo ""
echo "Successfully generated all code"
echo ""
