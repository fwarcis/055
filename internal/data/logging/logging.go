// Package logging...
package logging

import (
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
