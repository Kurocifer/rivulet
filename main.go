package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/kurocifer/rivulet/p2p"
)

func main() {

	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 2)

	go func() {
		log.Fatal(s2.Start())
	}()
	time.Sleep(time.Second * 2)

	// data := bytes.NewReader([]byte("Yeah we know Ulquiorra is him!"))
	// s2.Store("Espada Facts", data)
	// time.Sleep(time.Millisecond * 5)

	r, err := s2.Get("Espada Facts")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

// func onPeer(peer p2p.Peer) error {
// 	fmt.Println("What am I even doing ??????")
// 	peer.Close()
// 	return nil
// }

func makeServer(listenerAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenerAddr,
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.DefaultHandSake,
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOPts{
		StoreageRoot:      listenerAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	server := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = server.onPeer

	return server
}
