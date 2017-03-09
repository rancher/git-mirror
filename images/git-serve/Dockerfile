FROM ubuntu:16.04

RUN \
  apt-get update && \
  apt-get install -y git-core nginx fcgiwrap curl

RUN \
  curl -sSfL https://github.com/kelseyhightower/confd/releases/download/v0.12.0-alpha3/confd-0.12.0-alpha3-linux-amd64 -o /usr/bin/confd && \
  chmod +x /usr/bin/confd && \
  mkdir -p /etc/confd

ADD templates /etc/confd/templates/
ADD conf.d /etc/confd/conf.d/

ADD entrypoint /
ENTRYPOINT ["/entrypoint"]
