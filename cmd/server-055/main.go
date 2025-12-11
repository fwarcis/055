package main

import (
	"io"
	"log"
	"net"

	"055/internal/stream"
)

func main() {
	listener, err := net.Listen("tcp6", "localhost:6034")
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
