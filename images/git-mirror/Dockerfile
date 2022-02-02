FROM alpine:edge

RUN apk update && \
    apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
            git \
            yq

COPY entrypoint.sh /usr/local/bin/git-mirror
COPY config.yaml /etc/git-mirror/

RUN chmod +x /usr/local/bin/git-mirror

ENTRYPOINT ["/usr/local/bin/git-mirror"]

