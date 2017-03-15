package peers

import (
	"bytes"
	"github.com/nare469/gotorrent/parser"
	"math/rand"
	"net"
	"strconv"
)

func ConnectToPeers(toAttrs parser.TorrentAttrs, trAttrs parser.TrackerAttrs) {

	quit := make(chan bool)
	for _, peer := range trAttrs.Peers {
		go connectToPeer(peer, toAttrs, quit)
	}

	for _, _ = range trAttrs.Peers {
		<-quit
	}
	return
}

func connectToPeer(peer parser.Peer, toAttrs parser.TorrentAttrs, quit chan bool) {
	addr, err := net.ResolveTCPAddr("tcp", peer.HostName())
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		quit <- false
		return
	}

	conn.Write(createHandShake(toAttrs))

	header := make([]byte, 68)
	_, err = conn.Read(header)
	// TODO: Verify handshake
	if err != nil {
		return
	}

	peerConn := NewPeerConnection(peer, conn)
	go sendLoop(peerConn)
	go receiveLoop(peerConn)

	// TODO: invesigate channel to channel
	quit <- <-peerConn.quitChan
	return
}

func createHandShake(attrs parser.TorrentAttrs) []byte {
	var pstr = "BitTorrent protocol"
	var reserved = make([]byte, 8)
	var infoHash = attrs.InfoHash
	var peerId = "NS0001-"
	for i := len(peerId); i < 20; i++ {
		peerId += strconv.Itoa(rand.Intn(10))
	}

	handshakeBuffer := new(bytes.Buffer)

	handshakeBuffer.WriteByte(byte(len(pstr)))
	handshakeBuffer.WriteString(pstr)
	handshakeBuffer.Write(reserved)
	handshakeBuffer.Write(infoHash)
	handshakeBuffer.WriteString(peerId)

	return handshakeBuffer.Bytes()
}
