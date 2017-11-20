export APPNAME=server
export DEFAULT_PASS=bita123
export READ_PASS?=
export GO=$(shell which go)
export NODE=$(shell which nodejs)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export LINTER=$(BIN)/gometalinter.v1
export GOPATH=$(abspath $(ROOT)/../../..)
export GOBIN=$(ROOT)/bin
export DIFF=$(shell which diff)
export WATCH?=hello
export LONG_HASH?=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORT_HASH?=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMIT_DATE?=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMPDATE=$(shell date +%Y%m%d)
export COMMIT_COUNT?=$(shell git rev-list HEAD --count| cat)
export BUILD_DATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
export FLAGS="-X version.hash=$(LONG_HASH) -X version.short=$(SHORT_HASH) -X version.date=$(COMMIT_DATE) -X version.count=$(COMMIT_COUNT) -X version.build=$(BUILD_DATE)"
export LD_ARGS=-ldflags $(FLAGS)
export BUILD=cd $(ROOT) && $(GO) install -v $(LD_ARGS)
export DBPASS?=$(DEFAULT_PASS)
export DB_USER?=root
export DB_NAME?=clickyab
export RUSER?=$(APPNAME)
export RPASS?=$(DEFAULT_PASS)
export WORK_DIR=$(ROOT)/tmp
export LINTERCMD=$(LINTER) -e "vendor/.*" -e "tmp/.*" -e ".*gen.go" --cyclo-over=27 --line-length=180 --deadline=100s --disable-all --enable=structcheck --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=gosimple --enable=staticcheck --enable=unused
export UGLIFYJS=$(ROOT)/node_modules/.bin/uglifyjs
export GAD_SERVICES_MYSQL_WDSN=$(DB_USER):$(DBPASS)@tcp(127.0.0.1:3306)/$(DB_NAME)?charset=utf8&parseTime=true
export GAD_SERVICES_MYSQL_RDSN=dev:$(READ_PASS)@tcp(db-1.clickyab.ae:3306)/$(DB_NAME)?charset=utf8&parseTime=true

.PHONY: all clean

all:
	$(BUILD) clickyab.com/gad/cmd/...

needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

metalinter:
	GOPATH=$(ROOT)/tmp $(GO)  get -v gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(ROOT)/tmp $(LINTER) --install

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin

$(LINTER):
	@[ -f $(LINTER) ] || make -f $(ROOT)/Makefile metalinter

server: stylegen
	$(BUILD) clickyab.com/gad/cmd/server

impworker:
	$(BUILD) clickyab.com/gad/cmd/impworker

clickworker:
	$(BUILD) clickyab.com/gad/cmd/clickworker

convworker:
	$(BUILD) clickyab.com/gad/cmd/convworker

experiment:
	$(BUILD) clickyab.com/gad/cmd/experiment

run-server: server
	sudo setcap cap_net_bind_service=+ep $(BIN)/server
	$(BIN)/server

run-server-docker: restore server
	$(BIN)/server

run-impworker: impworker
	$(BIN)/impworker

run-clickworker: clickworker
	$(BIN)/clickworker

run-convworker: convworker
	$(BIN)/convworker

run-experiment: experiment
	sudo setcap cap_net_bind_service=+ep $(BIN)/experiment
	$(BIN)/experiment

mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql | true
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql | true
	echo 'FLUSH PRIVILEGES;' | mysql mysql | true
	echo 'DROP DATABASE IF EXISTS $(DB_NAME); CREATE DATABASE $(DB_NAME);' | mysql -u $(DB_USER) -p$(DBPASS)
	mysql -u $(DB_USER) -p$(DBPASS) -c $(DB_NAME) <$(ROOT)/db/structure.sql

