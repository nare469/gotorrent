package connection

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/nare469/gotorrent/parser"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Connect(attrs parser.TorrentAttrs) (x parser.TrackerAttrs, err error) {
	u, err := url.Parse(attrs.Announce)
	if u.Scheme == "udp" {
		x, err = connectUdp(attrs)
	} else {
		x, err = connectHttp(attrs)
	}
	return
}

func connectHttp(attrs parser.TorrentAttrs) (x parser.TrackerAttrs, err error) {
	req, err := http.NewRequest("GET", attrs.Announce, nil)

	if err != nil {
		return
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
		return
	}
	q.Add("left", strconv.FormatInt(len, 10))
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())

	if err != nil {
		return
	}

	defer resp.Body.Close()

	// TODO: Remove logging of the string sometime soon
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	x, err = parser.NewTrackerAttrs(strings.NewReader(s))

	return

}

func connectUdp(attrs parser.TorrentAttrs) (x parser.TrackerAttrs, err error) {
	u, err := url.Parse(attrs.Announce)
	addr, err := net.ResolveUDPAddr("udp", u.Host)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// TODO: Retries
	packet := new(bytes.Buffer)
	var protocolId uint64 = 0x41727101980
	var action uint32 = 0
	transactionId := rand.Uint32()

	binary.Write(packet, binary.BigEndian, protocolId)
	binary.Write(packet, binary.BigEndian, action)
	binary.Write(packet, binary.BigEndian, transactionId)
	_, err = conn.Write(packet.Bytes())
	if err != nil {
		return
	}

	response := make([]byte, 16)

	_, err = conn.Read(response)

	responseBuffer := bytes.NewBuffer(response)

	var actionResponse uint32
	var transactionResponse uint32
	var connectionId uint64

	err = binary.Read(responseBuffer, binary.BigEndian, &actionResponse)
	err = binary.Read(responseBuffer, binary.BigEndian, &transactionResponse)
	err = binary.Read(responseBuffer, binary.BigEndian, &connectionId)

	if err != nil {
		return
	}

	announceRequest := new(bytes.Buffer)

	return

}

func populateUDPAnnounce(buf *bytes.Buffer) {
	// TODO: this
}
