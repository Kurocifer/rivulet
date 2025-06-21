package main

import (
	"log"

	"github.com/kurocifer/rivulet/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":8080",
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.DefaultHandSake,
	}
	tr := p2p.NewTCPTransport(tcpTransportOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
