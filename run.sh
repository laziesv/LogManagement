#!/usr/bin/env sh
set -eu
cd "$(dirname "$0")"
sh scripts/init-env.sh
docker compose up --build -d
