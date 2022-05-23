FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y fcgiwrap \
                       git-core \
                       nginx

COPY entrypoint.sh /usr/local/bin/git-porter
ENTRYPOINT ["/usr/local/bin/git-porter"]
