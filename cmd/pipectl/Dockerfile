FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=pipectl

FROM alpine:3.21.0

RUN apk add --no-cache git

COPY --from=builder /app/.artifacts/pipectl /usr/local/bin/pipectl

ENTRYPOINT ["pipectl"]
