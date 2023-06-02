FROM golang:1.20-alpine3.18

# owner of cookies file
ENV UID=nobody GID=nobody

ENV CRED_PATH=credentials.json COOKIE_PATH=/shared/cookies
VOLUME /shared

WORKDIR /ntu
COPY *.go go.mod docker-entrypoint.sh ./
RUN go install

CMD ["/ntu/docker-entrypoint.sh"]
