# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/103341355?tag=v0.44.0-alpha-52-ge881cac
FROM ghcr.io/pipe-cd/piped-base@sha256:18a5a6af6507c8c70c5442d19651a2000672c78fc5dd91cfd3c0af26c899c1d3

ARG PIPED_UID=1000
ARG ROOT_GID=0
ARG PIPED_HOME=$HOME

USER root

COPY install-nss-wrapper.sh /installer/install-nss-wrapper.sh

RUN \
    apk add --no-cache \
        bash \
        gcc \
        libc-dev \
        musl-nscd-dev \
        make \
        cmake && \

    # Install glibc to be used for building nss_wrapper.
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.33-r0/glibc-2.33-r0.apk && \
    apk add glibc-2.33-r0.apk && \

    # Install nss_wrapper to add an random UID to "passwd" at runtime without having to directly modify /etc/passwd.
    /installer/install-nss-wrapper.sh && \
    # Remove what were used for installation.
    rm -rf /installer && \
    rm -f /var/cache/apk/* && \
    rm /glibc-2.33-r0.apk && \
    # Create the modifiable passwd file.
    grep -v -e ^default /etc/passwd > $PIPED_HOME/passwd && \
    # Grant access to the root group because all containers running on OKD belong to it.
    chown -R $PIPED_UID:$ROOT_GID $PIPED_HOME && \
    chmod 770 -R $PIPED_HOME

ENV LD_PRELOAD=/usr/local/lib64/libnss_wrapper.so
ENV NSS_WRAPPER_PASSWD=$PIPED_HOME/passwd
ENV NSS_WRAPPER_GROUP=/etc/group

USER $PIPED_UID
