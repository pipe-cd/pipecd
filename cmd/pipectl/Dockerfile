FROM alpine:3.21.0

RUN apk add --no-cache git

ADD .artifacts/pipectl /usr/local/bin/pipectl

ENTRYPOINT ["pipectl"]
