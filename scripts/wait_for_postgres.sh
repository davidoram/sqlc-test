#!/bin/bash
# Postges is not ready to accept connections immediately after starting, Even if you wait for port 5432 to be ready,
# it will still refuse connections. This script waits for the postgres to be ready to accept connections.
# See https://github.com/docker-library/postgres/issues/880
until PGPASSWORD=postgres psql -h localhost -U postgres -c "select 1" ; do sleep 2 ; echo "waiting for postgres ..."; done