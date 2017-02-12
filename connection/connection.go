package connection

import (
	"bytes"
	"encoding/binary"
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
	q.Add("left", strconv.FormatUint(len, 10))
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())

	if err != nil {
		return
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	x, err = parser.NewTrackerAttrsFromHttp(strings.NewReader(s))

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

	populateUDPAnnounce(connectionId, attrs, announceRequest)

	_, err = conn.Write(announceRequest.Bytes())

	announceResponse := make([]byte, 20+60)

	_, err = conn.Read(announceResponse)

	responseBuffer = bytes.NewBuffer(announceResponse)

	var interval uint32
	err = binary.Read(responseBuffer, binary.BigEndian, &actionResponse)
	err = binary.Read(responseBuffer, binary.BigEndian, &transactionResponse)
	err = binary.Read(responseBuffer, binary.BigEndian, &interval)

	x, err = parser.NewTrackerAttrsFromUdp(responseBuffer)

	return

}

func populateUDPAnnounce(connectionId uint64, attrs parser.TorrentAttrs, buf *bytes.Buffer) (err error) {
	var action uint32 = 1
	var transactionId uint32 = rand.Uint32()
	infoHash := string(attrs.InfoHash)
	peerId := string(attrs.PeerID)

	var downloaded uint64 = 0
	var left uint64
	left, err = attrs.Length()

	var uploaded uint64 = 0
	var event uint32 = 0
	var ipAddr uint32 = 0
	var key uint32 = 0
	var numWant uint32 = 10
	var port uint16 = 6881

	binary.Write(buf, binary.BigEndian, connectionId)
	binary.Write(buf, binary.BigEndian, action)
	binary.Write(buf, binary.BigEndian, transactionId)

	n, err := buf.WriteString(infoHash)
	if n != 20 || err != nil {
		return
	}
	n, err = buf.WriteString(peerId)
	if n != 20 || err != nil {
		return
	}

	// TODO: confirm transaction

	binary.Write(buf, binary.BigEndian, downloaded)
	binary.Write(buf, binary.BigEndian, left)
	binary.Write(buf, binary.BigEndian, uploaded)
	binary.Write(buf, binary.BigEndian, event)
	binary.Write(buf, binary.BigEndian, ipAddr)
	binary.Write(buf, binary.BigEndian, key)
	binary.Write(buf, binary.BigEndian, numWant)
	binary.Write(buf, binary.BigEndian, port)

	return
}
