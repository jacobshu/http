package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

const port = ":42069"

func c(s, c string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("error starting TCP listener: %v", err)
	}
	defer l.Close()

	fmt.Println("Listening to TCP traffic on", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)
		curr := ""

		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if curr != "" {
					ch <- curr
				}

				if errors.Is(err, io.EOF) {
					break
				}

				fmt.Printf("error: %v\n", err.Error())
				return
			}

			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s", curr, parts[i])
				curr = ""
			}
			curr += parts[len(parts)-1]
		}
	}()
	return ch
}
