#!/bin/sh

set -e
set -x

docker-compose -f deployments/docker-compose.development.yml up camforchat-app
