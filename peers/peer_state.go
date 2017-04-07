package peers

import (
	"errors"
	"fmt"
	"github.com/nare469/gotorrent/download_state"
	"github.com/nare469/gotorrent/parser"
	"net"
	"sync"
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
	mu      sync.Mutex
}

func NewPeerConnection(p parser.Peer, conn *net.TCPConn, attrs *parser.TorrentAttrs) *PeerConnection {
	pieces, _ := attrs.NumPieces()

	length, _ := attrs.PieceLength()
	length /= BLOCK_SIZE

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
		pieceInfo: &PieceInfo{
			data:    make([][]byte, length),
			counter: 0,
			index:   0,
		},
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
	p.pieceInfo.mu.Lock()
	fmt.Println(p.pieceInfo.counter)
	p.pieceInfo.data[p.pieceInfo.counter] = block
	p.pieceInfo.counter += 1

	if p.pieceInfo.counter == uint32(len(p.pieceInfo.data)) {
		p.pieceInfo.mu.Unlock()
		go download_state.WritePiece(p.pieceInfo.data, p.pieceInfo.index)
		p.choosePieceToRequest()
		return
	} else {
		p.pieceInfo.mu.Unlock()
	}

	p.requestChan <- p.pieceInfo.counter * uint32(len(block))
}

func (p *PeerConnection) choosePieceToRequest() error {
	for i, val := range p.bitfield {
		state := download_state.GetPieceState(uint32(i))
		if val && state == download_state.MISSING {
			p.pieceInfo.mu.Lock()
			p.pieceInfo.index = uint32(i)
			p.pieceInfo.counter = 0
			p.pieceInfo.mu.Unlock()

			download_state.SetPieceState(uint32(i), download_state.IN_PROGRESS)
			p.requestChan <- 0
			return nil
		}
	}
	return errors.New("Peer has no pieces")
}
