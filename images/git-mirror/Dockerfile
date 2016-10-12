FROM alpine:3.4

RUN apk add --update --no-cache git
ADD git-mirror /usr/bin/
ADD config.yaml /etc/git-mirror/
ENTRYPOINT ["/usr/bin/git-mirror"]
CMD ["-config-file", "/etc/git-mirror/config.yaml"]