#!/bin/sh

set -x
set -e

docker-compose -f deployments/docker-compose.yml down -v
