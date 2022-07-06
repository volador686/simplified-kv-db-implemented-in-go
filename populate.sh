#!/bin/bash

for shard in localhost:8080; do
    echo $shard
    for i in {1..10}; do
        curl "http://${shard}/set?key=${RANDOM}&value=${RANDOM}"
    done
done