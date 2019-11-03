#!/bin/sh

set -x
set -e

docker-compose -f deployments/docker-compose.yml run --rm camforchat-sqitch sqitch deploy default
