package parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Peer struct {
	ipAddr []byte
	port   []byte
}

type TrackerAttrs struct {
	raw   Item
	Peers []Peer
}

func (p *Peer) IPAddr() []byte {
	return p.ipAddr
}

func (p *Peer) Port() []byte {
	return p.port
}

func (p *Peer) HostName() (s string) {
	s = ""
	for _, b := range p.ipAddr {
		s = s + strconv.Itoa(int(b))
		s = s + "."
	}
	s = s[:len(s)-1]

	s = s + ":"
	s = s + strconv.FormatInt(int64(binary.BigEndian.Uint16(p.port)), 10)

	return s
}

func NewTrackerAttrsFromUdp(buf *bytes.Buffer) (attrs TrackerAttrs, err error) {
	var leechers uint32
	var seeders uint32

	peerReader := bufio.NewReader(buf)

	err = binary.Read(peerReader, binary.BigEndian, &leechers)
	err = binary.Read(peerReader, binary.BigEndian, &seeders)

	for i := uint32(0); i < seeders+leechers; i++ {
		ipAddr := make([]byte, 4)
		port := make([]byte, 2)

		_, err = peerReader.Read(ipAddr)
		if err != nil {
			return
		}

		_, err = peerReader.Read(port)
		if err != nil {
			return
		}

		attrs.Peers = append(attrs.Peers, Peer{
			ipAddr: ipAddr,
			port:   port,
		})

	}

	return
}

func NewTrackerAttrsFromHttp(r io.Reader) (attrs TrackerAttrs, err error) {
	var item Item
	item, err = parse(r)

	if err != nil {
		return
	}

	attrs.raw = item

	err = attrs.createPeerList()
	if err != nil {
		return
	}

	return
}

func (a *TrackerAttrs) createPeerList() (err error) {
	dict, ok := a.raw.(map[string]Item)

	if !ok {
		err = errors.New("Could not parse response")
		return
	}

	peerStr, ok := dict["peers"].(string)

	if !ok {
		err = errors.New("Could not parse response")
		return
	}

	peerReader := strings.NewReader(peerStr)
	for {
		ipAddr := make([]byte, 4)
		port := make([]byte, 2)

		_, err = peerReader.Read(ipAddr)

		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		_, err = peerReader.Read(port)

		if err != nil {
			return
		}

		a.Peers = append(a.Peers, Peer{
			ipAddr: ipAddr,
			port:   port,
		})

	}

	return
}
