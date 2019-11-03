#!/bin/sh

set -e
set -x

# waiting for postgres up
while ! psql -h$POSTGRES_HOST -U$POSTGRES_USER -d$POSTGRES_DB -c "select 1" > /dev/null
do
    echo "$(date) - still trying to connect to PostgreSQL host"
    sleep 1
done
echo "$(date) - connected successfully"
