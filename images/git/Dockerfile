FROM ubuntu:16.04

RUN \
  apt-get update && \
  apt-get install -y git-core nginx fcgiwrap

# TODO: remove
RUN \
  git clone --mirror https://github.com/rancher/infra-catalog /var/git/infra-catalog.git && \
  git clone --mirror https://github.com/rancher/rancher-catalog /var/git/rancher-catalog.git && \
  git clone --mirror https://github.com/rancher/community-catalog /var/git/community-catalog.git

ADD git.conf /etc/nginx/sites-available/
RUN \
  rm -rf /etc/nginx/sites-enabled/* && \
  ln -fs /etc/nginx/sites-available/git.conf /etc/nginx/sites-enabled/git.conf

ADD entrypoint /
ENTRYPOINT ["/entrypoint"]
