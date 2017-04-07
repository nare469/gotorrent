package main

import (
	"errors"
	"github.com/nare469/gotorrent/connection"
	"github.com/nare469/gotorrent/download_state"
	"github.com/nare469/gotorrent/parser"
	"github.com/nare469/gotorrent/peers"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		panic(errors.New("File argument needed"))
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	attrs, err := parser.NewTorrentAttrs(file)
	if err != nil {
		panic(err)
	}

	conn, err := connection.Connect(attrs)
	if err != nil {
		panic(err)
	}

	download_state.InitDownloadState(attrs)

	peers.ConnectToPeers(attrs, conn)
}
