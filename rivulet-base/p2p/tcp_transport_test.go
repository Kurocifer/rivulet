package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpTransportOpts := &TCPTransportOpts{
		ListenAddr: ":4000",
	}
	listnerAddr := ":4000"
	tr := NewTCPTransport(*tcpTransportOpts)

	assert.Equal(t, listnerAddr, tr.ListenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
