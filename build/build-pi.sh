#!/bin/bash
ARM_VERSION=7
echo "Building ARM image"
env GOOS=linux GOARCH=arm GOARM=$ARM_VERSION go build -o podbackup
podman build -t quay.io/vadimzharov/podbackup:latest .
podman push quay.io/vadimzharov/podbackup:latest
rm podbackup