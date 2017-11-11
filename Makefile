TEST?=./...
NAME = "$(shell awk -F\" '/^const Name/ { print $$2; exit }' version.go)"
VERSION = "$(shell awk -F\" '/^const Version/ { print $$2; exit }' version.go)"

ONIG_VERSION?=6.1.3

default: test

onigmo:
	curl -sLO https://github.com/k-takata/Onigmo/releases/download/Onigmo-${ONIG_VERSION}/onigmo-${ONIG_VERSION}.tar.gz
	tar xfz onigmo-${ONIG_VERSION}.tar.gz
	pushd onigmo-${ONIG_VERSION} && ./configure && make && sudo make install && popd

deps:
	go get -d -t ./...

depsdev:
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr

test:
	go test -v $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	go test -race $(TEST) $(TESTARGS)

ci:
	$(MAKE) onigmo
	$(MAKE) test
	go install .

dist:
	ghr v$(VERSION) pkg

.PHONY: default dist test test deps
