package stream

import (
	"strings"
)

const (
	EndOfPacket   = 0
	HeaderBodySep = "\x1e"
)

type Packet interface {
	Header() string
	Body() string
	Serialize() string
}

type validPacket struct {
	header        string
	body          string
	headerBodySep string
	endOfPacket   byte
}

func Deserialize(
	content string,
	headerBodySep string,
	endOfPacket byte,
) (Packet, error) {
	packetParts := strings.SplitN(content, headerBodySep, 2)
	if len(packetParts) != 2 {
		return nil, &WrongPacketFormatError{Content: content}
	}

	endOfPacketIndex := len(packetParts[1]) - 1
	return &validPacket{
		header: packetParts[0],
		body:   packetParts[1][:endOfPacketIndex],
	}, nil
}

func (pack *validPacket) Header() string {
	return pack.header
}

func (pack *validPacket) Body() string {
	return pack.body
}

func (pack *validPacket) Serialize() string {
	return pack.header + pack.headerBodySep +
		pack.body + string(pack.endOfPacket)
}
