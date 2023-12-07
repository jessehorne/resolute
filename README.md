![Resolute Logo](./assets/logo-wide.png)

A simple chat service backend designed for ephemeral, anonymous conversations.

### Project status: NOT READY FOR PRODUCTION

# Overview

Rësolute started as a learning project to further my understanding of secure and scalable communications systems.

# Features

* E2EE for messages
* Enforce WSS
* Messages are not stored on disk
* No data is collected

more coming soon...

# Usage

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
	if err := s.Listen("./cert.pem", "./key.pem"); err != nil {
		fmt.Println(err)
	}
}
```

Also check out `examples/`.

# Client API Documentation

Please see [Client-API.md](./Client-API.md).

# License

See `./LICENSE`

more coming soon...
