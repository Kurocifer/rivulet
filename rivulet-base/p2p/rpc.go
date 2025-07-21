package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// Message holds arbitrary data that is sent over each transport between two nodes.
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
