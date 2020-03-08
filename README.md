[![CircleCI](https://circleci.com/gh/spatialcurrent/go-walker/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-walker/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-walker)](https://goreportcard.com/report/spatialcurrent/go-walker)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-walker?status.svg)](https://godoc.org/github.com/spatialcurrent/go-walker) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-walker/blob/master/LICENSE)

# go-walker

## Description

**go-walker** is a library that supports walking a directory structure, including local and remote file system.  This library includes the `iterator`, `modeutil`, `oserror`, `pathutil`, and `walker` packages.  The `walker` package is the main package and the rest are supporting packages.

# Usage

**Go**

You can import **go-walker** as a library with:

```go
import (
  "github.com/spatialcurrent/go-walker/pkg/walker"
)
```

The easiest pattern is to create a walker that can be reused multiple times.

```go
w, err := walker.NewWalker(&walker.NewWalkerInput{
  SkipPath:      nil,
  SkipFile:      nil,
  SkipLink:      nil,
  ErrorLinkStat: nil,
  ErrorStat:     nil,
  ErrorWalk:     nil,
  Limit:         walker.NoLimit,
})
```

You can call the `Walk` method of the walker multiple times, and in concurrent goroutines.  The `Walk` method requires a root context, slice of root directories, and a walk function.  The context enables easier use of the walker across multiple API boundaries.  The `context` is not used by the walker itself, but is simply passed through to the walk function.

```go
n, err := w.Walk(context.Background(), []string{"testdata"}, func(ctx context.Context, p string, f File) error {
  // retrieve value from context if you like
  // do something with each file
  return nil
})
```

See [walker](https://godoc.org/github.com/spatialcurrent/go-walker/pkg/walker) in GoDoc for API documentation.

# Testing

To run Go tests use `make test_go` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-walker/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
