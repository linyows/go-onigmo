Onigmo bindings for Go
======================

It binds the regular expression library Onigmo to Go.

[![Travis](https://img.shields.io/travis/linyows/go-onigmo.svg?style=flat-square)][travis]
[![GitHub release](http://img.shields.io/github/release/linyows/go-onigmo.svg?style=flat-square)][release]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[travis]: https://travis-ci.org/linyows/go-onigmo
[release]: https://github.com/linyows/go-onigmo/releases
[license]: https://github.com/linyows/go-onigmo/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/linyows/go-onigmo

Benchmarks
----------

These are the benchmarks as they are defined in Go's regexp package.

```sh
$ go test -bench RE2 | sed 's/RE2/Regexp/' > before
$ go test -bench Onigmo | sed 's/Onigmo/Regexp/' > after
$ benchcmp before after
benchmark             old ns/op     new ns/op     delta
BenchmarkRegexp-4     25775         31043         +20.44%
```

Installation
------------

```sh
$ git clone git@github.com:linyows/go-onigmo.git && cd go-onigmo
$ make onigmo
```

To install, use `go get`:

```sh
$ go get -d github.com/linyows/go-onigmo
```

Contribution
------------

1. Fork ([https://github.com/linyows/go-onigmo/fork](https://github.com/linyows/go-onigmo/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

Author
------

[linyows](https://github.com/linyows)
