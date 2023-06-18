# syntax=docker/dockerfile:1
# TODO: use unprivileged image: nginxinc/nginx-unprivileged
FROM nginx:1.25-alpine3.17

EXPOSE 80
EXPOSE 443
ENV COOKIE_PATH=/shared/cookies.conf SOCK_PATH=/shared/sock

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
  apk add socat && \
  mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf~

COPY templates /etc/nginx/templates/
COPY docker-entrypoint.d /docker-entrypoint.d/

CMD ["nginx-debug", "-g", "daemon off;"]
