#!/bin/sh

# regularly reload NGINX configs so that latest cookies are used
# TODO: find a signal-based approach: let `login` container signal `proxy` container when it renews cookies
(while true; do
    sleep 600
    nginx -s reload
done &)
