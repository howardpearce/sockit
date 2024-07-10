// Listens to a socket
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := os.Args[1]
	msg := os.Args[2]

	fmt.Printf("Sending message \"%s\" on socket: %s\n", msg, socketPath)

	conn := must(net.Dial("unix", socketPath))
	defer conn.Close()

	_, err := conn.Write([]byte(msg))

	if err != nil {
		fmt.Printf("Error occurred when writing: %w\n", err)
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	fmt.Println("Received response:", string(buf[0:n]))

}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
