package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port string = ":42069"

// the commented-out code is left here to visualize the developement of the course-work

func getLinesChannel(conn net.Conn) <-chan string {
	var lines = make(chan string)

	go func() {
		readBuf := make([]byte, 8)
		var lineBuf bytes.Buffer

		for {
			n, err := conn.Read(readBuf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					if lineBuf.Len() > 0 {
						lines <- lineBuf.String()
					}
					close(lines)
					fmt.Printf("connection (%s) has been closed\n", conn.RemoteAddr().String())
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(readBuf[:n])
			if strings.Contains(str, "\n") {
				i := strings.Index(str, "\n")
				lineBuf.WriteString(str[:i])
				lines <- lineBuf.String()

				lineBuf.Reset()
				lineBuf.WriteString(str[i+1:])
				continue
			}
			lineBuf.WriteString(str)
		}
	}()
	return lines
}

func main() {
	// refactor: file-reading -> tcp-connection
	// file, err := os.Open("messages.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("listening on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("connection accepted from: %s\n", conn.RemoteAddr().String())

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
	}

}

// refactor this
// func main() {
// 	var buf = make([]byte, 8)
// 	var lineBuf bytes.Buffer
//
// 	file, err := os.Open("messages.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
//
// 	for {
// 		n, err := file.Read(buf)
// 		if err != nil {
// 			if errors.Is(err, io.EOF) {
// 				break
// 			}
// 			fmt.Printf("error: %s\n", err.Error())
// 			break
// 		}
// 		str := string(buf[:n])
// 		if strings.Contains(str, "\n") {
// 			i := strings.Index(str, "\n")
// 			lineBuf.WriteString(str[:i])
// 			res := lineBuf.String()
// 			fmt.Printf("read: %s\n", res)
//
// 			lineBuf.Reset()
// 			lineBuf.WriteString(str[i+1:])
// 			continue
// 		}
// 		lineBuf.WriteString(str)
// 	}
// }
