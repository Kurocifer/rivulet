package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/kurocifer/rivulet/p2p"
)

type FileServerOPts struct {
	StoreageRoot      string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOPts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store *Store
	quit  chan struct{}
}

func NewFileServer(opts FileServerOPts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StoreageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOPts: opts,
		store:          NewStore(storeOpts),
		quit:           make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Stop() {
	close(s.quit)
}

func (s *FileServer) onPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[peer.RemoteAddr().String()] = peer

	log.Println("connected with remote ", peer.RemoteAddr())

	// Yeah true doing this just to satisfy definition of the required function is not normal
	return nil
}
func (s *FileServer) loop() {
	defer func() {
		log.Println("File server stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Printf("%s sent %s", msg.From, string(msg.Payload))

		case <-s.quit:
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() {
	for _, addr := range s.BootstrapNodes {
		go func(addr string) {
			fmt.Println("attempting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootstrapNodes) > 0 {
		s.bootstrapNetwork()
	}

	s.loop()

	return nil
}
