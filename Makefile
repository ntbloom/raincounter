HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GW = $(HOMEDIR)pkg/gateway
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


### RUN ###

run-gateway: build
	@$(EXE) gateway

run-gateway-race: build-race
	@$(EXE)-race gateway


### CLEAN ###

clean-test:
	go clean -testcache

clean-files:
	- rm raincounter

clean: clean-test clean-files

