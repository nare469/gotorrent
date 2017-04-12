package download_state

import (
	"crypto/sha1"
	"github.com/nare469/gotorrent/logging"
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
		logging.Info.Println("Initializing download_state")
		s = &state{
			pieces: make([]byte, numPieces),
			attrs:  &attrs,
		}
		os.Mkdir("gotorrent_pieces", 0755)
		go completionWorker()

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
	logging.Info.Println("Writing piece", index, "to file")
	file, err := os.Create(pieceFileName(index))

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
	logging.Info.Println("Verifying piece", index)
	file, err := os.Open(pieceFileName(index))

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
		logging.Info.Println("Verification successful on piece", index)
		SetPieceState(index, COMPLETE)
	} else {
		logging.Error.Println("Verification failed on piece", index, ", removing file")
		os.Remove(pieceFileName(index))
		SetPieceState(index, MISSING)
	}
}

func completionWorker() {
	for {
		select {
		case <-time.After(10 * time.Second):
			logging.Info.Println("Completion Worker running")

			incomplete := false
			for i := 0; i < len(s.pieces); i++ {
				if GetPieceState(uint32(i)) != COMPLETE {
					incomplete = true
				}
			}
			if incomplete {
				continue
			}
			logging.Info.Println("Completion Worker detected fully downloaded file")
			mergePieces()
		}
	}
}

func mergePieces() {
	fileName, err := s.attrs.FileName()
	if err != nil {
		return
	}

	mergedFile, err := os.Create(fileName)

	for i := 0; i < len(s.pieces); i++ {
		pieceFile, err := os.Open(pieceFileName(uint32(i)))
		if err != nil {
			return
		}
		io.Copy(mergedFile, pieceFile)
		pieceFile.Close()
	}
	mergedFile.Close()
	return
}

func pieceFileName(index uint32) string {
	return "gotorrent_pieces/piece_" + strconv.Itoa(int(index))
}
