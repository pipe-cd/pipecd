FROM ubuntu:18.04

ARG BAZELISK_VERSION=1.5.0

RUN \
    # This makes add-apt-repository available
    apt-get update && \
    apt-get -y install \
        make \
        python \
        python-pkg-resources \
        software-properties-common \
        unzip \
        apt-transport-https \
        curl \
        liblz4-tool \
        build-essential \
        ca-certificates && \
    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \
    # Install bazelisk
    curl -sSL -o /usr/bin/bazelisk https://github.com/bazelbuild/bazelisk/releases/download/v${BAZELISK_VERSION}/bazelisk-linux-amd64 && \
    chmod +x /usr/bin/bazelisk
