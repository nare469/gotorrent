package peers

import (
	"bytes"
	"encoding/binary"
	"github.com/nare469/gotorrent/logging"
	"time"
)

func sendChoke(peerConn *PeerConnection) (err error) {
	logging.Info.Println("Sending CHOKE message to", peerConn.peer.HostName())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, CHOKE)
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amChoking = true
	return
}

func sendUnchoke(peerConn *PeerConnection) (err error) {
	logging.Info.Println("Sending UNCHOKE message to", peerConn.peer.HostName())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, UNCHOKE)
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amChoking = false
	return
}

func sendInterested(peerConn *PeerConnection) (err error) {
	logging.Info.Println("Sending INTERESTED message to", peerConn.peer.HostName())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, INTERESTED)
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amInterested = true
	return
}

func sendUninterested(peerConn *PeerConnection) (err error) {
	logging.Info.Println("Sending UNINTERESTED message to", peerConn.peer.HostName())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, UNINTERESTED)
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amInterested = false
	return
}

func sendRequest(peerConn *PeerConnection, index, begin, length uint32) (err error) {
	logging.Info.Println("Requesting index", begin, "from piece", index, "from peer", peerConn.peer.HostName())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(13))
	binary.Write(buf, binary.BigEndian, REQUEST)
	binary.Write(buf, binary.BigEndian, index)
	binary.Write(buf, binary.BigEndian, begin)
	binary.Write(buf, binary.BigEndian, length)
	_, err = peerConn.conn.Write(buf.Bytes())
	if err != nil {
		return
	}
	peerConn.state.amInterested = false
	return

}

func sendKeepAlive(peerConn *PeerConnection) (err error) {
	logging.Info.Println("Sending KEEP-ALIVE message to peer", peerConn.peer.HostName())
	// Send [0 0 0 0] as a keep-alive message
	b := make([]byte, 4)
	peerConn.conn.Write(b)
	return
}

func sendLoop(peerConn *PeerConnection) {
	sendInterested(peerConn)
	sendUnchoke(peerConn)

	for {
		select {
		case begin := <-peerConn.requestChan:
			sendRequest(peerConn, peerConn.pieceInfo.index, begin, BLOCK_SIZE)
		case <-time.After(time.Minute * 2):
			sendKeepAlive(peerConn)
		}
	}
}
