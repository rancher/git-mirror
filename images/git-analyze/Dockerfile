FROM golang:1.7-alpine AS builder

RUN apk add --update --no-cache git
RUN go get github.com/sirupsen/logrus
COPY ./code /usr/src/github.com/rancher/git-analyze
WORKDIR /usr/src/github.com/rancher/git-analyze
RUN go build


FROM alpine:3.5

RUN apk add --update --no-cache ca-certificates
COPY --from=builder /usr/src/github.com/rancher/git-analyze/git-analyze /usr/bin/
ENTRYPOINT ["/usr/bin/git-analyze"]
