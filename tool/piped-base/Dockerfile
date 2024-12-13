FROM alpine:3.21.0

ARG PIPED_USER=piped
ARG PIPED_USER_GROUP=piped
ARG PIPED_UID=1000
ARG PIPED_GID=1000

ENV HOME=/home/${PIPED_USER}
ENV PIPED_TOOLS_DIR="${HOME}/.piped/tools"
ENV PATH="${PIPED_TOOLS_DIR}:${PATH}"

RUN \
    addgroup -S -g $PIPED_GID $PIPED_USER_GROUP && \
    adduser -S -u $PIPED_UID -G $PIPED_USER_GROUP -h $HOME $PIPED_USER && \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        jq \
        curl && \
    update-ca-certificates && \
    mkdir -p ${PIPED_TOOLS_DIR} && \
    rm -f /var/cache/apk/* && \
    chown -R $PIPED_USER:$PIPED_USER_GROUP $HOME && \
    chmod 770 -R $HOME

USER $PIPED_USER
