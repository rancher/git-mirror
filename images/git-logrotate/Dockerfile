FROM alpine:3.5

RUN addgroup -g 40000 rancher

RUN apk add --update --no-cache logrotate curl bash tini

RUN \
  curl -sSfL https://github.com/kelseyhightower/confd/releases/download/v0.16.0/confd-0.16.0-linux-amd64 -o /usr/bin/confd && \
  chmod +x /usr/bin/confd && \
  mkdir -p /etc/confd

COPY templates /etc/confd/templates/
COPY conf.d /etc/confd/conf.d/
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/sbin/crond", "-f"]
