// Package logging...
package logging

import (
	"fmt"
	"log"

	"055/internal/data/stream"
)

func HandleForLogging(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

type ReceivingError struct {
	Sender  stream.Stream
	BaseErr error
}

func (err *ReceivingError) Error() string {
	return err.Sender.Address() + ": " + err.BaseErr.Error()
}

type SendingError struct {
	Receiver stream.Stream
	Packet   stream.Packet
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
