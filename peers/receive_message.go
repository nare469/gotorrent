package peers

import (
	"encoding/binary"
	"github.com/nare469/gotorrent/logging"
	"io"
)

func receiveLoop(peerConn *PeerConnection) {
	buf := make([]byte, 4)
	for {
		_, err := io.ReadFull(peerConn.conn, buf)
		if err != nil {
			peerConn.quitChan <- true
		}
		length := binary.BigEndian.Uint32(buf)
		if length != 0 {
			rest := make([]byte, length)
			io.ReadFull(peerConn.conn, rest)
			switch rest[0] {
			case CHOKE:
				logging.Info.Println("Received CHOKE from", peerConn.peer.HostName())
				peerConn.state.peerChoking = true
			case UNCHOKE:
				logging.Info.Println("Received UNCHOKE from", peerConn.peer.HostName())
				peerConn.state.peerChoking = false
			case INTERESTED:
				logging.Info.Println("Received INTERESTED from", peerConn.peer.HostName())
				peerConn.state.peerChoking = false
				peerConn.state.peerInterested = true
			case UNINTERESTED:
				logging.Info.Println("Received UNINTERESTED from", peerConn.peer.HostName())
				peerConn.state.peerInterested = false
			case HAVE:
				pieceIndex := binary.BigEndian.Uint32(rest[1:])
				logging.Info.Println("Received HAVE from", peerConn.peer.HostName(), "with piece", pieceIndex)
				peerConn.setHasPiece(pieceIndex)
				peerConn.canReceiveBitfield = false
			case BITFIELD:
				logging.Info.Println("Received BTFIELD from", peerConn.peer.HostName())
				if peerConn.canReceiveBitfield {
					peerConn.setBitfield(rest[1:])
				}
				peerConn.canReceiveBitfield = false
				peerConn.choosePieceToRequest()
			case REQUEST:
				logging.Info.Println("Received REQUEST from", peerConn.peer.HostName())
			case PIECE:
				pieceIndex := binary.BigEndian.Uint32(rest[1:5])
				offset := binary.BigEndian.Uint32(rest[5:9]) / BLOCK_SIZE
				logging.Info.Println("Received PIECE", pieceIndex, "at position", offset, "from", peerConn.peer.HostName())
				if pieceIndex == peerConn.pieceInfo.index && offset == peerConn.pieceInfo.counter {
					peerConn.receiveBlock(rest[9:])
				}
			case CANCEL:
				logging.Info.Println("Received CANCEL from", peerConn.peer.HostName())
			}
		}

	}
}
