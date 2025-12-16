// Package client...
package client

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"055/internal/data/logging"
	"055/internal/data/stream"
)

func WatchErrors(cancel context.CancelFunc, errChan chan error) {
	for err := range errChan {
		logging.HandleForLogging(err)

		if stream.IsDisconnectCond(err) {
			cancel()
			return
		}
	}
}

func RunSending(ctx context.Context, stm stream.Stream, errChan chan error) {
	var err error
	input, sent := "", 0

	for !stream.IsDisconnectCond(err) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		input, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			errChan <- err
			continue
		}

		header := "message"
		body := input
		sent, err = stm.Send(header, body)
		if err != nil {
			errChan <- &stream.SendingError{
				Receiver: stm,
				Packet: stm.Deserialize(header, body),
				Sent: sent,
				BaseErr: err,
			}
		}
	}
}

func RunReceiving(ctx context.Context, stm stream.Stream, errChan chan error) {
	var err error
	var packet *stream.Packet

	for !stream.IsDisconnectCond(err) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		packet, err = stm.Receive()
		if err != nil {
			errChan <- &stream.ReceivingError{Sender: stm, BaseErr: err}
			continue
		}

		fmt.Println(stm.Address() + ": " + (*packet).Body())
	}
}

