package serverUtils

import (
	crypt "github.com/kurocifer/rivulet/rivulet-base/crypto"
	"github.com/kurocifer/rivulet/rivulet-base/p2p"
	"github.com/kurocifer/rivulet/rivulet-base/server"
	"github.com/kurocifer/rivulet/rivulet-base/store"
)

func MakeServer(listenAddr string, nodes ...string) *server.FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.DefaultHandSake,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := server.FileServerOpts{
		EncKey:            crypt.NewEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: store.CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := server.NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}
