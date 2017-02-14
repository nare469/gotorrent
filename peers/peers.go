package peers

import (
	"bytes"
	"fmt"
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
	fmt.Println("Connection")
	fmt.Println(peer.HostName())
	addr, err := net.ResolveTCPAddr("tcp", peer.HostName())
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	fmt.Println("Connection Done")

	if err != nil {
		quit <- false
		return
	}

	conn.Write(createHandShake(toAttrs))

	header := make([]byte, 68)
	n, err := conn.Read(header)
	fmt.Println(n)
	fmt.Println(header)
	if err != nil {
		return
	}

	defer conn.Close()
	quit <- true
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
	fmt.Println(peerId)

	handshakeBuffer := new(bytes.Buffer)

	handshakeBuffer.WriteByte(byte(len(pstr)))
	handshakeBuffer.WriteString(pstr)
	handshakeBuffer.Write(reserved)
	handshakeBuffer.Write(infoHash)
	handshakeBuffer.WriteString(peerId)

	return handshakeBuffer.Bytes()
}
