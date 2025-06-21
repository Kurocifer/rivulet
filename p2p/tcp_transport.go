package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over an established TCP connection
type TCPPeer struct {
	// comm is the underlying connection of the peer
	conn net.Conn
	// If we dial a peer and retrieve a connectin => outbound == true
	// if we accept from a peer and retrieve a connection => outbound == false
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(tcptransferOpts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: tcptransferOpts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}

		fmt.Printf("new incoming connection %+v", conn)
		go t.handleConnection(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	if err := t.HandShakeFunc(peer); err != nil {
		conn.Close()
		return
	}

	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", msg)
	}
}
