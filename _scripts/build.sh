#!/bin/sh

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/..")


cmdDir="${root}/backend/cmd/idx"
if [ ! -d "$cmdDir" ] ; then
    echo "Command directory $cmdDir does not exist"
fi
cd "$cmdDir" || exit 1


git --version  >/dev/null 2>&1
GIT_IS_AVAILABLE=$?
if [ $GIT_IS_AVAILABLE -eq 0 ] &&  [ -z "$GIT_TAG" ]; then 
    GIT_TAG=$(git describe --tag || echo 'latest')
    GIT_HASH=$(git rev-parse --verify HEAD)
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    BUILD_TIME=$(date -Isec)
    BUILD_HOST=$(hostname)
    BUILD_USER=$(whoami)
fi

if [ -z "$GIT_TAG" ];    then echo "GIT_TAG not set "; exit 1; fi
if [ -z "$GIT_HASH" ];   then echo "GIT_HASH not set "; exit 1; fi
if [ -z "$GIT_BRANCH" ]; then echo "GIT_BRANCH not set "; exit 1; fi
if [ -z "$BUILD_TIME" ]; then echo "BUILD_TIME not set "; exit 1; fi
if [ -z "$BUILD_HOST" ]; then echo "BUILD_HOST not set "; exit 1; fi
if [ -z "$BUILD_USER" ]; then echo "BUILD_USER not set "; exit 1; fi


# GOMOD=${GOMOD:-"go.mod"}
# echo "Using Go mod file: ${GOMOD}"

depDir="${root}/_local/bin"
if [ ! -d "${depDir}" ]; then 
    mkdir -p "${depDir}" || exit 1
fi

export CGO_ENABLED=0
go build \
    -installsuffix 'static' \
    -ldflags "-w -s \
        -X github.com/varunamachi/idx/core.GitTag=${GIT_TAG}
        -X github.com/varunamachi/idx/core.GitHash=${GIT_HASH}
        -X github.com/varunamachi/idx/core.GitBranch=${GIT_BRANCH}
        -X github.com/varunamachi/idx/core.BuildTime=${BUILD_TIME}
        -X github.com/varunamachi/idx/core.BuildHost=${BUILD_HOST}
        -X github.com/varunamachi/idx/core.BuildUser=${BUILD_USER}    
    "\
    -o "${depDir}/idx" || exit 2



