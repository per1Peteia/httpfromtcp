package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const addr string = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		data, err := buf.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatal(err)
		}
	}
}
