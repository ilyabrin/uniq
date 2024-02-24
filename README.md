# UNIQ
[![Go Tests](https://github.com/ilyabrin/uniq/actions/workflows/test.yml/badge.svg)](https://github.com/ilyabrin/uniq/actions/workflows/test.yml)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/ilyabrin/uniq.svg?style=flat)]()
[![PR's Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](http://makeapullrequest.com)
[![GitHub Release](https://img.shields.io/github/release/ilyabrin/uniq.svg?style=flat)]()

Generate unique ID with given length

## Installation and usage

### Installation

`go get -u github.com/ilyabrin/uniq`

### Usage

```go
package main

import "github.com/ilyabrin/uniq"

const length = 10
const maxUpCase = 6

func main() {
	println("Generated ID:", uniq.New(length, maxUpCase))
}
```
Example: https://go.dev/play/p/i6vKfk-evij

## License

MIT
