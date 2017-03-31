package download_state

import (
	"os"
	"strconv"
	"sync"
)

const (
	MISSING byte = iota
	IN_PROGRESS
	COMPLETE
)

type state struct {
	pieces map[uint32]byte
	mu     sync.RWMutex
}

var (
	s    *state
	once sync.Once
)

func InitDownloadState() *state {
	once.Do(func() {
		s = &state{
			pieces: make(map[uint32]byte),
		}
		os.Mkdir("gotorrent_pieces", 0755)

	})
	return s
}

func GetPieceState(piece uint32) byte {
	return s.pieces[piece]
}

func SetPieceState(piece uint32, state byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pieces[piece] = state
}

func WritePiece(data [][]byte, index uint32) (err error) {
	file, err := os.Create("gotorrent_pieces/piece_" + strconv.Itoa(int(index)))
	defer file.Close()

	if err != nil {
		return
	}

	for _, value := range data {
		file.Write(value)
	}

	SetPieceState(index, COMPLETE)
	return
}
