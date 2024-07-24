#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"

docker compose -p idx_test -f "$scriptDir/pg-dc.yaml" "$@"