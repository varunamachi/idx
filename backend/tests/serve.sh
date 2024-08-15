#!/usr/bin/env bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/../..")


export IDX_MAIL_PROVIDER=IDX_SIMPLE_MAIL_SERVICE_CLIENT_PROVIDER
export PG_URL="postgresql://idx:idxp@localhost:5432/test-data?sslmode=disable"

"$root/_scripts/run.sh" serve