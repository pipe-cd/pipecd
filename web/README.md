# PipeCD web

## Directory structure

```bash
src
├── __fixtures__ # dummy models
├── api # API clients
├── components # shared components
│  └── comp-name
│     ├── comp-name # component's components
│     ├── index.tsx
│     ├── index.test.ts
│     └── index.stories.ts
├── constants # shared constants
├── hooks # shared hooks
├── middlewares # redux middlewares
├── mocks # API mock files
│  └── services
├── modules # redux modules
│  └── module-name
│     ├── index.ts
│     └── index.test.ts
├── styles # shared styles
├── types # application types
└── utils # utils files
```

## Development

### Running with Mocks(msw)
We use `msw` for mocking API, so you can see UI without running API server.

```bash
make run/web
```

The app will be available at http://localhost:9090.

### Connect Real API server
If you want to connect the real API server, additional settings on the `.env` file are needed.

First, create your own `.env` file based on the `.env.example` file.

```bash
cp .env.example .env
```

Add your API endpoint to the `.env` file like this:

```
API_ENDPOINT=https://{API_ADDRESS}
```

Set `ENABLE_MOCK` to false explicitly.

```
ENABLE_MOCK=false
```

If API server has authorization by cookie, set `API_COOKIE` to the cookie you have already obtained through other clients
(typically you need to send some kind of request from an authenticated client and peek at the request header in some way).

```
API_COOKIE={COOKIE}
```
