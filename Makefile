HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SBC = $(HOMEDIR)sbc
COMMON = $(HOMEDIR)common
IGNORE = grep -v "no test files"

build-rainbase:
	@go build -v $(SBC)/rainbase.go

run-rainbase: build-rainbase
	./rainbase

test-common:
	@go test  $(COMMON)/...

test-rainbase: test-common
	@go test $(SBC)/...

clean:
	go clean -testcache

