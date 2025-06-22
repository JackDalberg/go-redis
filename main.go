package main

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAoF("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	err = aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Printf("Invalid command: %v\n", command)
			return
		}
		handler(args)
	})
	if err != nil {
		fmt.Printf("Could not read Aof: %v", err)
		return
	}

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
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Printf("Invalid command: %v\n", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if slices.Contains(ModifiesDB, command) {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
