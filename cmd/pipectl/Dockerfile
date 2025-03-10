# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.23.3 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=pipectl BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

FROM alpine:3.21.3

RUN apk add --no-cache git

COPY --from=builder /app/.artifacts/pipectl /usr/local/bin/pipectl

ENTRYPOINT ["pipectl"]
