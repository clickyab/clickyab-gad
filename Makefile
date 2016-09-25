#@IgnoreInspection BashAddShebang
export APPNAME=gad
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOPATH=$(ROOT):$(ROOT)/vendor
export LONGHASH=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(shell git log -n1 --pretty="format:%cd"| sed -e "s/ /-/g")
export COMMITCOUNT=$(shell git rev-list HEAD --count| cat)
export BUILDDATE=$(shell date| sed -e "s/ /-/g")
export FLAGS="-X shared/config.hash=$(LONGHASH) -X shared/config.short=$(SHORTHASH) -X shared/config.date=$(COMMITDATE) -X shared/config.count=$(COMMITCOUNT) -X shared/config.build=$(BUILDDATE)"
export LDARG=-ldflags $(FLAGS)
export GB=$(BIN)/gb
export BUILD=$(GB) build $(LDARG)
export DBPASS?=$(DEFAULT_PASS)
export DUSER?=$(APPNAME)
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
	$(BIN)/gb build $(LDARG) server

run-server: server
	sudo setcap cap_net_bind_service=+ep $(BIN)/server
	$(BIN)/server