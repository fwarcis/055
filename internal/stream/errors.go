package stream

import "fmt"

type WrongPacketFormatError struct {
	Content string
}

func (err *WrongPacketFormatError) Error() string {
	return fmt.Sprintf("conn stream: \"%s\": wrong packet format", err.Content)
}
