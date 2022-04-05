FROM golang:1.17.6-alpine3.15

RUN apk update && apk add git

COPY . /app

RUN cd /app && \
  go build -o /gh-release . && \
  chmod +x /gh-release

ENTRYPOINT ["/gh-release"]
