#!/bin/sh

set -e
set -x

docker-compose -f deployments/docker-compose.development.yml down -v
