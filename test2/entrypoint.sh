#!/bin/bash
echo "entrypoint script start"

/usr/bin/telemetry_supervisor.sh &

exec "$@"

