package stream

import (
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
)

func IsAnyEOF(err error) bool {
	return err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, net.ErrClosed) || errors.Is(err, net.ErrWriteToConnected) ||
		errors.Is(err, syscall.EPIPE)
}

type WrongPacketFormatError struct {
	Content string
}

func (err *WrongPacketFormatError) Error() string {
	return fmt.Sprintf("stream: \"%q\": wrong packet format", err.Content)
}
