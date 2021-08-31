HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GW = $(HOMEDIR)pkg/gateway
SERVER = $(HOMEDIR)pkg/server
COMMON = $(HOMEDIR)pkg/common
EXE = ./raincounter

# docker-compose, for testing
COMPOSEFILE = $(HOMEDIR)pkg/test/docker-compose.yaml
COMPOSE = docker-compose -f $(COMPOSEFILE)
COMPOSEFLAGS  = --remove-orphans
COMPOSEFLAGS += -d

#TESTFLAGS  = -v

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

# server
test-server:
	@$(COMPOSE) up $(COMPOSEFLAGS)
	@sleep 2
	@- go test $(TESTFLAGS) $(SERVER)/...
	@$(COMPOSE) down


test-server-race: clean-test test-common-race
	@go clean -testcache
	@go test $(TESTFLAGS) -race $(SERVER)/...

### RUN ###

run-gateway: build
	@$(EXE) gateway

run-gateway-race: build-race
	@$(EXE)-race gateway


### CLEAN ###

clean-test:
	go clean -testcache
	- rm /tmp/rainbase.db /tmp/raincloud.db

clean-files:
	- rm raincounter raincounter-race

clean-docker:
	- $(COMPOSE) down
	- docker volume prune -f | tail -1

clean: clean-test clean-files clean-docker

