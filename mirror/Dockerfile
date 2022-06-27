FROM alpine:latest

COPY entrypoint.sh /usr/local/bin/git-mirror

RUN apk update && \
    apk add git \
            mcookie && \
    chmod +x /usr/local/bin/git-mirror

ENTRYPOINT ["/usr/local/bin/git-mirror"]
