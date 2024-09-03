#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"

export PG_USER="idx"
export PG_PASSWORD="idxp"
export PG_DB="idx-test"
export PG_PORT="5432"

docker compose -p idx_test -f "$scriptDir/pg-dc.yaml" "$@"