#!/usr/bin/dumb-init /bin/bash
set -euo pipefail

MYSQL_USER=${MY_USER:-clickyab_master}
MYSQL_PASSWORD=${MY_PASS:-oRgLsGydHQnZzvdNfHM6}
MYSQL_DB=${MY_DB:-clickyab}
MYSQL_HOST=${MY_HOST:-192.168.100.11}
MYSQL_PORT=${MY_PORT:-3306}
REDIS_HOST=${REDIS_HOST:-192.168.100.30}
REDIS_PORT=${REDIS_PORT:-2223}
RABBIT_PASS=${AMQP_PASS:-eeTheej2}
RABBIT_USER=${AMQP_USER:-cy}
RABBIT_HOST=${AMQP_HOST:-192.168.100.30}
RABBIT_PORT=${AMQP_POST:-5672}

# TODO : env re-write must be done here
export GAD_DEVEL_MODE=false
export GAD_CLICKYAB_UNDER_FLOOR=true
export GAD_SITE=a.clickyab.com
export GAD_PROTO=http
export GAD_REDIS_ADDRESS=${REDIS_HOST}:${REDIS_PORT}
export GAD_REDIS_PASSWORD=${REDIS_PASSWORD:-bita123}
export GAD_REDIS_DATABASE=${REDIS_DATABASE:-7}
export GAD_MYSQL_RDSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?parseTime=true&charset=utf8"
export GAD_MYSQL_WDSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?parseTime=true&charset=utf8"
export GAD_MYSQL_DATABASE="${MYSQL_DB}"
export GAD_PROFILE=disable
export GAD_AMQP_DSN="amqp://${RABBIT_USER}:${RABBIT_PASS}@${RABBIT_HOST}:${RABBIT_PORT}/"
export GAD_SLACK_ACTIVE=true

exec "$@"