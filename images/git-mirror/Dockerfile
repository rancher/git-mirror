FROM golang:1.7-alpine AS builder

RUN apk add --update --no-cache git
RUN go get github.com/sirupsen/logrus
RUN go get gopkg.in/yaml.v2
RUN go get github.com/coreos/etcd/client
RUN go get github.com/gorilla/mux
RUN go get golang.org/x/net/context
COPY ./code /usr/src/github.com/rancher/git-mirror
WORKDIR /usr/src/github.com/rancher/git-mirror
RUN go build

FROM alpine:3.4

RUN apk add --update --no-cache git
COPY --from=builder /usr/src/github.com/rancher/git-mirror/git-mirror /usr/bin/
COPY config.yaml /etc/git-mirror/
ENTRYPOINT ["/usr/bin/git-mirror"]
# CMD ["-config-file", "/etc/git-mirror/config.yaml"]
