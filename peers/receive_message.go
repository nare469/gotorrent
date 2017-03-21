package peers

import (
	"encoding/binary"
	"fmt"
	"io"
)

func receiveLoop(peerConn *PeerConnection) {
	buf := make([]byte, 4)
	for {
		_, err := io.ReadFull(peerConn.conn, buf)
		if err != nil {
			fmt.Println(err.Error())
			peerConn.quitChan <- true
		}
		length := binary.BigEndian.Uint32(buf)
		if length != 0 {
			rest := make([]byte, length)
			io.ReadFull(peerConn.conn, rest)
			switch rest[0] {
			case CHOKE:
				fmt.Println("CHOKE")
				peerConn.state.peerChoking = true
			case UNCHOKE:
				fmt.Println("UNCHOKE")
				peerConn.state.peerChoking = false
			case INTERESTED:
				fmt.Println("INTERESTED")
				peerConn.state.peerInterested = true
			case UNINTERESTED:
				fmt.Println("uninter")
				peerConn.state.peerInterested = false
			case HAVE:
				pieceIndex := binary.BigEndian.Uint32(rest[1:])
				fmt.Println("HAVE ", pieceIndex)
				peerConn.haveChan <- pieceIndex
			case BITFIELD:
				fmt.Println("BITFIELD")
				peerConn.bitfield = rest
			case REQUEST:
				fmt.Println("REQUEST")
			case PIECE:
				fmt.Println("PIECE")
			case CANCEL:
				fmt.Println("CANCEL")
			}
		}

	}
}
