#!/bin/sh

set -e
set -x

docker-compose -f deployments/docker-compose.development.yml run --rm camforchat-sqitch sqitch deploy default
