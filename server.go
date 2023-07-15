package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	fmt.Println("Listening to port: 6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept and listen for connections
	conn, err := l.Accept()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from client ", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte("+ok\r\n"))
	}
}
