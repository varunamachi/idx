#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/..")

envFile="${scriptDir}/common.env"
if [[ -f  "${envFile}" ]]; then
    set -o allexport
    # shellcheck disable=SC1090
    source "${envFile}"
    set +o allexport
fi

echo "Building...."
bash "${root}/_scripts/build-dev.sh" || exit 2


export MK_USER_TOKEN="73a39932-f368-11ec-9e7b-3f1b9f8c1041"
echo "Running...."
echo
"$root/_local/bin/idx" "$@"