#!/bin/sh

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/..")
libxDir=$(readlink -f "$scriptDir/../../libx")


GIT_TAG="$(git describe --tag || echo 'latest')"
GIT_HASH="$(git rev-parse --verify HEAD)"
GIT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
BUILD_TIME="$(date -Isec)"
BUILD_HOST="$(hostname)"
BUILD_USER="$(whoami)"

export GIT_TAG
export GIT_HASH
export GIT_BRANCH
export BUILD_TIME
export BUILD_HOST
export BUILD_USER

# rsync -ac "${root}/go.mod" "${root}/go.dev.mod" || exit 101
# rsync -ac "${root}/go.sum" "${root}/go.dev.sum" || exit 102
# go mod edit -replace github.com/varunamachi/libx="${libxDir}" go.dev.mod \
#     || exit 103
# go mod tidy -modfile "go.dev.mod"
# export GOMOD="go.dev.mod"

"$scriptDir/build.sh"