HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GW = $(HOMEDIR)pkg/gateway
SERVER = $(HOMEDIR)pkg/server
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

# gateway

test-gateway:
	@go test $(TESTFLAGS) $(GW)/...

test-gateway-race: clean-test test-common-race
	@go clean -testcache
	@go test $(TESTFLAGS) -race $(GW)/...

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
	@-go test $(TESTFLAGS) $(SERVER)/webdb/

test-webdb-race:
	@-go test -race $(TESTFLAGS) $(SERVER)/webdb/

test-receiver:
	@-go test $(TESTFLAGS) $(SERVER)/receiver/

test-receiver-race:
	@-go test -race $(TESTFLAGS) $(SERVER)/receiver/

test-rest:
	@-go test $(TESTFLAGS) $(SERVER)/rest/

test-rest-race:
	@-go test -race $(TESTFLAGS) $(SERVER)/rest/

### RUN ###

run-gateway: build
	@$(EXE) gateway

run-gateway-race: build-race
	@$(EXE)-race gateway

run-receiver: build
	@$(EXE) receiver

run-receiver-race: build-race
	@$(EXE) race receiver

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

