package main

import (
	"fmt"
	"io"
	"os"
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

func c(s, c string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()

		var curr string
		b := make([]byte, 8)
		var e error
		for e != io.EOF {
			_, e = f.Read(b)
			parts := strings.Split(fmt.Sprintf("%s", b), "\n")
			clear(b)

			if len(parts) == 1 {
				curr = fmt.Sprintf("%s%s", curr, parts[0])
			} else {
				for i, p := range parts {
					if i == len(parts)-1 {
						ch <- curr
						curr = p
					} else {
						curr = fmt.Sprintf("%s%s", curr, p)
					}
				}
			}
		}
	}()

	return ch
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening messages: %v", err)
	}

	ch := getLinesChannel(f)
	for c := range ch {
		fmt.Printf("read: %s\n", c)
	}
}
