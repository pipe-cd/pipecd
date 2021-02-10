# PipeCD web

## Setup

```bash
bazelisk build //pkg/app/web:build_api //pkg/app/web:build_model # generate models and API client from proto files. Also will install dependencies by yarn
```

## Development

### Running with Mocks(msw)

We use `msw` for mocking API, so you can see UI without running API server.

```bash
yarn dev
```

The app will be available at http://localhost:9090.

### Connect Real API server

```bash
cp .env.example .env
```

Add your API endpoint to the env file like this:

```
API_ENDPOINT=https://api.pipecd.dev
```

If API server has authorization by cookie, you can use `API_COOKIE` for adding cookie to request.
