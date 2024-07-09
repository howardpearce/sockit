// Listens to a socket
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := os.Args[1]
	//msg := os.Args[2]

	//fmt.Printf("Sending message \"%s\" on socket: %s\n", msg, socketPath)

	//socket := must(net.Listen("unix", socketPath))

	conn := must(net.Dial("unix", socketPath))
	defer conn.Close()

	_, err := conn.Write([]byte("hello!"))

	if err != nil {
		fmt.Printf("Error occurred when writing: %w\n", err)
	}

}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
