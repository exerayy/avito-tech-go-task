#!/bin/sh

TIMEOUT=30


n=0
while [ $n -le $TIMEOUT ]
do
    pg_isready -h localhost > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "pg ready"
        break
    fi
    n=$(( n+1 ))
    echo "pg not ready. retry $n" 
    sleep 1
done
