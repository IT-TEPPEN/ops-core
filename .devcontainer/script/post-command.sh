#!/usr/bin/env bash

# load .env
set -a
source .devcontainer/.env
set +a

# set git config
git config --global user.name "${GIT_NAME}"
git config --global user.email "${GIT_EMAIL}"

# fix docker socket permissions
sudo chown $(whoami) /var/run/docker.sock
