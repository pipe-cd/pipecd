# check=skip=SecretsUsedInArgOrEnv
# ghcr.io/pipe-cd/piped-base:v0.58.0-193-g4253e54
FROM ghcr.io/pipe-cd/piped-base@sha256:054891ce4d2216def0b5a4067cca264e0733582700be8b0e65cc5707e151c464

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
    # Install nss_wrapper to add an random UID to "passwd" at runtime without having to directly modify /etc/passwd.
    /installer/install-nss-wrapper.sh && \
    # Remove what were used for installation.
    rm -rf /installer && \
    rm -f /var/cache/apk/* && \
    # Create the modifiable passwd file.
    grep -v -e ^default /etc/passwd > $PIPED_HOME/passwd && \
    # Grant access to the root group because all containers running on OKD belong to it.
    chown -R $PIPED_UID:$ROOT_GID $PIPED_HOME && \
    chmod 770 -R $PIPED_HOME

ENV LD_PRELOAD=/usr/local/lib/libnss_wrapper.so
ENV NSS_WRAPPER_PASSWD=$PIPED_HOME/passwd
ENV NSS_WRAPPER_GROUP=/etc/group

USER $PIPED_UID
