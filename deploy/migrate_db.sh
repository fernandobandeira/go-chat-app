#!/bin/sh

set -e

ENTRYPOINT="/docker-entrypoint-initdb.d"

find "$ENTRYPOINT" -iname '*.sql' -exec psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "{}" \;