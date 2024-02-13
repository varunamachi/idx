#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"

docker-compose -f "$scriptDir/pg.docker-compose.yml" -p "fake-data" "$@"