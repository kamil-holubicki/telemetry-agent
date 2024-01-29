#!/bin/bash

#while :
#do
#	echo "script-2: Press [CTRL+C] to stop.."
#	sleep 2
#done


for i in {1..30}
do
   echo "mysqld $i"
   sleep 2
done

echo "mysqld dies..."
