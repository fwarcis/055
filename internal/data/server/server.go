// Package server...
package server

import (
	"055/internal/data/logging"
	"055/internal/data/statuses"
	"055/internal/data/stream"
	"context"
	"net"
	"sync"
)

// TODO: stream.Stream wrapper for communication logic with the protocol

func WatchErrors(cancel context.CancelFunc, errChan chan error) {
	for err := range errChan {
		logging.HandleForLogging(err)

		if stream.IsAnyEOF(err) {
			cancel()
			return
		}
	}
}

func RunListening(
	ctx context.Context,
	listener net.Listener,
	wg *sync.WaitGroup,
	errChan chan error,
) {
	streams := []stream.Stream{}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
		}
		defer conn.Close()
		sender := stream.NewConnStream(
			conn, stream.HeaderBodySep, stream.EndOfPacket)

		wg.Go(func() { RunSharingWithOthers(ctx, sender, &streams, errChan) })

		streams = append(streams, sender)
	}
}

func RunSharingWithOthers(
	ctx context.Context,
	sender stream.Stream,
	receivers *[]stream.Stream,
	errChan chan error,
) {
	var err error
	var packet *stream.Packet

	for !stream.IsAnyEOF(err) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		packet, err = sender.Receive()
		if err != nil {
			sender.Send(string(statuses.Error), err.Error())
			errChan <- &logging.ReceivingError{Sender: sender, BaseErr: err}
			continue
		}
		sender.Send(string(statuses.Success), "message received")

		for i := range *receivers {
			currentStream := (*receivers)[i]
			header, body := (*packet).Header(), (*packet).Body()
			_, err = currentStream.Send(header, body)
			if err != nil {
				errChan <- &logging.ReceivingError{Sender: sender, BaseErr: err}
			}
		}
	}
}
