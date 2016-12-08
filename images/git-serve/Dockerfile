FROM ubuntu:16.04

RUN \
  apt-get update && \
  apt-get install -y git-core nginx fcgiwrap cron logrotate

ADD git.conf /etc/nginx/sites-available/
RUN \
  rm -rf /etc/nginx/sites-enabled/* && \
  ln -fs /etc/nginx/sites-available/git.conf /etc/nginx/sites-enabled/git.conf
ADD nginx.conf /etc/nginx/

ADD entrypoint nginx.logrotate /
ENTRYPOINT ["/entrypoint"]
