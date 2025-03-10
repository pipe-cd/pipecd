FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum  ./
RUN go mod download
COPY . ./
RUN go build -o /plan-preview .

FROM ghcr.io/pipe-cd/pipectl:v0.43.1
COPY --from=builder /plan-preview /
ENV PATH=$PATH:/app/cmd/pipectl
RUN chmod +x /plan-preview
ENTRYPOINT ["/plan-preview"]
