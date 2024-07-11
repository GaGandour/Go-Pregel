#!/bin/bash

IMAGE_NAME_OR_ID="pregel"

# Get the list of container IDs based on the specified image
CONTAINER_IDS=$(docker ps --filter "ancestor=$IMAGE_NAME_OR_ID" -q)

# Loop through each container ID and stop it
for CONTAINER_ID in $CONTAINER_IDS; do
    docker stop $CONTAINER_ID
done
