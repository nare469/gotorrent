package parser

import (
	"crypto/sha1"
	"errors"
	"io"
	"math"
)

// TODO: Rename to torrent_attributes

type TorrentAttrs struct {
	Announce string
	InfoHash []byte
	PeerID   []byte
	raw      Item
}

func NewTorrentAttrs(r io.Reader) (attrs TorrentAttrs, err error) {
	var item Item
	item, err = parse(r)

	if err != nil {
		return
	}

	attrs.raw = item

	dict, ok := item.(map[string]Item)
	if !ok {
		err = errors.New("Invalid torrent file")
	}
	attrs.Announce = dict["announce"].(string)

	info := encode(dict["info"])
	hash := sha1.New()
	hash.Write([]byte(info))

	attrs.InfoHash = hash.Sum(nil)
	attrs.PeerID = []byte("Narendran Srinivasan")

	return
}

func (me *TorrentAttrs) FileName() (name string, err error) {
	dict, ok := me.raw.(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	info, ok := dict["info"].(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	if val, ok := info["name"]; ok {
		name, ok = val.(string)
		if !ok {
			err = errors.New("Invalid Torrent File")
		}
		return
	}

	err = errors.New("Multiple files aren't supported yet")
	return
}

func (me *TorrentAttrs) Length() (length uint64, err error) {
	// TODO: Make this uniform i.e. functions for all or data members for all
	dict, ok := me.raw.(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	info, ok := dict["info"].(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	if val, ok := info["length"]; ok {
		x, ok := val.(int64)
		if !ok {
			err = errors.New("Invalid Torrent File")
		}
		length = uint64(x)
		return
	}

	files, ok := info["files"]

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	length = 0

	filesArr, ok := files.([]Item)

	for _, file := range filesArr {
		fileDict, ok := file.(map[string]Item)

		if !ok {
			err = errors.New("Invalid torrent file")
		}

		fileSize, ok := fileDict["length"].(uint64)

		if !ok {
			err = errors.New("Invalid torrent file")
		}

		length += fileSize
	}
	return
}

func (me *TorrentAttrs) PieceLength() (length uint64, err error) {
	dict, ok := me.raw.(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	info, ok := dict["info"].(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	if val, ok := info["piece length"]; ok {
		x, ok := val.(int64)
		if !ok {
			err = errors.New("Invalid Torrent File")
		}
		length = uint64(x)
	}
	return
}

func (me *TorrentAttrs) NumPieces() (num uint64, err error) {
	length, err := me.Length()
	pieceLength, err := me.PieceLength()
	if err != nil {
		return
	}
	num = uint64(math.Ceil(float64(length) / float64(pieceLength)))
	return
}

func (me *TorrentAttrs) PieceHash() (hash []byte, err error) {
	dict, ok := me.raw.(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	info, ok := dict["info"].(map[string]Item)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	hashItem, ok := info["pieces"]

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	hashString, ok := hashItem.(string)

	hash = []byte(hashString)

	if !ok {
		err = errors.New("Invalid torrent file")
	}

	return
}
