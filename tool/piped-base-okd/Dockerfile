# check=skip=SecretsUsedInArgOrEnv
# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/459025001?tag=v0.52.0-97-g3e3df6d
FROM ghcr.io/pipe-cd/piped-base@sha256:fc2bde38211a80ee0df6d90ed168c817f30a7d0dbe8a0aff872a93ffb7d40a66

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
        jq \
        cmake && \
    # Install glibc to be used for building nss_wrapper.
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.35-r1/glibc-2.35-r1.apk && \
    apk add glibc-2.35-r1.apk && \
    # Install nss_wrapper to add an random UID to "passwd" at runtime without having to directly modify /etc/passwd.
    /installer/install-nss-wrapper.sh && \
    # Remove what were used for installation.
    rm -rf /installer && \
    rm -f /var/cache/apk/* && \
    rm /glibc-2.35-r1.apk && \
    # Create the modifiable passwd file.
    grep -v -e ^default /etc/passwd > $PIPED_HOME/passwd && \
    # Grant access to the root group because all containers running on OKD belong to it.
    chown -R $PIPED_UID:$ROOT_GID $PIPED_HOME && \
    chmod 770 -R $PIPED_HOME

ENV LD_PRELOAD=/usr/local/lib64/libnss_wrapper.so
ENV NSS_WRAPPER_PASSWD=$PIPED_HOME/passwd
ENV NSS_WRAPPER_GROUP=/etc/group

USER $PIPED_UID
