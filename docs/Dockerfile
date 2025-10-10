FROM golang:1.25.0-alpine3.22 AS builder
COPY main.go .
RUN go build -o /server main.go

FROM alpine:3.22.1
RUN apk --no-cache add ca-certificates

COPY --from=builder /server ./
RUN chmod +x ./server

COPY public /public

EXPOSE 8080
ENTRYPOINT ["./server"]
