#!/bin/sh

set -e
set -x

docker-compose -p test -f deployments/docker-compose.test.yml up -d camforchat-db
docker-compose -p test -f deployments/docker-compose.test.yml run --rm camforchat-sqitch /scripts/check_pg.sh
docker-compose -p test -f deployments/docker-compose.test.yml run --rm camforchat-sqitch sqitch deploy test
