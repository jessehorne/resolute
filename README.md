![Resolute Logo](./assets/logo-wide.png)

A simple chat service backend designed for ephemeral, anonymous conversations.

### Project status: NOT READY FOR PRODUCTION

# Overview

Rësolute started as a learning project to further my understanding of secure and scalable communications systems. It aims to  be used for short term communications where message history isn't important. It also purposely excludes any form of long-term on-disk user identification. All ID's follow a standard format but are cryptographically secure and stored only in memory. This memory **is not yet properly protected using something like memguard**, but that functionality is planned.

All message content is E2E encrypted from client to client, along with enforced TLS for connections to the server. Currently, 2048-bit RSA key pairs are generated for every connection made (see `pkg/v1/client/client.go`) in the client library. These are stored in memory but **not currently properly protected**. This E2EE is up to clients discretion, however. I just provide the ability to share public keys in the protocol. In the future, I'd like to remove as much data as possible that is unencrypted. Please see the client API documentation for more details on what data the server itself is aware of.

**This software does not yet replace your INSERT SECURE COMMS APP HERE. Hopefully one day it will, but not likely soon.

## Rooms

On a server instance, any user, without authentication nor authorization, can create a Room. A room can be a conversation between 1+ members. Creating a room automatically subscribes your client to any messages sent by other users in that room. Public keys, usernames and so on are stored in a Rooms data structure for the lifetime of that room.

To join a room, the creator of the room must first generate one-time or "forever" keys. These keys work to authenticate someone. They are 32 byte random base64 encoded strings, as are all IDs. The library I am using for this is called uniuri which uses crypto/rand's rand.Read. One-time keys can be used once and then are erased. Forever keys can be used by anyone any number of times. The creator of a room can send a packet to erase all one-time and forever keys in one go.

## Messages

Messages are encrypted on the client using each destination-users public key in a Room and are sent to the server from the client where the server then relays each message to the appropriate users. Please see the client API documentation for more information.

**more coming soon...**

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
