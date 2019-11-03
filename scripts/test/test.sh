#!/bin/sh

set -x
set -e

docker-compose -f deployments/docker-compose.test.yml -p test run --rm camforchat-app go test -v -cover ./...
