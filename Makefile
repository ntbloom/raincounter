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

build-rainbase-race:
	@go build -race -o rainbase-race -v $(SBC)/rainbase.go

test-rainbase: test-common
	@go test $(SBC)/...

test-rainbase-race: test-common-race
	@go test -race $(SBC)/...

run-rainbase: build-rainbase
	./rainbase

race-rainbase-run: build-rainbase-race
	./rainbase-race




clean:
	@- go clean -testcache
	@- rm rainbase rainbase-race

