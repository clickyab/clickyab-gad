#!/bin/bash -x
set -euo pipefail


echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/gad/bin" >> /home/develop/.zshrc

cd /home/develop/go/src/clickyab.com/gad

make all
