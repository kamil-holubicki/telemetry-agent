#!/bin/bash
echo "entrypoint script start"

#/usr/bin/script-1.sh &

supervisord -n &

exec "$@"

