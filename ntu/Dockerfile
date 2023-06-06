FROM golang:1.20-alpine3.18

# owner of cookies file
ENV UID=nobody GID=nobody

ENV CRED_PATH=credentials.json COOKIE_PATH=/shared/cookies
VOLUME /shared

# check if login succeeds
HEALTHCHECK --interval=5s --timeout=5s --retries=5 CMD test -s "$COOKIE_PATH"

WORKDIR /ntu
COPY *.go go.mod docker-entrypoint.sh ./
RUN go install -v

CMD ["/ntu/docker-entrypoint.sh"]
