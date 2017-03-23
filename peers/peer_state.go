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

const BLOCK_SIZE = 16384

type PeerConnection struct {
	peer               parser.Peer
	conn               *net.TCPConn
	attrs              *parser.TorrentAttrs
	bitfield           []bool
	canReceiveBitfield bool
	quitChan           chan bool
	state              *PeerState
}

type PeerState struct {
	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool
}

func NewPeerConnection(p parser.Peer, conn *net.TCPConn, attrs *parser.TorrentAttrs) *PeerConnection {
	pieces, _ := attrs.NumPieces()
	return &PeerConnection{
		peer:  p,
		conn:  conn,
		attrs: attrs,
		state: &PeerState{
			amChoking:      true,
			amInterested:   false,
			peerChoking:    true,
			peerInterested: false,
		},
		quitChan:           make(chan bool),
		canReceiveBitfield: true,
		bitfield:           make([]bool, pieces),
	}
}

func (p *PeerConnection) setBitfield(bitfield []byte) {
	i := 0
	for _, b := range bitfield {
		j := 7
		for b != 0 {
			if i+j < len(p.bitfield) {
				if b%2 == 1 {
					p.bitfield[i+j] = true
				} else {
					p.bitfield[i+j] = false
				}
			}
			b = b >> 1
			j--
		}
		i += 8
	}
}

func (p *PeerConnection) setHasPiece(piece uint32) {
	p.bitfield[piece] = true
}

func (p *PeerConnection) choosePieceToRequest() {

}
