FROM nginx:alpine

COPY entrypoint.sh /usr/local/bin/git-porter

RUN apk update && \
    apk add fcgiwrap \
            fcgiwrap-openrc \
            git \
            git-daemon \
            tini && \
    echo 'FCGI_CHILDREN="5"' >> /etc/conf.d/fcgiwrap && \
    chmod +x /usr/local/bin/git-porter

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/local/bin/git-porter"]
