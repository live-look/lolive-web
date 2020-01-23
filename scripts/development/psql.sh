#!/bin/sh

set -e
set -x

docker-compose -f deployments/docker-compose.yml run --rm camforchat-db psql -U postgres -h camforchat-db camforchat
