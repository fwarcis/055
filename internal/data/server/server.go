// Package server...
package server

import (
	"context"
	"net"
	"sync"

	"055/internal/data/logging"
	"055/internal/data/statuses"
	"055/internal/data/stream"
)

// TODO: stream.Stream wrapper for communication logic with the protocol

func WatchErrors(cancel context.CancelFunc, errChan chan error) {
	for err := range errChan {
		logging.HandleForLogging(err)

		if stream.IsDisconnectCond(err) {
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
	sent := 0

	for !stream.IsDisconnectCond(err) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		packet, err = sender.Receive()
		if err != nil {
			errChan <- &logging.ReceivingError{Sender: sender, BaseErr: err}
			sent, err = sender.Send(string(statuses.Error), err.Error())
			if err != nil {
				errChan <- &logging.SendingError{Receiver: sender, Sent: sent, BaseErr: err}
			}
			continue
		}
		sent, err := sender.Send(string(statuses.Success), "message received")
		if err != nil {
			errChan <- &logging.SendingError{Receiver: sender, Sent: sent, BaseErr: err}
		}

		for i := range *receivers {
			currentStream := (*receivers)[i]
			header, body := (*packet).Header(), (*packet).Body()
			sent, err = currentStream.Send(header, body)
			if err != nil {
				errChan <- &logging.SendingError{
					Receiver: currentStream,
					Packet: *packet,
					Sent: sent,
					BaseErr: err,
				}
			}
		}
	}
}
