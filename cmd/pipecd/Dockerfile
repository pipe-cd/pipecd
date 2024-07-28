FROM alpine:3.20

ARG GOOGLE_CLOUD_SDK_VERSION=477.0.0

RUN \
    apk add --no-cache \
        python3 \
        curl && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar -zxvf ./google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm ./google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    ./google-cloud-sdk/install.sh --quiet

ENV PATH="/google-cloud-sdk/bin:${PATH}"

RUN mkdir -p web/static

ADD .artifacts/pipecd /usr/local/bin/pipecd
ADD .artifacts/web-static web/static

ENTRYPOINT ["pipecd"]
