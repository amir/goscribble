package mpdclient

import (
	"net/textproto"
	"strings"
)

const (
	TCP = "tcp"
)

// Commands
const (
	CURRENTSONG = "currentsong"
	STATUS      = "status"
)

type Client struct {
	conn *textproto.Conn
}
type Values map[string]string

func Dial(address, password string) (c *Client, err error) {
	c = new(Client)

	if c.conn, err = textproto.Dial(TCP, address); err != nil {
		return
	}

	var version string
	if version, err = c.conn.ReadLine(); err != nil {
		return nil, err
	}

	if version[0:6] != "OK MPD" {
		return nil, textproto.ProtocolError("Handshake Error")
	}

	return
}

func (c *Client) readValues() (values Values, err error) {
	values = make(Values)
	for {
		line, err := c.conn.ReadLine()
		if err != nil {
			return nil, err
		}
		if line == "OK" {
			break
		}
		z := strings.Index(line, ":")
		if z < 0 {
			return nil, textproto.ProtocolError("Can't parse line: " + line)
		}
		key := line[0:z]
		values[key] = line[z+2:]
	}
	return
}

func (c *Client) CurrentSong() (s *Song, err error) {
	id, err := c.conn.Cmd(CURRENTSONG)
	if err != nil {
		return nil, err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	values, err := c.readValues()
	if err != nil {
		return
	}

	s = readCurrentSong(values)
	return
}

func (c *Client) Status() (s *Status, err error) {
	id, err := c.conn.Cmd(STATUS)
	if err != nil {
		return nil, err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	values, err := c.readValues()
	if err != nil {
		return
	}

	s = readStatus(values)
	return
}
