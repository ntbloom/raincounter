HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GW = $(HOMEDIR)pkg/gateway
SERVER = $(HOMEDIR)pkg/server
COMMON = $(HOMEDIR)pkg/common
EXE = ./raincounter

### BUILD ###

build:
	@go build -v

build-race: clean
	@go build -race -o $(EXE)-race

### TEST ###

# common
test-common:
	@go test  $(COMMON)/...

test-common-race:
	@go test -race $(COMMON)/...

# gateway

test-gateway: test-common
	@go test $(GW)/...

test-gateway-race: clean-test test-common-race
	@go clean -testcache
	@go test -race $(GW)/...

# server
test-server: test-common
	@go test $(SERVER)/...

test-server0race: clean-test test-common-race
	@go clean -testcache
	@go test -race $(SERVER)/...

### RUN ###

run-gateway: build
	@$(EXE) gateway

run-gateway-race: build-race
	@$(EXE)-race gateway


### CLEAN ###

clean-test:
	go clean -testcache

clean-files:
	- rm raincounter raincounter-race

clean: clean-test clean-files

