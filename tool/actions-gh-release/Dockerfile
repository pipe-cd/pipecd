FROM golang:1.26rc3-alpine3.22

RUN apk update && apk add git

COPY . /app

RUN cd /app && \
  go build -o /gh-release . && \
  chmod +x /gh-release

ENTRYPOINT ["/gh-release"]
