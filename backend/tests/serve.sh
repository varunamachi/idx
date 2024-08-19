#!/usr/bin/env bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/../..")


export IDX_MAIL_PROVIDER=IDX_SIMPLE_MAIL_SERVICE_CLIENT_PROVIDER
export IDX_SIMPLE_SRV_SEND_URL="http://localhost:9999/api/v1/send"
export PG_URL="postgresql://idx:idxp@localhost:5432/idx-test?sslmode=disable"
export IDX_ROLE_MAPPING='super:Super'

"$root/_scripts/run.sh" serve