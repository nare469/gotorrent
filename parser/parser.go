package parser

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type Item interface{}

func parseInt(buf *bufio.Reader, delim byte) (data int64, err error) {
	var str []byte
	for {
		var c byte
		c, err = buf.ReadByte()

		if err != nil {
			return
		}

		if c == delim {
			break
		}

		if !(c >= '0' && c <= '9') {
			err = errors.New("Invalid character in integer")
			return
		}

		str = append(str, c)
	}

	data, err = strconv.ParseInt(string(str), 10, 64)
	return
}

func parseString(buf *bufio.Reader) (data string, err error) {

	length, err := parseInt(buf, ':')

	if err != nil {
		return
	}

	str := make([]byte, length)

	_, err = io.ReadFull(buf, str)

	if err != nil {
		return
	}

	data = string(str)
	return
}

func parseFromReader(buf *bufio.Reader) (it Item, err error) {
	var c byte
	c, err = buf.ReadByte()
	if err != nil {
		return
	}
	switch {
	case c >= '0' && c <= '9':
		err = buf.UnreadByte()
		if err != nil {
			return
		}

		it, err = parseString(buf)
		if err != nil {
			return
		}
	case c == 'i':
		it, err = parseInt(buf, 'e')
	case c == 'l':
		list := make([]Item, 0, 8)

		for {
			c, err = buf.ReadByte()
			if err != nil {
				return
			}

			if c == 'e' {
				break
			}

			err = buf.UnreadByte()
			if err != nil {
				return
			}

			var el Item
			el, err = parseFromReader(buf)
			if err != nil {
				return
			}

			list = append(list, el)
		}
		it = list
	case c == 'd':
		dict := make(map[string]Item)

		for {
			c, err = buf.ReadByte()
			if err != nil {
				return
			}

			if c == 'e' {
				break
			}

			err = buf.UnreadByte()
			if err != nil {
				return
			}

			var key string
			var value Item
			key, err = parseString(buf)
			if err != nil {
				return
			}

			value, err = parseFromReader(buf)
			if err != nil {
				return
			}

			dict[key] = value
		}
		it = dict
	}
	return
}

func parse(r io.Reader) (it Item, err error) {
	buf := newBufioReader(r)
	return parseFromReader(buf)
}

func newBufioReader(r io.Reader) *bufio.Reader {
	return bufio.NewReader(r)
}
