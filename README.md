Resolute
===

A simple chat service backend designed for anonymous, short-term conversations.

*Please note: This is a learning project. Use at your own risk. While I'm working to build in E2EE and other security/privacy measures over time, this project is unfinished and untested.*

Currently, resolute is best suited for short term conversations where chat history and user identification isn't needed.

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

# Client API Documentation

Please see [Client-API.md](./Client-API.md).

# License

See `./LICENSE`

more coming soon...
