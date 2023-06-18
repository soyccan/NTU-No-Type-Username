# syntax=docker/dockerfile:1
FROM golang:1.20-alpine3.18

ENV CRED_PATH=credentials.json COOKIE_PATH=/shared/cookies.conf SOCK_PATH=/shared/sock

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
  apk add socat

# check if login succeeds
HEALTHCHECK --interval=5s --timeout=5s --retries=5 CMD test -s "$COOKIE_PATH"

WORKDIR /ntu
COPY *.go go.mod ./
RUN go install -v

COPY docker-entrypoint.sh ./
CMD ["/ntu/docker-entrypoint.sh"]
