#!/bin/bash

while true
do
  echo "    -> telemetry agent start"
  /usr/bin/telemetry_agent.sh
  echo "    -> telemetry agent died"
done
