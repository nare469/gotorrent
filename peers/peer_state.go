package peers

import (
	"errors"
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
	requestChan        chan uint32
	state              *PeerState
	pieceInfo          *PieceInfo
}

type PeerState struct {
	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool
}

type PieceInfo struct {
	data    [][]byte
	counter uint32
	index   uint32
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
		requestChan:        make(chan uint32),
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

func (p *PeerConnection) receiveBlock(block []byte) {
	p.pieceInfo.data[p.pieceInfo.counter] = block
	p.pieceInfo.counter += 1
	if p.pieceInfo.counter == uint32(len(p.pieceInfo.data)) {
		return
	}
	p.requestChan <- p.pieceInfo.counter
}

func (p *PeerConnection) choosePieceToRequest() error {
	if p.pieceInfo != nil {
		return errors.New("Peer already requesting")
	}
	for i, val := range p.bitfield {
		if val {
			length, _ := p.attrs.PieceLength()
			length /= BLOCK_SIZE
			p.pieceInfo = &PieceInfo{
				data:    make([][]byte, length),
				counter: 0,
				index:   uint32(i),
			}

			p.requestChan <- 0
			return nil
		}
	}
	return errors.New("Peer has no pieces")
}
