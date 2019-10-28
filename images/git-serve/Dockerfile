FROM ubuntu:18.04

RUN \
  apt-get update && \
  apt-get install -y git-core nginx fcgiwrap curl && \
  rm -rf /var/lib/apt/lists/*

RUN \
  curl -sSfL https://github.com/kelseyhightower/confd/releases/download/v0.16.0/confd-0.16.0-linux-amd64 -o /usr/bin/confd && \
  chmod +x /usr/bin/confd && \
  mkdir -p /etc/confd

COPY templates /etc/confd/templates/
COPY conf.d /etc/confd/conf.d/

COPY entrypoint /
ENTRYPOINT ["/entrypoint"]
