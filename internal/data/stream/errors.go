package stream

import (
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
)

func IsDisconnectCond(err error) bool {
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

type ReceivingError struct {
	Sender  Stream
	BaseErr error
}

func (err *ReceivingError) Error() string {
	return err.Sender.Address() + ": " + err.BaseErr.Error()
}

type SendingError struct {
	Receiver Stream
	Packet   Packet
	Sent     int
	BaseErr  error
}

func (err *SendingError) Error() string {
	return fmt.Sprintf(
		"%s (%dB sent): %s: %s",
		err.Receiver.Address(),
		err.Sent,
		err.Packet.Serialize(),
		err.BaseErr.Error(),
	)
}
