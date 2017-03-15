package peers

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func makeInterestedMessage() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, INTERESTED)
	fmt.Println(buf.Bytes())
	return buf.Bytes()
}

func makeUnchokeMessage() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, UNCHOKE)
	fmt.Println(buf.Bytes())
	return buf.Bytes()
}

func sendLoop(peerConn *PeerConnection) {
	go receiveLoop(peerConn)
	peerConn.conn.Write(makeInterestedMessage())
	peerConn.conn.Write(makeUnchokeMessage())
}
