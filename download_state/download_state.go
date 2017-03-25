package download_state

import (
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

func State() {
	once.Do(func() {
		s = &state{
			pieces: make(map[uint32]byte),
		}
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
