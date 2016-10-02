export APPNAME=server
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOPATH=$(ROOT):$(ROOT)/vendor
export WATCH?=hello
export LONGHASH=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export COMMITCOUNT=$(shell git rev-list HEAD --count| cat)
export BUILDDATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
export FLAGS="-X version.hash=$(LONGHASH) -X version.short=$(SHORTHASH) -X version.date=$(COMMITDATE) -X version.count=$(COMMITCOUNT) -X version.build=$(BUILDDATE)"
export LDARG=-ldflags $(FLAGS)
export BUILD=$(BIN)/gb build $(LDARG)
export DBPASS?=$(DEFAULT_PASS)
export DB_USER?=root
export DB_NAME?=clickyab
export RUSER?=$(APPNAME)
export RPASS?=$(DEFAULT_PASS)
export WORK_DIR=$(ROOT)/tmp

.PHONY: all gb clean

all: $(GB)
	$(BUILD)

needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

notroot :
	@[ "$(shell id -u)" != "0" ] || exit 1

gb: notroot
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get -u -v github.com/constabulary/gb/...

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin

$(GB): notroot
	@[ -f $(BIN)/gb ] || make gb

server: $(GB)
	$(BUILD) server

run-server: server
	sudo setcap cap_net_bind_service=+ep $(BIN)/server
	$(BIN)/server


mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql
	echo 'FLUSH PRIVILEGES;' | mysql mysql
	echo 'CREATE DATABASE $(DB_NAME);' | mysql -u $(DB_USER) -p$(DBPASS)