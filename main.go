package main

import (
	"log"

	"github.com/kurocifer/rivulet/p2p"
)

func main() {

	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	if err := s2.Start(); err != nil {
		log.Fatal(err)
	}
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
		StoreageRoot:      "rivulet",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	server := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = server.onPeer

	return server
}
