// Listens to a socket
package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	socketPath := os.Args[1]

	if _, err := os.Stat(socketPath); err == nil {
		if err := syscall.Unlink(socketPath); err != nil {
			fmt.Printf("Error occurred while unlinking existing socketfile: %w", err)
		}
	}

	fmt.Println("Listening on socket:", socketPath)

	socket := must(net.ListenUnix("unix", &net.UnixAddr{socketPath, "unix"}))

	for {
		conn := must(socket.Accept())

		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 4096)

			n := must(conn.Read(buf))
			fmt.Printf("Message Received size=%d: %s\n", n, string(buf[:]))
			conn.Write([]byte("OK"))
		}(conn)
	}

}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
