FROM golang:1.13.4-alpine3.10 AS builder
COPY main.go .
RUN go build -o /server main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates

COPY --from=builder /server ./
RUN chmod +x ./server

COPY public /public

EXPOSE 8080
ENTRYPOINT ["./server"]
