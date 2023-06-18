#!/bin/sh

trap "exit" INT TERM
trap "kill 0" EXIT

# regularly login to prevent session timeout
while true; do
    /go/bin/ntu
    sleep 3600

    # signal the proxy container to reload config
    socat UNIX-CONNECT:"$SOCK_PATH" -
done
