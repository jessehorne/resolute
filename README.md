Resolute
===

A simple chat service backend designed for anonymity.

*Please note: This is a learning project. Use at your own risk.*

# Usage

## Create Server

```go
package main

import (
	"fmt"

	"github.com/jessehorne/resolute/pkg/v1/resolute"
)

func main() {
	host := "127.0.0.1:5656"
	
	fmt.Println("listening on:", host)
	
	s := resolute.NewServer("/v1", host)
	if err := s.Listen(); err != nil {
		fmt.Println(err)
	}
}

```

# License

See `./LICENSE`

more coming soon...
