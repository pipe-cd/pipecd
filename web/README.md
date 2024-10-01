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

If the API server has authorization by cookie, set `API_COOKIE` to the cookie you have already obtained through other clients
(typically you need to send some kind of request from an authenticated client and peek at the request header in some way).

```
API_COOKIE={COOKIE}
```

Take Chrome for example;
1. Access to the existing UI.
2. Open the developer tools and go to the network panel.
3. Find the `GetMe` request and select it.
4. Copy the whole value of the `Cookie` in "Request Headers" and paste it to `API_COOKIE={COOKIE}` in the `.env` file.

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/play-environment-get-me.png)

TIP: If you don't want to step up (or don't have) a PipeCD controlplane API server, you can login to [https://play.pipecd.dev](https://play.pipecd.dev/login?project=play) and use its API with your authenticated account.
