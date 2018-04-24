#!/bin/bash
set -x
set -eo pipefail

# This job is from jenkins. so kill it if it is a pull request
exit_message() {
    echo ${1:-'exiting...'}
    code=${2:-1}
    if [[ "#{code}" == "0" ]]; then
        echo "${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
        echo "Build was OK, but it's not the correct branch(${APP}:${BRANCH}.${COMMIT_COUNT} By ${CHANGE_AUTHOR}). ignore this." >> ${OUT_LOG}
        echo "green" >> ${OUT_LOG_COLOR}
    else
        echo "${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
        echo "Build was NOT OK (${APP}:${BRANCH}.${COMMIT_COUNT} By ${CHANGE_AUTHOR}). Verify with dev team." >> ${OUT_LOG}
        echo "red" > ${OUT_LOG_COLOR}
    fi
    exit ${code}
}

APP=${APP:-}
BRANCH=${BRANCH_NAME:-master}
BRANCH=${CHANGE_TARGET:-${BRANCH}}
CACHE_ROOT=${CACHE_ROOT:-/var/lib/jenkins/cache}

[ -z ${CHANGE_AUTHOR} ] && PUSH="--push"
[ -z ${APP} ] && exit_message "The APP is not defined." # WTF, the APP_NAME is important


SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

SOURCE_DIR=${1:-}
[ -z ${SOURCE_DIR} ] && exit_message "Must pass the source directory as the first parameter" 1
SOURCE_DIR=$(cd "${SOURCE_DIR}/" && pwd)

BUILD_DIR=${2:-${SOURCE_DIR}-build}
CACHE_DIR=${CACHE_ROOT}/${APP}-${BRANCH}
ENV_DIR=$(mktemp -d)

mkdir -p "${BUILD_DIR}" "${CACHE_DIR}" "${ENV_DIR}"
BUILD=$(cd "${BUILD_DIR}/" && pwd)
CACHE=$(cd "${CACHE_DIR}/" && pwd)
VARS=$(cd "${ENV_DIR}/" && pwd)

#chown $(id -u):$(id -g) $ENV_DIR
#chown $(id -u):$(id -g) $CACHE_DIR

BUILD_PACKS_DIR=$(mktemp -d)

#GIT_DIR=$(cd "${SOURCE_DIR}/" && pwd)
pushd ${SOURCE_DIR}
GIT_WORK_TREE=${BUILD} git checkout -f HEAD

export LONGHASH=$(git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMPDATE=$(date +%Y%m%d)
export COMMITCOUNT=$(git rev-list HEAD --count| cat)
export BUILDDATE=$(date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
popd

env -0 | while IFS='=' read -r -d '' n v; do
    echo "${v}">"${VARS}/${n}";
done< <(env -0)

TEMPORARY=$(mktemp -d)

cat > ${TEMPORARY}/Rockerfile <<EOF
FROM gliderlabs/herokuish

MOUNT {{ .Build }}:/tmp/app
MOUNT {{ .EnvDir }}:/tmp/env
MOUNT {{ .Target }}:/tmp/build
MOUNT {{ .Cache }}:/tmp/cache

RUN /bin/herokuish buildpack build && rm -rf /app/pkg && rm -rf /app/tmp
EXPORT /app/bin /app

FROM ubuntu:16.04
IMPORT /app

RUN apt-get update && apt-get install -y ca-certificates tzdata && apt-get clean

ENV TZ=Asia/Tehran
RUN ln -snf /usr/share/zoneinfo/\$TZ /etc/localtime && echo \$TZ > /etc/timezone

TAG registry.clickyab.ae/clickyab/{{ .App }}:{{ .Version }}
PUSH registry.clickyab.ae/clickyab/{{ .App }}:{{ .Version }}
TAG registry.clickyab.ae/clickyab/{{ .App }}:latest
PUSH registry.clickyab.ae/clickyab/{{ .App }}:latest
EOF

pushd ${TEMPORARY}
TARGET=$(mktemp -d)
rocker build ${PUSH} -var Build=${BUILD} -var EnvDir=${VARS} -var Cache=${CACHE} -var Target=${TARGET} -var Version=${BRANCH}.${COMMITCOUNT} -var App=${APP}

popd

echo "${VARS}" >> /tmp/kill-me
echo "${TARGET}" >> /tmp/kill-me
echo "${TEMPORARY}" >> /tmp/kill-me
echo "${BUILD_DIR}" >> /tmp/kill-me
echo "${BUILD_PACKS_DIR}" >> /tmp/kill-me

[ -z ${CHANGE_AUTHOR} ] || exit_message "Build OK" 0

echo "${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
echo "The branch ${BRANCH} build finished, try to deploy it" >> ${OUT_LOG}
echo "If there is no report after this for successful deploy, it means the deply failed. report it please." >> ${OUT_LOG}
if [[  "${BRANCH}" == "master"  ]]; then
    for WRKTYP in webserver impression click
    do
        kubectl -n ${APP} set image deployment  ${APP}-${WRKTYP} ${APP}-${BRANCH}=registry.clickyab.ae/clickyab/${APP}:${BRANCH}.${COMMITCOUNT} --record
    done

elif [[  "${BRANCH}" == "dev"  ]]; then
    kubectl -n ${APP} set image deployment  ${APP}-webserver-${BRANCH} ${APP}-${BRANCH}=registry.clickyab.ae/clickyab/${APP}:${BRANCH}.${COMMITCOUNT} --record
fi
echo "..." >> ${OUT_LOG}
echo "Deploy done successfully to image registry.clickyab.ae/clickyab/${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
echo "green" >> ${OUT_LOG_COLOR}
