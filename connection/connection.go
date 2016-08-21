package connection

import (
	"fmt"
	"github.com/nare469/gotorrent/parser"
	"io"
	"net/http"
	"os"
	"strconv"
)

func Connect(attrs parser.TorrentAttrs) {
	req, err := http.NewRequest("GET", attrs.Announce, nil)

	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("info_hash", string(attrs.InfoHash))
	q.Add("peer_id", string(attrs.PeerID))
	q.Add("port", "6881")
	q.Add("downloaded", "0")
	q.Add("uploaded", "0")
	q.Add("event", "started")
	q.Add("compact", "1")
	len, err := attrs.Length()

	if err != nil {
		panic(err)
	}
	q.Add("left", strconv.FormatInt(len, 10))
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := http.Get(req.URL.String())

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Println(resp)

	_, err = io.Copy(os.Stdout, resp.Body)

}
