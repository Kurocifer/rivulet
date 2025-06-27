package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer represents the remote node over an established TCP connection
type TCPPeer struct {
	// The underlying connection of the peer
	net.Conn
	// If we dial a peer and retrieve a connectin => outbound == true
	// if we accept from a peer and retrieve a connection => outbound == false
	outbound bool

	Wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

// // Close implements the Peer interface.
// func (p *TCPPeer) Close() error {
// 	return p.conn.Close()
// }

// // RemoteAddr implements the peer interfeace, and will return the remote
// // address of it's underlying connection.
// func (p *TCPPeer) RemoteAddr() net.Addr {
// 	return p.conn.RemoteAddr()
// }

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
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

// Close implements the transport interface. It closes the connection set up for transport
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) ListeAddr() string {
	return t.ListenAddr
}

// Dial implements the transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("Dialed ", conn.RemoteAddr().String())

	go t.handleConnection(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Printf("TCP accept error: %v\n", err)
		}

		go t.handleConnection(conn, false)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)
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
		if err != nil {
			// if strings.Contains(err.Error(), "use of closed network connection") {
			// 	fmt.Println("Peer connection discontinued")
			// 	return
			// }

			if errors.Is(err, net.ErrClosed) {
				fmt.Println("Peer connection discontinued")
				return
			}
			fmt.Printf("TCP error %v\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			peer.Wg.Add(1)
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.Wg.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}
		// wait for other peers (go routines) to read from the connection
		// before proceeding to a next read cycle
		t.rpcch <- rpc
	}
}
