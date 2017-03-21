package peers

import (
	"github.com/nare469/gotorrent/parser"
	"net"
)

const (
	CHOKE byte = iota
	UNCHOKE
	INTERESTED
	UNINTERESTED
	HAVE
	BITFIELD
	REQUEST
	PIECE
	CANCEL
)

type PeerConnection struct {
	peer     parser.Peer
	conn     *net.TCPConn
	bitfield []byte
	quitChan chan bool
	haveChan chan uint32
	state    *PeerState
}

type PeerState struct {
	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool
}

func NewPeerConnection(p parser.Peer, conn *net.TCPConn) *PeerConnection {
	return &PeerConnection{
		peer: p,
		conn: conn,
		state: &PeerState{
			amChoking:      true,
			amInterested:   false,
			peerChoking:    true,
			peerInterested: false,
		},
		quitChan: make(chan bool),
		haveChan: make(chan uint32),
	}
}
