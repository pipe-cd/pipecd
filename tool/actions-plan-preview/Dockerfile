FROM golang:1.22.4-alpine3.20 as builder
WORKDIR /app
COPY go.mod go.sum  ./
RUN go mod download
COPY . ./
RUN go build -o /plan-preview .

FROM ghcr.io/pipe-cd/pipectl:v0.43.1
COPY --from=builder /plan-preview /
ENV PATH $PATH:/app/cmd/pipectl
RUN chmod +x /plan-preview
ENTRYPOINT ["/plan-preview"]
