#!/bin/sh

trap "exit" INT TERM
trap "kill 0" EXIT

touch "${COOKIE_PATH}"
chown -R $UID:$GID "${COOKIE_PATH}"
chmod -R 600 "${COOKIE_PATH}"

# regularly login to prevent session timeout
while true; do
    /go/bin/ntu
    sleep 3600
done
