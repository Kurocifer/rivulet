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

type TCPTransport struct {
	listenerAddress string
	listener        net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenerAddr string) *TCPTransport {
	return &TCPTransport{
		listenerAddress: listenerAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenerAddress)
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
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Println(conn.RemoteAddr().String())
	fmt.Printf("New incoming connection %+v\n", peer)
}
