package p2p

// Peer is an interface representing a remote node
type Peer interface {
	Close() error
}

// Transport is anything that handles communication
// between two nodes in the netowrk. It can be of the for
// (UDP, websockets, TCP, ...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
