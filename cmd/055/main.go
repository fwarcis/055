package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"055/internal/data/client"
	"055/internal/data/stream"
)

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Go(func() { client.WatchErrors(cancel, errChan) })

	if len(os.Args) >= 3 {
		fmt.Println("usage: 055 [ADDRESS]")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		errChan <- err
		close(errChan)
		wg.Wait()
		os.Exit(1)
	}
	defer conn.Close()

	stm := stream.NewConnStream(
		conn, stream.HeaderBodySep, stream.EndOfPacket)

	wg.Go(func() { client.RunReceiving(ctx, stm, errChan) })
	wg.Go(func() { client.RunSending(ctx, stm, errChan) })

	wg.Wait()
}
