package parser

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type TorrentAttrs struct {
	announce string
	infoHash string
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
	attrs.announce = dict["announce"].(string)

	info := encode(dict["info"])
	fmt.Println(info)
	hash := sha1.New()
	hash.Write([]byte(info))

	attrs.infoHash = hex.EncodeToString(hash.Sum(nil))

	fmt.Println(attrs.infoHash)

	return
}

func (t TorrentAttrs) Announce() string {
	return t.announce
}
