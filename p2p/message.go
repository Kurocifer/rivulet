package p2p

import "net"

// Message holds arbitrary data that is sent over each transport between two nodes.
type Message struct {
	From    net.Addr
	Payload []byte
}
