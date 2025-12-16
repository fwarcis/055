package main

import (
	"context"
	"net"
	"os"
	"sync"

	"055/internal/data/server"
)

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Go(func() { server.WatchErrors(cancel, errChan) })

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		errChan <- err
		close(errChan)
		wg.Wait()
		os.Exit(1)
	}
	defer listener.Close()

	wg.Go(func() { server.RunListening(ctx, listener, &wg, errChan) })

	wg.Wait()
}
