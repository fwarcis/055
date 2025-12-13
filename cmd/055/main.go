package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"055/internal/protocol"
	"055/internal/stream"
)

func main() {
	if len(os.Args) >= 3 {
		fmt.Println("usage: 055 [ADDRESS]")
		os.Exit(1)
	}

	cfg := NewConfig()
	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()
	stm := stream.NewConnectionStream(
		conn, stream.EndOfPacket, stream.PacketPartsSeparator)

	go send(stm, conn.RemoteAddr().String())
	go receive(stm)

	for {
		time.Sleep(50 * time.Millisecond)
	}
}

func send(stm stream.Stream, remoteAddr string) {
	for {
		fmt.Print(remoteAddr + ": ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			continue
		}

		sent, err := stm.Send(stream.Packet{Header: string(protocol.Message), Body: input})
		if err != nil && err != io.EOF {
			log.Printf(err.Error()+" (%d sent)\n", sent)
		} else if err == io.EOF {
			log.Println("connection closed")
			return
		}
	}
}

func receive(stm stream.Stream) {
	for {
		packet, err := stm.Receive()
		time.Sleep(50 * time.Millisecond)
		if err != nil && err != io.EOF {
			log.Println(err.Error())
			continue
		} else if err == io.EOF {
			log.Println("connection closed")
			return
		}

		if packet.Header == string(protocol.Message) {
			fmt.Println("Message: ", packet.Body)
		}
	}
}
