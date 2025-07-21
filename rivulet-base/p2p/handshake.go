package p2p

type HandshakeFunc func(Peer) error

func DefaultHandSake(Peer) error { return nil }
