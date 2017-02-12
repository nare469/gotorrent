package peers

import (
	"fmt"
	"github.com/nare469/gotorrent/parser"
	"net"
	"time"
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
	fmt.Println(peer.Port())
	fmt.Println(peer.IPAddr())
	conn, err := net.DialTimeout("udp", peer.HostName(), 10*time.Second)

	fmt.Println("Connecting to peer" + peer.HostName())

	if err != nil {
		quit <- false
		return
	}

	defer conn.Close()
	quit <- true
	return
}
