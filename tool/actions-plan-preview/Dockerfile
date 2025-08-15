FROM golang:1.25.0-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum  ./
RUN go mod download
COPY . ./
RUN go build -o /plan-preview .

FROM ghcr.io/pipe-cd/pipectl:v0.53.0
COPY --from=builder /plan-preview /
ENV PATH=$PATH:/app/cmd/pipectl
RUN chmod +x /plan-preview
ENTRYPOINT ["/plan-preview"]
