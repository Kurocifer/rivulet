package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	serverUtils "github.com/kurocifer/rivulet/rivulet-base/serverUtils"
)

var espadas = []string{
	"Stark",
	"Baragan",
	"Haribel",
	"Ulquirra",
	"Nnoitra",
	"Grimmjow",
	"Zommari",
	"Szayel",
	"Aaroniero",
	"Yammy",
}

// func makeServer(listenAddr string, nodes ...string) *FileServer {
// 	tcptransportOpts := p2p.TCPTransportOpts{
// 		ListenAddr:    listenAddr,
// 		HandShakeFunc: p2p.DefaultHandSake,
// 		Decoder:       p2p.DefaultDecoder{},
// 	}
// 	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

// 	fileServerOpts := FileServerOpts{
// 		EncKey:            newEncryptionKey(),
// 		StorageRoot:       listenAddr + "_network",
// 		PathTransformFunc: CASPathTransformFunc,
// 		Transport:         tcpTransport,
// 		BootstrapNodes:    nodes,
// 	}

// 	s := NewFileServer(fileServerOpts)

// 	tcpTransport.OnPeer = s.OnPeer

// 	return s
// }

func main() {
	// s1 := serverUtils.MakeServer(":3000", "")
	s2 := serverUtils.MakeServer(":7000", "")
	s3 := serverUtils.MakeServer(":5000", ":3000", ":7000")

	// go func() { log.Fatal(s1.Start()) }()
	// time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	for i, name := range espadas {
		key := fmt.Sprintf("%s_Espada_%d.bleach", name, i)
		data := bytes.NewReader([]byte("Yare Yare go is really awesome"))
		s3.Store(key, data)

		if err := s3.Sstore.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}

// func onPeer(peer p2p.Peer) error {
// 	fmt.Println("What am I even doing ??????")
// 	peer.Close()
// 	return nil
// }
