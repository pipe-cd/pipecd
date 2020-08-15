FROM alpine:3.10.1

ENV PIPED_BIN_DIR=/usr/local/piped
ENV PATH="${PIPED_BIN_DIR}:${PATH}"
ENV CLOUDSDK_PYTHON=python3
ENV PATH /google-cloud-sdk/bin:$PATH

ADD install-helm.sh /installer/install-helm.sh
ADD install-kubectl.sh /installer/install-kubectl.sh
ADD install-kustomize.sh /installer/install-kustomize.sh
ADD install-gcloud.sh /installer/install-gcloud.sh

RUN \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        python3 \
        py3-crcmod \
        gnupg \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir ${PIPED_BIN_DIR} && \
    # Pre-install helm.
    /installer/install-helm.sh && \
    # Pre-install kubectl.
    /installer/install-kubectl.sh && \
    # Pre-install kustomize.
    /installer/install-kustomize.sh && \
    # Pre-install gcloud.
    /installer/install-gcloud.sh && \
    # Delete installer directory.
    rm -rf /installer && \
    rm -f /var/cache/apk/*
