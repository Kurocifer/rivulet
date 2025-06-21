package p2p

import (
	"fmt"
	"net"
	"strings"
)

// TCPPeer represents the remote node over an established TCP connection
type TCPPeer struct {
	// comm is the underlying connection of the peer
	conn net.Conn
	// If we dial a peer and retrieve a connectin => outbound == true
	// if we accept from a peer and retrieve a connection => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements the Peer interface.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(tcptransferOpts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: tcptransferOpts,
		rpcch:            make(chan RPC),
	}
}

// Consume implements the transport interface, which will return a read-only channel
// for reading incoming messages recieved from another peer in the netwrok
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
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

		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err = t.HandShakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)

		if strings.Contains(err.Error(), "use of closed network connection") {
			fmt.Println("Peer connection discontinued")
			return
		}

		if err != nil {
			fmt.Printf("TCP error %v\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()

		t.rpcch <- rpc
	}
}
