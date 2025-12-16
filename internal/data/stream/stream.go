// Package stream...
package stream

import (
	"bufio"
	"net"
)

type Stream interface {
	Receive() (*Packet, error)
	Send(header, body string) (sent int, err error)
	Address() string
	Deserialize(header, body string) Packet
}

type ConnStream struct {
	conn          net.Conn
	reader        bufio.Reader
	writer        bufio.Writer
	headerBodySep string
	endOfPacket   byte
}

func NewConnStream(
	conn net.Conn,
	packetPartsSep string,
	endOfPacket byte,
) *ConnStream {
	return &ConnStream{
		conn:          conn,
		reader:        *bufio.NewReader(conn),
		writer:        *bufio.NewWriter(conn),
		headerBodySep: packetPartsSep,
		endOfPacket:   endOfPacket,
	}
}

func (strm *ConnStream) Receive() (*Packet, error) {
	content, err := strm.reader.ReadString(strm.endOfPacket)
	if err != nil {
		return nil, err
	}

	packet, err := Deserialize(content, strm.headerBodySep, strm.endOfPacket)
	if err != nil {
		return nil, err
	}

	return &packet, err
}

func (strm *ConnStream) Send(header, body string) (sent int, err error) {
	packet := validPacket{header, body, strm.headerBodySep, strm.endOfPacket}
	content := packet.Serialize()

	sent, err = strm.writer.WriteString(content)
	if err != nil {
		return sent, err
	}

	err = strm.writer.Flush()
	if err != nil {
		return sent, err
	}

	return sent, nil
}

func (strm *ConnStream) Deserialize(header, body string) Packet {
	return &validPacket{header, body, strm.headerBodySep, strm.endOfPacket}
}

func (strm *ConnStream) Address() string {
	return strm.conn.RemoteAddr().String()
}
