#!/bin/bash

IMAGE_NAME_OR_ID="pregel"

# Get the list of container IDs based on the specified image
CONTAINER_IDS=$(docker ps --filter "ancestor=$IMAGE_NAME_OR_ID" -q)

mkdir -p ../../docker_logs
rm -f ../../docker_logs/*

# Loop through each container ID and stop it
for CONTAINER_ID in $CONTAINER_IDS; do
    docker stop $CONTAINER_ID
done

for CONTAINER_ID in $CONTAINER_IDS; do
    echo "Getting logs for container $CONTAINER_ID"
    container_name=$(docker inspect --format='{{.Name}}' "$CONTAINER_ID" | sed 's/\///')
    log_file="../../docker_logs/${container_name}.log"
    docker logs "$CONTAINER_ID" > "$log_file" 2>&1
done
