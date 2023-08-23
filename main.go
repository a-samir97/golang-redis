package main

import (
	"fmt"
	"net"
	"strings"
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
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		// check about type of the client input
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array type")
			continue
		}
		// check array length
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length greater than 0")
			continue
		}
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		result := handler(args)
		writer.Write(result)
	}
}
