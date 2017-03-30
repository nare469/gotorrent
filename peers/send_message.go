package peers

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func sendChoke(peerConn *PeerConnection) (err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, CHOKE)
	fmt.Println(buf.Bytes())
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amChoking = true
	return
}

func sendUnchoke(peerConn *PeerConnection) (err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, UNCHOKE)
	fmt.Println(buf.Bytes())
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amChoking = false
	return
}

func sendInterested(peerConn *PeerConnection) (err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, INTERESTED)
	fmt.Println(buf.Bytes())
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amInterested = true
	return
}

func sendUninterested(peerConn *PeerConnection) (err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, UNINTERESTED)
	fmt.Println(buf.Bytes())
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amInterested = false
	return
}

func sendRequest(peerConn *PeerConnection, index, begin, length uint32) (err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(13))
	binary.Write(buf, binary.BigEndian, REQUEST)
	binary.Write(buf, binary.BigEndian, index)
	binary.Write(buf, binary.BigEndian, begin)
	binary.Write(buf, binary.BigEndian, length)
	fmt.Println(buf.Bytes())
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amInterested = false
	return

}

func sendLoop(peerConn *PeerConnection) {
	sendInterested(peerConn)
	sendUnchoke(peerConn)

	for {
		select {
		case begin := <-peerConn.requestChan:
			fmt.Println("Sending request")
			fmt.Println(peerConn.pieceInfo.index)
			sendRequest(peerConn, peerConn.pieceInfo.index, begin, BLOCK_SIZE)

		}
	}
}
