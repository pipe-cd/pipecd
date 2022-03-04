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
)

for dir in ${goProtoDirs[*]}; do
  echo ""
  echo "- ${dir}"
  echo "deleting previously generated Go files..."
  find ${dir} -name "*.pb.go" -o -name "*.pb.validate.go" -type f -delete
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
    ${dir}/*.proto
  echo "successfully generated"
done

# Generate JS files from proto files.
jsProtoDirs=(
  "pkg/model"
  "pkg/app/server/service/webservice"
)
jsOutDirs=(
  "pkg/app/web/model"
  "pkg/app/web/api_client"
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
    --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:. \
    ${inDir}/*.proto
  mv ${inDir}/*.js ${outDir}
  mv ${inDir}/*.d.ts ${outDir}

  find ${outDir} -type f -exec sed -i 's:.*validate_pb.*::g' {} \;
  find ${outDir} -type f -exec sed -i "s:'.*pkg:'pipecd\/pkg\/app\/web:g;" {} \;
  find ${outDir} -type f -exec sed -i "s:'.*\/model\/:'pipecd\/pkg\/app\/web\/model\/:g;" {} \;
  echo "successfully generated"

echo ""
echo "Successfully generated all code"
echo ""

done
