#!/usr/bin/env bash
set -euo pipefail

if [ "$1" = 'server' ]; then
cat >>/etc/supervisor/conf.d/$1.conf <<-EOCONF
[program:$1]
command=/gad/$1
stderr_logfile = /var/log/$1-stderr.log
stdout_logfile = /var/log/$1-stdout.log
EOCONF
fi


if [ "$1" = 'clickworker' ]; then
cat >>/etc/supervisor/conf.d/$1.conf <<-EOCONF
[program:$1]
command=/gad/$1
stderr_logfile = /var/log/$1-stderr.log
stdout_logfile = /var/log/$1-stdout.log
EOCONF
fi


if [ "$1" = 'impworker' ]; then
cat >>/etc/supervisor/conf.d/$1.conf <<-EOCONF
[program:$1]
command=/gad/$1
stderr_logfile = /var/log/$1-stderr.log
stdout_logfile = /var/log/$1-stdout.log
EOCONF
fi


if [ "$1" = 'convworker' ]; then
cat >>/etc/supervisor/conf.d/$1.conf <<-EOCONF
[program:$1]
command=/gad/$1
stderr_logfile = /var/log/$1-stderr.log
stdout_logfile = /var/log/$1-stdout.log
EOCONF
fi

exec "$@"