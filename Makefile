HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
RAINBASE = $(HOMEDIR)pkg/rainbase
RAINCLOUD = $(HOMEDIR)pkg/raincloud
COMMON = $(HOMEDIR)pkg/common
EXE = ./raincounter
EXERACE = $(EXE)-race

# docker-compose, for testing
COMPOSEFILE = $(HOMEDIR)pkg/test/docker-compose.yaml
COMPOSE = docker-compose -f $(COMPOSEFILE)
COMPOSEWAIT = 10
COMPOSEFLAGS  = --build --remove-orphans
COMPOSEFLAGS += -d

FRONTEND_COMPOSEFILE = $(HOMEDIR)docker/docker-compose.yml
FRONTEND_COMPOSE = docker-compose -f $(FRONTEND_COMPOSEFILE)

TESTFLAGS   = -p 1
#TESTFLAGS  = -timeout 10s
#TESTFLAGS += -v

SQLFLAGS  = -h localhost
SQLFLAGS += -U postgres
SQLFLAGS += raincounter
DUMMY_DATA = $(HOMEDIR)pkg/test/dummy.sql
CLEAR_SQL = $(HOMEDIR)pkg/test/clear.sql

# for the front end
DOCKERDIR = $(HOMEDIR)docker

### DEPLOY ###

DEVCFG = $(HOMEDIR)config/insecure.yml
DEVFLAGS = --config $(DEVCFG)

DEVRUN = $(EXE) $(DEVFLAGS)

dev-server:
	@$(DEVRUN) server

dev-receiver:
	@$(DEVRUN) receiver

dev-rainbase:
	@$(DEVRUN) rainbase

### BUILD ###

build:
	@go build -v
	@# add the build dependencies to the front-end docker toolchain
	@cp $(EXE) $(DOCKERDIR)
	@cp $(HOMEDIR)pkg/test/pgschema/schema.sql $(DOCKERDIR)/pgschema/00-schema.sql
	@cp $(HOMEDIR)pkg/test/dummy.sql $(DOCKERDIR)/pgschema/99-dummy.sql

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

frontend-up:
	@$(FRONTEND_COMPOSE) up $(COMPOSEFLAGS)

frontend-down:
	@$(FRONTEND_COMPOSE) down

# postgresql control

psql:
	psql $(SQLFLAGS)

define enter_data
	@psql $(SQLFLAGS) -f $(DUMMY_DATA) > /dev/null
endef

remove-data:
	psql $(SQLFLAGS) -f $(CLEAR_SQL)

# server
enter-data:
	@$(call enter_data)

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
	@$(EXERACE) gateway

run-receiver: build
	@$(EXE) receiver

run-receiver-race: build-race
	@$(EXERACE) race receiver

run-server: build
	@$(EXE) server

run-server-race: build-race
	@$(EXERACE) server

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

