package mysql

import (
	"fmt"
	"io"
	"net"
)

type PacketIO struct {
	Conn     net.Conn
	Sequence uint8
}

func (p *PacketIO) RemoteAddr() net.Addr {
	return p.Conn.RemoteAddr()
}

func (p *PacketIO) LocalAddr() net.Addr {
	return p.Conn.LocalAddr()
}

// ReadPacket loads []byte from Conn
func (p *PacketIO) ReadPacket() ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(p.Conn, header); err != nil {
		// log.Error("read header error %s", err.Error())
		return nil, ErrBadConn
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	if length < 1 {
		err := fmt.Errorf("invalid payload length %d", length)
		// log.Error(err.Error())
		return nil, err
	}

	sequence := header[3]
	if sequence != p.Sequence {
		err := fmt.Errorf("invalid sequence %d != %d", sequence, p.Sequence)
		// log.Error(err.Error())
		return nil, err
	}
	p.Sequence++

	data := make([]byte, length)
	if _, err := io.ReadFull(p.Conn, data); err != nil {
		// log.Error("read payload data error %s", err.Error())
		return nil, ErrBadConn
	}

	if length < MaxPayloadLen {
		return data, nil
	}

	var (
		buf []byte
		err error
	)
	buf, err = p.ReadPacket() // 递归调用
	if err != nil {
		// log.Error("read packet error %s", err.Error())
		return nil, ErrBadConn
	}
	return append(data, buf...), nil
}

// WritePacket sends []byte to Conn
// data already have header
func (c *PacketIO) WritePacket(data []byte) error {
	length := len(data) - 4
	for length >= MaxPayloadLen {
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff
		data[3] = c.Sequence
		n, err := c.Conn.Write(data[:4+MaxPayloadLen])
		if err != nil {
			// log.Error("write error %s", err.Error())
			return ErrBadConn
		}
		if n != (4 + MaxPayloadLen) {
			// log.Error("write error, write data number %d != %d", n, (4 + MaxPayloadLen))
			return ErrBadConn
		}
		c.Sequence++
		length -= MaxPayloadLen
		data = data[MaxPayloadLen:]
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)
	data[3] = c.Sequence
	n, err := c.Conn.Write(data)
	if err != nil {
		// log.Error("write error %s", err.Error())
		return ErrBadConn
	}
	if n != len(data) {
		// log.Error("write error, write data number %d != %d", n, (4 + MaxPayloadLen))
		return ErrBadConn
	}
	c.Sequence++
	return nil
}
