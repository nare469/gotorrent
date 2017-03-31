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

	// TODO: Rename x and y
	x, err := parser.NewTorrentAttrs(file)
	if err != nil {
		panic(err)
	}

	y, err := connection.Connect(x)
	if err != nil {
		panic(err)
	}

	download_state.InitDownloadState()

	peers.ConnectToPeers(x, y)
}
