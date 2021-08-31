export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ATOMIX_BENCHMARKS_VERSION := latest

all: build

build: # @HELP build the source code
build:
	GOOS=linux GOARCH=amd64 go build -o build/_output/atomix-benchmarks ./cmd/atomix-benchmarks

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/atomix-benchmarks/...

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}
