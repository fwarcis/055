package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"055/internal/config"
	"055/internal/stream"
)

func main() {
	network, address := config.NETWORK, config.ADDRESS
	if len(os.Args) == 2 {
		network = os.Args[1]
	} else if len(os.Args) == 3 {
		network = os.Args[1]
		address = os.Args[2]
	} else if len(os.Args) > 3 {
		fmt.Println("usage: 055 [NETWORK] [ADDRESS]")
		os.Exit(1)
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer listener.Close()

	streams := []stream.Stream{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		stm := stream.NewConnectionStream(
			conn, stream.EndOfPacket, stream.PacketPartsSeparator)
		streams = append(streams, stm)

		go distribute(streams, stm)

	}
}

func distribute(streams []stream.Stream, stm stream.Stream) {
	defer stm.Close()
	for {
		packet, err := stm.Receive()
		if err != nil && err != io.EOF {
			log.Println(err.Error())
			continue
		} else if err == io.EOF {
			log.Println("connection closed")
			return
		}

		for i := range streams {
			if streams[i] == stm {
				continue
			}
			sent, err := streams[i].Send(*packet)
			if err != nil {
				log.Printf(err.Error()+" (%d bytes sent)\n", sent)
			}
		}
	}
}
