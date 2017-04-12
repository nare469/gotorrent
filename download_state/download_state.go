package download_state

import (
	"crypto/sha1"
	"fmt"
	"github.com/nare469/gotorrent/parser"
	"io"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

const (
	MISSING byte = iota
	IN_PROGRESS
	COMPLETE
)

type state struct {
	pieces    []byte
	numPieces int
	mu        sync.RWMutex
	attrs     *parser.TorrentAttrs
}

var (
	s    *state
	once sync.Once
)

func InitDownloadState(attrs parser.TorrentAttrs) *state {
	numPieces, _ := attrs.NumPieces()
	once.Do(func() {
		s = &state{
			pieces: make([]byte, numPieces),
			attrs:  &attrs,
		}
		os.Mkdir("gotorrent_pieces", 0755)

	})
	return s
}

func GetPieceState(piece uint32) byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.pieces[piece]
}

func SetPieceState(piece uint32, state byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pieces[piece] = state
}

func WritePiece(data [][]byte, index uint32) (err error) {
	file, err := os.Create("gotorrent_pieces/piece_" + strconv.Itoa(int(index)))

	if err != nil {
		return
	}

	for _, value := range data {
		file.Write(value)
	}
	file.Close()

	verifyPiece(index)

	return
}

func verifyPiece(index uint32) {
	filePath := "gotorrent_pieces/piece_" + strconv.Itoa(int(index))
	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	h := sha1.New()

	if _, err := io.Copy(h, file); err != nil {
		return
	}
	hStr := h.Sum(nil)
	file.Close()

	hash, err := s.attrs.PieceHash()

	if err != nil {
		return
	}

	result := reflect.DeepEqual(hStr, hash[20*index:20*(index+1)])

	if result {
		SetPieceState(index, COMPLETE)
	} else {
		os.Remove(filePath)
		SetPieceState(index, MISSING)
	}
}

func completionWorker() {
	for {
		select {
		case <-time.After(10 * time.Second):
			for i := 0; i < len(s.pieces); i++ {
				if GetPieceState(uint32(i)) != COMPLETE {
					continue
				}
			}
			fmt.Println("DONE")
		}
	}
}
