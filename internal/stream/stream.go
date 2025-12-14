// Package stream...
package stream

import (
	"bufio"
	"net"
	"strings"
)

type Packet struct {
	Header string
	Body   string
}

type Stream interface {
	Receive() (*Packet, error)
	Send(packet Packet) (sent int, err error)
	Close() error
}

const (
	EndOfPacket   = 0
	HeaderBodySep = "\t"
)

type ConnStream struct {
	conn          net.Conn
	reader        bufio.Reader
	writer        bufio.Writer
	endOfPacket   byte
	headerBodySep string
}

func NewConnectionStream(
	conn net.Conn,
	endOfPacket byte,
	packetPartsSep string,
) Stream {
	return &ConnStream{
		conn:          conn,
		reader:        *bufio.NewReader(conn),
		writer:        *bufio.NewWriter(conn),
		endOfPacket:   endOfPacket,
		headerBodySep: packetPartsSep,
	}
}

func (stm *ConnStream) Receive() (*Packet, error) {
	packetContent, err := stm.reader.ReadString(stm.endOfPacket)
	if err != nil {
		return nil, err
	}

	packetParts := strings.SplitN(packetContent, stm.headerBodySep, 2)
	if len(packetParts) != 2 {
		return nil, &WrongPacketFormatError{Content: packetContent}
	}

	endOfPacketIndex := len(packetParts[1]) - 1
	return &Packet{
		Header: packetParts[0],
		Body:   packetParts[1][:endOfPacketIndex],
	}, nil
}

func (stm *ConnStream) Send(packet Packet) (sent int, err error) {
	packetContent := (packet.Header + stm.headerBodySep +
		packet.Body + string(stm.endOfPacket))
	if packet.Header == "" || packet.Body == "" {
		return 0, &WrongPacketFormatError{Content: packetContent}
	}

	sent, err = stm.writer.WriteString(packetContent)
	if err != nil {
		return sent, err
	}
	err = stm.writer.Flush()
	if err != nil {
		return sent, err
	}
	return sent, nil
}

func (stm *ConnStream) Close() error {
	return stm.conn.Close()
}
