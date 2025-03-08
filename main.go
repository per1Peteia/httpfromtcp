package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	buf := make([]byte, 8)
	var lineBuf bytes.Buffer

	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		n, err := file.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		str := string(buf[:n])
		if strings.Contains(str, "\n") {
			i := strings.Index(str, "\n")
			lineBuf.WriteString(str[:i])
			res := lineBuf.String()
			fmt.Printf("read: %s\n", res)

			lineBuf.Reset()
			lineBuf.WriteString(str[i+1:])
			continue
		}
		lineBuf.WriteString(str)
	}
}
