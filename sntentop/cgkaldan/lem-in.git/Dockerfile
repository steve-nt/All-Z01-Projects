############################################
# 1. Build stage: compile the Go + G3N app
############################################
FROM golang:1.24-bookworm AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      libgl1-mesa-dev \
      libglfw3 \
      libglfw3-dev \
      xorg-dev \
      pkg-config \
      libopenal-dev \
      libvorbis-dev \
      libogg-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/antfarm

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /go/bin/antfarm .

############################################
# 2. Runtime stage: slim image with runtime deps
############################################
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      libgl1-mesa-glx \
      libglfw3 \
      libopenal1 \
      libvorbis0a \
      libogg0 \
      xvfb \
      xauth \
      x11vnc \
      supervisor && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      libgl1-mesa-glx \
      libglfw3 \
      xvfb \
      xauth && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/antfarm /usr/local/bin/antfarm
COPY --from=builder /go/src/antfarm/examples /usr/local/bin/examples

# Simple entrypoint
RUN echo '#!/bin/bash\n' \
    'set -e\n' \
    '# Disable audio\n' \
    'export ALSA_CONFIG_PATH=/dev/null\n' \
    'export ALSOFT_DRIVERS=null\n' \
    'export PULSE_SERVER=none\n' \
    'export SDL_AUDIODRIVER=dummy\n' \
    '# Assume DISPLAY is already set by the docker run command\n' \
    'exec ./antfarm "$@"' > /entrypoint.sh && \
    chmod +x /entrypoint.sh


WORKDIR /usr/local/bin
ENTRYPOINT ["/entrypoint.sh"]