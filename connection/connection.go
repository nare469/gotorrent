package connection

import (
	"github.com/nare469/gotorrent/parser"
	"net"
	"net/url"
)

func Connect(attrs parser.TorrentAttrs) {
	trackerURL, err := url.Parse(attrs.Announce())
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial(trackerURL.Scheme, trackerURL.Host)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
}
