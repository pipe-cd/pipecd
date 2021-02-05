FROM python:alpine3.13

RUN \
    apk add --no-cache \
        curl && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-324.0.0-linux-x86.tar.gz && \
    tar -zxvf ./google-cloud-sdk-324.0.0-linux-x86.tar.gz && \
    ./google-cloud-sdk/install.sh --quiet

ENV PATH="/google-cloud-sdk/bin:${PATH}"
