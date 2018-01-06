#!/usr/bin/env bash
set -eo pipefail

ROOT="$(readlink -f $(dirname ${BASH_SOURCE[0]})/)"
sleep 1
${ROOT}/server
