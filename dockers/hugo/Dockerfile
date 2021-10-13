FROM alpine:3.10

ARG VERSION=0.88.1
ARG PACKAGE=hugo_extended_${VERSION}_Linux-64bit.tar.gz
ARG URL=https://github.com/gohugoio/hugo/releases/download/v${VERSION}/${PACKAGE}

RUN apk add --update --no-cache --virtual .build-deps wget && \
    apk add --update --no-cache ca-certificates git libc6-compat libstdc++ bash npm && \
    wget --quiet "${URL}" -P /tmp/hugo && \
    tar -xzvf "/tmp/hugo/${PACKAGE}" -C /tmp/hugo && \
    mv /tmp/hugo/hugo /usr/local/bin && \
    rm -rf /tmp/hugo && \
    apk del .build-deps
