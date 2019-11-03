#!/bin/sh

set -x
set -e

docker-compose -p test -f deployments/docker-compose.test.yml run --rm camforchat-app golint -set_exit_status ./...
