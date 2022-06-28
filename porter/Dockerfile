FROM nginx:alpine

COPY entrypoint.sh /usr/local/bin/git-porter

RUN apk update && \
    apk add fcgiwrap \
            git-daemon \
            spawn-fcgi \
            tini && \
    chmod +x /usr/local/bin/git-porter

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/local/bin/git-porter"]
