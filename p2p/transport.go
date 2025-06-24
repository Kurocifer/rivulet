package p2p

import "net"

// Peer is an interface representing a remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything that handles communication
// between two nodes in the netowrk. It can be of the for
// (UDP, websockets, TCP, ...)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	ListeAddr() string
}
