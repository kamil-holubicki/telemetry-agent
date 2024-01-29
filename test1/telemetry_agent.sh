#!/bin/bash

for i in {1..8}
do
   echo "Telemetry agent $i"
   sleep 2
done

echo "telemetry agent dies..."

exit 1
