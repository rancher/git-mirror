FROM ubuntu:20.04

COPY entrypoint.sh /usr/local/bin/git-porter

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y fcgiwrap \
                       git-core \
                       nginx && \
    echo 'FCGI_CHILDREN="5"' >> /etc/default/fcgiwrap && \
    chmod +x /usr/local/bin/git-porter

ENTRYPOINT ["/usr/local/bin/git-porter"]
