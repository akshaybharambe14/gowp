# GOWP

[![PkgGoDev](https://pkg.go.dev/badge/github.com/akshaybharambe14/gowp)](https://pkg.go.dev/github.com/akshaybharambe14/gowp)
[![Build and Test Status](https://github.com/akshaybharambe14/gowp/workflows/Build%20and%20test/badge.svg)](https://github.com/akshaybharambe14/gowp/actions?query=workflow%3A%22Build+and+test%22)
[![PkgGoDev](https://goreportcard.com/badge/github.com/akshaybharambe14/gowp)](https://goreportcard.com/report/github.com/akshaybharambe14/gowp)

Package gowp provides synchronization, concurrency limiting, error propagation, and Context cancellation for a group of workers/goroutines. Designed with performance in mind.

## Features

- Context cancellation, won't process new jobs if parent context gets cancelled
- Error propagation, will exit and won't accept any more jobs if any error occurs, returns first error
- Concurrency limiting

## Installation

```bash
go get -u github.com/akshaybharambe14/gowp
```
