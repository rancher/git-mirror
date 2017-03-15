FROM alpine:3.5

RUN apk add --update --no-cache ca-certificates
ADD code/git-analyze /usr/bin/
ENTRYPOINT ["/usr/bin/git-analyze"]