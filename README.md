# UNIQ

Generate unique ID with given length

## Installation and usage

### Installation

`go get -u github.com/ilyabrin/uniq`

### Usage

```go
package main

import (
    "fmt"
    "github.com/ilyabrin/uniq"
)

func main() {
    generatedID := uniq.New(10,4)
    fmt.Println("Generated ID:", generatedID)
}
```

## License

MIT
