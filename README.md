# gordon

Go released binaries manager

[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gordon)](https://goreportcard.com/report/github.com/kyoh86/gordon)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/gordon.svg)](https://codecov.io/gh/kyoh86/gordon)

## What's this?

`homebrew` is awesome.

If you created tools by `go` and release it on GitHub Releases, you can use `gordon` instead of `homebrew`.

```console
$ gordon install kyoh86/richgo
$ gordon update
$ gordon uninstall kyoh86/richgo
$ gordon dump GordonDumpFile
$ gordon restore GordonDumpFile
```

`gordon install` will download asset for local OS/architecture from GitHub Releases,
and link to the executables in the asset from `$HOME/.local/bin`.

CAUTION: now its version is 0.x (not stable).

## Install

```console
$ go get github.com/kyoh86/gordon@latest
```

If you see this error:

```
go: cannot use path@version syntax in GOPATH mode
then run
```

```console
$ GO111MODULE=on go get github.com/kyoh86/gordon@latest
```

## Usage

```console
$ gordon --help
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
