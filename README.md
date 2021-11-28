# GOWP

[![PkgGoDev](https://pkg.go.dev/badge/github.com/akshaybharambe14/gowp)](https://pkg.go.dev/github.com/akshaybharambe14/gowp)
[![Build and Test Status](https://github.com/akshaybharambe14/gowp/workflows/Build%20and%20test/badge.svg)](https://github.com/akshaybharambe14/gowp/actions?query=workflow%3A%22Build+and+test%22)
[![PkgGoDev](https://goreportcard.com/badge/github.com/akshaybharambe14/gowp)](https://goreportcard.com/report/github.com/akshaybharambe14/gowp)

Package gowp (Go Worker-Pool) provides concurrency limiting, error propagation, and Context cancellation for a group of workers/goroutines.

## Features

- Context cancellation, won't process new tasks if parent context gets cancelled.
- Error propagation, will return the first error encountered.
- Exit on error, if specified, it won't process further tasks if an error is encountered.
- Concurrency limiting

## Installation

```bash
go get -u github.com/akshaybharambe14/gowp
```

## Why?

Goroutines are cheap in terms of memory, but not free. If you want to achieve extreme performance, you need to limit the number of goroutines you use. Also, in real world applications, you need to take care of the failures. This package does that for you.

This package works best with short bursts of tasks. We have other packages that are meant for running huge number of tasks.

### Why yet another worker pool implementation

I wanted to build a perfect worker pool implementation with above specified features. We have other implementations, but I think they do a lot of work in background if used to execute fixed set of tasks. Gowp outperforms some of them ([see benchmarks](#Benchmarks)).

## Benchmarks

For following benchmarks, a single operation means running 10 tasks over 4 workers.
Packages compared (results are in same sequence as packages are listed):

- github.com/akshaybharambe14/gowp
- github.com/gammazero/workerpool
- github.com/alitto/pond

```bash
$ go test -bench=. -benchmem github.com/akshaybharambe14/gowp/benchmarks
goos: windows
goarch: amd64
pkg: github.com/akshaybharambe14/gowp/benchmarks
cpu: Intel(R) Core(TM) i5-4200M CPU @ 2.50GHz
Benchmark_simple_gowp-4                    79420             13663 ns/op             704 B/op         13 allocs/op
Benchmark_simple_workerpool-4              40402             29716 ns/op             952 B/op         17 allocs/op
Benchmark_simple_pond-4                    57692             20689 ns/op             968 B/op         20 allocs/op
PASS
ok      github.com/akshaybharambe14/gowp/benchmarks     4.323s
```

## Examples

see [package examples](https://pkg.go.dev/github.com/akshaybharambe14/gowp#pkg-examples)

## Contact

[Akshay Bharambe](https://twitter.com/akshaybharambe1)

---

If this is not something you are looking for, you can check other similar packages on [go.libhunt.com](https://go.libhunt.com/categories/493-goroutines).

Do let me know if you have any feedback. Leave a ⭐ if you like this work.
