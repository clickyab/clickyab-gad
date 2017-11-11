#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/go/bin" >> /home/develop/.zshrc

cd /home/develop/go/src/clickyab.com/gad
make -f /home/develop/go/src/clickyab.com/gad/Makefile mysql-setup
make -f /home/develop/go/src/clickyab.com/gad/Makefile rabbitmq-setup

chown -R develop:develop . /home/develop
sudo -u develop /home/develop/go/src/clickyab.com/gad/bin/provision_user.sh