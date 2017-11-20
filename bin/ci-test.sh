#!/usr/bin/env bash
set -euo pipefail
set -x
TARGET=${1:-test}

OLD_ROOT="$(readlink -f $(dirname ${BASH_SOURCE[0]})/../)"
GOPATH="$(mktemp -d)"
ROOT="${GOPATH}/src/clickyab.com/gad"
mkdir -p ${ROOT}

cp -R ${OLD_ROOT}/* ${ROOT}
cd ${ROOT}
make -f ${ROOT}/Makefile ${TARGET} && PASSED="true"
rm -rf ${GOPATH}

if [ -z "${PASSED:-}" ];then
    exit -1
fi;