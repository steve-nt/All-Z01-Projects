#!/usr/bin/env bash

podman build -f Containerfile -t sphinx-image-fedora
podman run --rm --volume $PWD:/app sphinx-image-fedora
