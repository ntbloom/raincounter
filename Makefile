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

TESTFLAGS   = -p 1
#TESTFLAGS  = -timeout 10s
#TESTFLAGS += -v

SQLFLAGS  = -h localhost
SQLFLAGS += -U postgres
SQLFLAGS += raincounter
DUMMY_DATA = $(HOMEDIR)pkg/test/dummy.sql

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

# docker control

docker-up:
	@$(COMPOSE) up $(COMPOSEFLAGS)

docker-down:
	@$(COMPOSE) down

docker-cycle: docker-down docker-up

docker-pglogs:
	@$(COMPOSE) logs -f postgresql

# postgresql control

psql:
	psql $(SQLFLAGS)

define enter_data
	@psql $(SQLFLAGS) -f $(DUMMY_DATA) > /dev/null
endef

# server
test-server: test-webdb test-receiver test-rest

test-server-race: clean-test test-common-race test-webdb-race test-receiver-race test-rest-race

test-webdb:
	@-go test $(TESTFLAGS) $(RAINCLOUD)/webdb/

test-webdb-race:
	@-go test -race $(TESTFLAGS) $(RAINCLOUD)/webdb/

test-receiver:
	$(call enter_data)
	@-go test $(TESTFLAGS) $(RAINCLOUD)/receiver/

test-receiver-race:
	$(call enter_data)
	@-go test -race $(TESTFLAGS) $(RAINCLOUD)/receiver/

test-rest:
	$(call enter_data)
	@-go test $(TESTFLAGS) $(RAINCLOUD)/api/

test-rest-race:
	$(call enter_data)
	@-go test -race $(TESTFLAGS) $(RAINCLOUD)/api/

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
	@go clean -testcache
	@- rm /tmp/rainbase.db /tmp/raincloud.db

clean-files:
	@- rm raincounter raincounter-race

clean-docker:
	- $(COMPOSE) down
	- docker volume prune -f | tail -1

clean: clean-test clean-files

