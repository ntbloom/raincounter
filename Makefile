HOMEDIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SBC = $(HOMEDIR)sbc

rainbase:
	@go build -v $(SBC)/rainbase.go

test-rainbase:
	@go test -v $(SBC)/... | grep -v "no test files"