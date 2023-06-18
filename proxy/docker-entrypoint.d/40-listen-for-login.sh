#!/bin/sh

if [ -S "$SOCK_PATH" ]; then
    rm "$SOCK_PATH"
fi

# reload NGINX configs after a periodic login is done
# (signaled by a unix socket connection),
# so that latest cookies are used
socat -d -d -lf /dev/stderr UNIX-LISTEN:"$SOCK_PATH" EXEC:'nginx -s reload' &
