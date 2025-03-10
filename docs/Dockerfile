FROM golang:1.24.1-alpine3.21 AS builder
COPY main.go .
RUN go build -o /server main.go

FROM alpine:3.21.3
RUN apk --no-cache add ca-certificates

COPY --from=builder /server ./
RUN chmod +x ./server

COPY public /public

EXPOSE 8080
ENTRYPOINT ["./server"]
