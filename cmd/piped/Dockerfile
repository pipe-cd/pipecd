FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=piped

# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/321463679?tag=v0.50.0-26-ga8527d2
FROM ghcr.io/pipe-cd/piped-base@sha256:9960b45a5aa822ae45ca2966056d8d2e98795b51681df25afd1fecf96360981c

COPY --from=builder /app/.artifacts/piped /usr/local/bin/piped

ENTRYPOINT ["piped"]