rabbitmq-setup: needroot
	[ "1" -eq "$(shell rabbitmq-plugins enable rabbitmq_management | grep 'Plugin configuration unchanged' | wc -l)" ] || (rabbitmqctl stop_app && rabbitmqctl start_app)
	rabbitmqctl add_user $(RUSER) $(RPASS) || rabbitmqctl change_password $(RUSER) $(RPASS)
	rabbitmqctl set_user_tags $(RUSER) administrator
	rabbitmqctl set_permissions -p / $(RUSER) ".*" ".*" ".*"
	wget -O /usr/bin/rabbitmqadmin http://127.0.0.1:15672/cli/rabbitmqadmin
	chmod a+x /usr/bin/rabbitmqadmin
	rabbitmqadmin declare queue name=dlx-queue
	rabbitmqadmin declare exchange name=dlx-exchange type=topic
	rabbitmqctl set_policy DLX ".*" '{"dead-letter-exchange":"dlx-exchange"}' --apply-to queues
	rabbitmqadmin declare binding source="dlx-exchange" destination_type="queue" destination="dlx-queue" routing_key="#"

lint: $(LINTER)
	cd $(ROOT) && $(LINTERCMD) ./...

uglifyjs:
	npm install uglifyjs

$(UGLIFYJS):
	@[ -f $(UGLIFYJS) ] || make uglifyjs

uglify: $(UGLIFYJS)
	rm -rf $(ROOT)/tmp/embed
	mkdir -p $(ROOT)/tmp/embed
	cp $(ROOT)/statics/show.js $(ROOT)/tmp/embed/show.js
	$(NODE) $(UGLIFYJS) $(ROOT)/statics/show.js -o $(ROOT)/tmp/embed/show-min.js
	cp $(ROOT)/statics/conversion/clickyab-tracking.js $(ROOT)/tmp/embed/clickyab-tracking.js
	$(NODE) $(UGLIFYJS) $(ROOT)/statics/conversion/clickyab-tracking.js -o $(ROOT)/tmp/embed/clickyab-tracking-min.js
	cp $(ROOT)/statics/vastAD.js $(ROOT)/tmp/embed/vastAD.js
	$(NODE) $(UGLIFYJS) $(ROOT)/statics/vastAD.js -o $(ROOT)/tmp/embed/vastAD-min.js

go-bindata:
	$(BUILD) github.com/jteeuwen/go-bindata/go-bindata

embed: go-bindata uglify
	cd $(ROOT)/tmp/embed/ && $(BIN)/go-bindata -o $(ROOT)/statics_src/static-no-lint.gen.go -nomemcopy=true -pkg=statics ./...

create-imp-table :
	echo 'CREATE TABLE impressions$(IMPDATE)  LIKE impressions20161108; ' | mysql -u $(DB_USER) -p$(DBPASS) -c $(DB_NAME)

conditional-restore:
	$(DIFF) $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done || make restore

docker-build: conditional-restore all

ansible:
	ansible-playbook -vvvv -i $(ROOT)/contrib/deploy/hosts.ini $(ROOT)/contrib/deploy/staging.yaml

stylegen:
	GOPATH=$(ROOT)/tmp $(GO) get -v github.com/kib357/less-go/...
	$(BIN)/lessc -i $(ROOT)/routes/native_style.less -o $(ROOT)/routes/native_style.css
	echo "// Code generated by lessc. DO NOT EDIT.\n// Source: /routes/native_style.less\n\npackage routes\n\nconst style=\`" > $(ROOT)/routes/native_style.gen.go
	cat $(ROOT)/routes/native_style.css >>  $(ROOT)/routes/native_style.gen.go
	echo "\`" >> $(ROOT)/routes/native_style.gen.go
	rm $(ROOT)/routes/native_style.css

$(ROOT)/contrib/IP-COUNTRY-REGION-CITY-ISP.BIN:
	mkdir -p $(ROOT)/contrib
	cd $(ROOT)/contrib && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY-ISP.BIN.gz
	cd $(ROOT)/contrib && gunzip IP-COUNTRY-REGION-CITY-ISP.BIN.gz
	cd $(ROOT)/contrib && rm -f IP-COUNTRY-REGION-CITY-ISP.BIN.md5 && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY-ISP.BIN.md5
	cd $(ROOT)/contrib && md5sum -c IP-COUNTRY-REGION-CITY-ISP.BIN.md5

ip2location: $(ROOT)/contrib/IP-COUNTRY-REGION-CITY-ISP.BIN
	cp $(ROOT)/contrib/IP-COUNTRY-REGION-CITY-ISP.BIN $(BIN)

.PHONY: ip2location