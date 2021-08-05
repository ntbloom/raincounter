HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SBC = $(HOMEDIR)sbc
COMMON = $(HOMEDIR)common

# common
test-common:
	@go test  $(COMMON)/...

test-common-race:
	@go test -race $(COMMON)/...


# rainbase
build-rainbase:
	@go build -v $(SBC)/rainbase.go

build-rainbase-race: clean
	@go build -race -o rainbase-race -v $(SBC)/rainbase.go

test-rainbase: test-common
	@go test $(SBC)/...

test-rainbase-race: clean-test test-common-race
	@go clean -testcache
	@go test -race $(SBC)/...

run-rainbase: build-rainbase
	./rainbase

run-rainbase-race: build-rainbase-race
	./rainbase-race


clean-test:
	go clean -testcache

clean-files:
	- rm rainbase rainbase-race

clean: clean-test clean-files

