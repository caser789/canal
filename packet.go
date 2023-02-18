package main

import (
	"bufio"
	"bytes"
	"github.com/pingcap/errors"
	"io"
	"net"
)

func NewPacket(conn net.Conn) *Packet {
	c := new(Packet)
	c.Conn = conn

	c.br = bufio.NewReaderSize(c, 65536) // 64kb
	c.reader = c.br

	c.copyNBuf = make([]byte, 16*1024)

	return c
}

type Packet struct {
	net.Conn

	br     *bufio.Reader
	reader io.Reader

	header   [4]byte
	copyNBuf []byte

	Sequence uint8
}

// WritePacket data already has 4 bytes header
// will modify data inplace
func (p *Packet) WritePacket(data []byte) error {
	length := len(data) - 4

	for length >= MaxPayloadLen {
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff
		data[3] = p.Sequence

		n, err := p.Write(data[:4+MaxPayloadLen])
		if err != nil {
			return errors.Wrapf(ErrBadConn, "Write(payload portion) failed. err %v", err)
		}

		if n != (4 + MaxPayloadLen) {
			return errors.Wrapf(ErrBadConn, "Write(payload portion) failed. only %v bytes written, while %v expected", n, 4+MaxPayloadLen)
		}

		p.Sequence++
		length -= MaxPayloadLen
		data = data[MaxPayloadLen:]
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)
	data[3] = p.Sequence

	n, err := p.Write(data)
	if err != nil {
		return errors.Wrapf(ErrBadConn, "Write failed. err %v", err)
	}

	if n != len(data) {
		return errors.Wrapf(ErrBadConn, "Write failed. only %v bytes written, while %v expected", n, len(data))
	}

	p.Sequence++
	return nil
}

func (p *Packet) ReadPacket() ([]byte, error) {
	return p.ReadPacketReuseMem(nil)
}

func (p *Packet) ReadPacketReuseMem(dst []byte) ([]byte, error) {
	// Here we use `sync.Pool` to avoid allocate/destroy buffers frequently.
	buf := GetBytesBuffer()
	defer func() {
		PutBytesBuffer(buf)
	}()

	if err := p.ReadPacketTo(buf); err != nil {
		return nil, errors.Trace(err)
	}

	readBytes := buf.Bytes()
	readSize := len(readBytes)
	var result []byte
	if len(dst) > 0 {
		result = append(dst, readBytes...)
		// if read block is big, do not cache buf anymore
		if readSize > TooBigBlockSize {
			buf = nil
		}
	} else {
		if readSize > TooBigBlockSize {
			// if read block is big, use read block as result and do not cache buf anymore
			result = readBytes
			buf = nil
		} else {
			result = append(dst, readBytes...)
		}
	}

	return result, nil
}

// ReadPacketTo read body to w (without header)
func (p *Packet) ReadPacketTo(w io.Writer) error {
	// [1] read header
	if _, err := io.ReadFull(p.reader, p.header[:4]); err != nil {
		return errors.Wrapf(ErrBadConn, "io.ReadFull(header) failed. err %v", err)
	}

	length := int(uint32(p.header[0]) | uint32(p.header[1])<<8 | uint32(p.header[2])<<16)
	sequence := p.header[3]

	if sequence != p.Sequence {
		return errors.Errorf("invalid sequence %d != %d", sequence, p.Sequence)
	}

	p.Sequence++

	if buf, ok := w.(*bytes.Buffer); ok {
		// Allocate the buffer with expected length directly instead of call `grow` and migrate data many times.
		buf.Grow(length)
	}

	// [2] read body
	n, err := p.copyN(w, p.reader, int64(length))
	if err != nil {
		return errors.Wrapf(ErrBadConn, "io.CopyN failed. err %v, copied %v, expected %v", err, n, length)
	}

	if n != int64(length) {
		return errors.Wrapf(ErrBadConn, "io.CopyN failed(n != int64(length)). %v bytes copied, while %v expected", n, length)
	}

	// [3] there are more
	if length < MaxPayloadLen {
		return nil
	}

	if err := p.ReadPacketTo(w); err != nil {
		return errors.Wrap(err, "ReadPacketTo failed")
	}

	return nil
}

func (p *Packet) copyN(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
	for n > 0 {
		bcap := cap(p.copyNBuf)
		if int64(bcap) > n {
			bcap = int(n)
		}
		buf := p.copyNBuf[:bcap]

		rd, err := io.ReadAtLeast(src, buf, bcap)
		n -= int64(rd)

		if err != nil {
			return written, errors.Trace(err)
		}

		wr, err := dst.Write(buf)
		written += int64(wr)
		if err != nil {
			return written, errors.Trace(err)
		}
	}

	return written, nil
}

func (c *Packet) Close() error {
	c.Sequence = 0
	if c.Conn != nil {
		return errors.Wrap(c.Conn.Close(), "Conn.Close failed")
	}
	return nil
}
