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
├── mocks # API mock files
│  └── services
├── queries # react query modules
│  └── module-name
│     ├── query-hook.ts # use query hook
│     └── mutation-hook.ts # use mutation hook
├── styles # shared styles
├── types # application types
└── utils # utils files
```

## Development

### Prerequisites

- [NodeJS v20 or later](https://nodejs.org/en/)
- [Yarn](https://yarnpkg.com/)

### Running with Mocks(msw)

First time running, you need to install dependencies.

```bash
make update/web-deps
```

We use `msw` for mocking API, so you can see UI without running API server.

```bash
make run/web
```

The app will be available at http://localhost:9090.

### Connect Real API server
If you want to connect to a real API server, additional settings on the `.env` file are needed.

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

For local development, you can set API_ENDPOINT=http://localhost:8080 after running local server following [here](../CONTRIBUTING.md)

TIP: If you don't want to step up (or don't have) a PipeCD controlplane API server, you can use [https://play.pipecd.dev](https://play.pipecd.dev) as API_ENDPOINT, and interact with `play` project with your authenticated account.
