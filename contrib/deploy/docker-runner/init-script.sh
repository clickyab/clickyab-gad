#!/usr/bin/dumb-init /bin/bash
set -euo pipefail

MYSQL_USER=${MY_USER:-clickyab_master}
MYSQL_PASSWORD=${MY_PASS:-oRgLsGydHQnZzvdNfHM6}
MYSQL_DB=${MY_DB:-clickyab}
MYSQL_HOST=${MY_HOST:-192.168.100.11}
MYSQL_PORT=${MY_PORT:-3306}
REDIS_HOST=${REDIS_HOST:-192.168.100.30}
REDIS_PORT=${REDIS_PORT:-2222}
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
export GAD_REDIS_SIZE=${REDIS_SIZE:-250}
export GAD_MYSQL_MAX_CONNECTION=${MYSQL_MAX_CONNECTION:-250}
export GAD_MYSQL_MAX_CONNECTION=${MYSQL_MAX_IDLE_CONNECTION:-30}
export GAD_MYSQL_RDSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?parseTime=true&charset=utf8"
export GAD_MYSQL_WDSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?parseTime=true&charset=utf8"
export GAD_MYSQL_DATABASE="${MYSQL_DB}"
export GAD_PROFILE=disable
export GAD_AMQP_DSN="amqp://${RABBIT_USER}:${RABBIT_PASS}@${RABBIT_HOST}:${RABBIT_PORT}/"
export GAD_SLACK_ACTIVE=true
export GAD_PHP_CODE_ROOT=/app/clickyab-server/a/
export GAD_PHP_CODE_FPM=127.0.0.1:9000

if [ "$1" = '/app/bin/server' ];
then

cat >/app/clickyab-server/library/db.php <<-EOCONF
<?php

function client_addr_2(){
    if(\$_SERVER['HTTP_CF_CONNECTING_IP']) return \$_SERVER['HTTP_CF_CONNECTING_IP'];
    return str_in_db(\$_SERVER['REMOTE_ADDR']);
}

\$mysql_connect = mysqli_connect ( "${MYSQL_HOST}", "${MYSQL_USER}", "${MYSQL_PASSWORD}", "${MYSQL_DB}" );
mysqli_set_charset(\$mysql_connect,'utf8');

define("REDIS_HOST", "${REDIS_HOST}");
define("REDIS_PORT", 2222); // this is different from gad for now
define("REDIS_PASS", "${GAD_REDIS_PASSWORD}");

EOCONF
	chown www-data:www-data /app/clickyab-server/library/db.php
	service php7.0-fpm start
	exec "$@"
else
	exec "$@"
fi;