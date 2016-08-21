package parser

import (
	"errors"
	"io"
	"strings"
)

type peer struct {
	ipAddr []byte
	port   []byte
}

type TrackerAttrs struct {
	raw   Item
	peers []peer
}

func (p *peer) HostName() string {
	//TODO: this
	return ""
}

func NewTrackerAttrs(r io.Reader) (attrs TrackerAttrs, err error) {
	var item Item
	item, err = parse(r)

	if err != nil {
		return
	}

	attrs.raw = item

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

		a.peers = append(a.peers, peer{
			ipAddr: ipAddr,
			port:   port,
		})

	}

	return
}
