HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
RAINBASE = $(HOMEDIR)pkg/rainbase
RAINCLOUD = $(HOMEDIR)pkg/raincloud
COMMON = $(HOMEDIR)pkg/common
EXE = ./raincounter

# docker-compose, for testing
COMPOSEFILE = $(HOMEDIR)pkg/test/docker-compose.yaml
COMPOSE = docker-compose -f $(COMPOSEFILE)
COMPOSEWAIT = 10
COMPOSEFLAGS  = --remove-orphans
COMPOSEFLAGS += -d

TESTFLAGS  = -timeout 10s
TESTFLAGS += -p 1
#TESTFLAGS += -v

### BUILD ###

build:
	@go build -v

build-race: clean
	@go build -race -o $(EXE)-race

### TEST ###

# common
test-common:
	@go test $(TESTFLAGS) $(COMMON)/...

test-common-race:
	@go test $(TESTFLAGS) -race $(COMMON)/...

test-all: clean test-common test-gateway test-server

test-all-race: clean test-gateway-race test-server-race

# rainbase

test-gateway:
	@go test $(TESTFLAGS) $(RAINBASE)/...

test-gateway-race: clean-test test-common-race
	@go clean -testcache
	@go test $(TESTFLAGS) -race $(RAINBASE)/...

docker-up:
	@$(COMPOSE) up $(COMPOSEFLAGS)

docker-down:
	@$(COMPOSE) down

docker-cycle: docker-down docker-up

docker-pglogs:
	@$(COMPOSE) logs -f postgresql

psql:
	psql -U postgres -h localhost raincounter

# server
test-server: test-webdb test-receiver test-rest

test-server-race: clean-test test-common-race test-webdb-race test-receiver-race test-rest-race

test-webdb:
	@-go test $(TESTFLAGS) $(RAINCLOUD)/webdb/

test-webdb-race:
	@-go test -race $(TESTFLAGS) $(RAINCLOUD)/webdb/

test-receiver:
	@-go test $(TESTFLAGS) $(RAINCLOUD)/receiver/

test-receiver-race:
	@-go test -race $(TESTFLAGS) $(RAINCLOUD)/receiver/

test-rest:
	@-go test $(TESTFLAGS) $(RAINCLOUD)/rest/

test-rest-race:
	@-go test -race $(TESTFLAGS) $(RAINCLOUD)/rest/

### RUN ###

run-gateway: build
	@$(EXE) gateway

run-gateway-race: build-race
	@$(EXE)-race gateway

run-receiver: build
	@$(EXE) receiver

run-receiver-race: build-race
	@$(EXE) race receiver

run-server: build
	@$(EXE) server

run-server-race: build-race
	@$(EXE) race server

### LINT ###

lint:
	@golangci-lint run


### CLEAN ###

clean-test:
	go clean -testcache
	- rm /tmp/rainbase.db /tmp/raincloud.db

clean-files:
	- rm raincounter raincounter-race

clean-docker:
	- $(COMPOSE) down
	- docker volume prune -f | tail -1

clean: clean-test clean-files

