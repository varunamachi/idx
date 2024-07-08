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

cd "${root}/backend/tests/cmd/tester" || exit 1

echo "Building...."
export CGO_ENABLED=0
go build \
    -installsuffix 'static' \
    -o "${root}/_local/bin/tester" || exit 2



echo "Running...."
echo
"$root/_local/bin/tester" "$@"

