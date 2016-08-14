package main

import (
	"errors"
	"github.com/nare469/gotorrent/connection"
	"github.com/nare469/gotorrent/parser"
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

	x, err := parser.NewTorrentAttrs(file)
	if err != nil {
		panic(err)
	}

	connection.Connect(x)
}
