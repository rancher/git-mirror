FROM alpine:3.5

RUN apk add --update --no-cache logrotate curl bash

RUN \
  curl -sSfL https://github.com/kelseyhightower/confd/releases/download/v0.12.0-alpha3/confd-0.12.0-alpha3-linux-amd64 -o /usr/bin/confd && \
  chmod +x /usr/bin/confd && \
  mkdir -p /etc/confd

ADD templates /etc/confd/templates/
ADD conf.d /etc/confd/conf.d/

ADD entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
CMD ["logrotate", "-f", "/rotate.cfg"]