FROM nginxinc/nginx-unprivileged:1.25-alpine3.17

EXPOSE 80
EXPOSE 443
ENV COOKIE_PATH=

COPY templates /etc/nginx/
COPY docker-entrypoint.d/ /docker-entrypoint.d/

CMD ["nginx-debug", "-g", "daemon off;"]
