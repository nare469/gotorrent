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
			continue
		}
		length := binary.BigEndian.Uint32(buf)
		if length != 0 {
			rest := make([]byte, length)
			io.ReadFull(peerConn.conn, rest)
			switch rest[0] {
			case CHOKE:
			case UNCHOKE:
			case INTERESTED:
			case UNINTERESTED:
			case HAVE:
				pieceIndex := binary.BigEndian.Uint32(rest[1:])
				fmt.Println("Have piece ", pieceIndex)
			case BITFIELD:
			case REQUEST:
			case PIECE:
			case CANCEL:
			}
		}

	}
}
