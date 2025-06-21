package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error reading from client: %v", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte("+OK\r\n"))
	}
}
