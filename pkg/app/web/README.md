# web

## Development

```
yarn
bazelisk build //pkg/app/web:build_api
yarn dev
```

## Run API server

1. Run `cmd/api` server

See `cmd/api` README.

2. Setup proxy server(envoy, grpcwebproxy)

```sh
grpcwebproxy \
   --backend_addr=localhost:9081 \
   --run_tls_server=false
```

3. Copy `.env.example` and set `API_ENDPOINT`

```sh
API_ENDPOINT=http://localhost:8080 # 8080 is grpcwebproxy default port
```
