#!/bin/sh

touch "${COOKIE_PATH}"
chown -R $UID:$GID "${COOKIE_PATH}"
chmod -R 600 "${COOKIE_PATH}"

# login & exit
/go/bin/ntu
