package main

import (
	"fmt"
	"log"

	"github.com/kurocifer/rivulet/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":8080",
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.DefaultHandSake,
		OnPeer:        onPeer,
	}
	tr := p2p.NewTCPTransport(tcpTransportOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}

func onPeer(peer p2p.Peer) error {
	fmt.Println("What am I even doing ??????")
	peer.Close()
	return nil
}
