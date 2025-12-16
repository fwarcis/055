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

		if stream.IsAnyEOF(err) {
			cancel()
			return
		}
	}
}

func RunSending(ctx context.Context, stm stream.Stream, errChan chan error) {
	var err error
	input := ""

	for !stream.IsAnyEOF(err) {
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

		_, err = stm.Send("message", input)
		if err != nil {
			errChan <- err
		}
	}
}

func RunReceiving(ctx context.Context, stm stream.Stream, errChan chan error) {
	var err error
	var packet *stream.Packet

	for !stream.IsAnyEOF(err) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		packet, err = stm.Receive()
		if err != nil {
			errChan <- err
			continue
		}

		fmt.Println(stm.Address() + ": " + (*packet).Body())
	}
}

