package peers

import (
	"github.com/nare469/gotorrent/parser"
	"net"
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
	addr, err := net.ResolveUDPAddr("udp", peer.HostName())

	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		quit <- false
		return
	}

	defer conn.Close()
	quit <- true
	return
}
