TEST?=./...
ifeq ("$(shell uname)","Darwin")
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
else
NCPU ?= $(shell cat /proc/cpuinfo | grep processor | wc -l)
endif
TEST_ARGS=-v
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)

ONIG_VERSION?=6.2.0

default: build

deps:
	go get -u golang.org/x/lint/golint

test:
	go test $(TEST) $(TEST_ARGS) $(TEST_OPTIONS)
	go test -race $(TEST) -coverprofile=coverage.txt -covermode=atomic

lint:
	golint -set_exit_status $(TEST)

order: deps lint
	go mod tidy
	git diff go.mod

install:
	go install .

onigmo:
	test -d tmp || mkdir tmp
	cd ./tmp && \
		curl -sLO https://github.com/k-takata/Onigmo/releases/download/Onigmo-${ONIG_VERSION}/onigmo-${ONIG_VERSION}.tar.gz && \
		tar xfz onigmo-${ONIG_VERSION}.tar.gz && \
		cd onigmo-${ONIG_VERSION} && ./configure && make && sudo make install

.PHONY: default lint test deps
