FROM alpine:latest

COPY entrypoint.sh /usr/local/bin/git-porter

RUN apk update && \
    apk add fcgiwrap \
            fcgiwrap-openrc \
            git \
            nginx \
            tini && \
    echo 'FCGI_CHILDREN="5"' >> /etc/conf.d/fcgiwrap && \
    chmod +x /usr/local/bin/git-porter

ENTRYPOINT ["/usr/local/bin/git-porter"]
