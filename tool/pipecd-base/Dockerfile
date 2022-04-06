FROM alpine:3.13

ARG GOOGLE_CLOUD_SDK_VERSION=324.0.0

RUN \
    apk add --no-cache \
        python3 \
        curl && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86.tar.gz && \
    tar -zxvf ./google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86.tar.gz && \
    ./google-cloud-sdk/install.sh --quiet

ENV PATH="/google-cloud-sdk/bin:${PATH}"
