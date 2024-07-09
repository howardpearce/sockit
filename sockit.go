// Listens to a socket
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
)

var block bool = false

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Please provide two paths as input.\n")
		os.Exit(1)
	}

	bindSocketPath := os.Args[1]
	dialSocketPath := os.Args[2]

	unlinkIfExists(bindSocketPath)
	unlinkIfExists(dialSocketPath)

	fmt.Println("Listening on sockets:", bindSocketPath, ",", dialSocketPath)

	bindSocket := must(net.Listen("unix", bindSocketPath))

	go func() {
		for {

			// Receive A message
			bindConn := must(bindSocket.Accept())
			buf := make([]byte, 4096)
			n := must(bindConn.Read(buf))
			fmt.Printf("Message Received size=%d: %s\n", n, string(buf[:]))

			if block {
				// Forward the message
				dialConn := must(net.Dial("unix", dialSocketPath))
				fmt.Printf("Writing message to other socket.\n")
				_, err := dialConn.Write(buf)
				if err != nil {
					fmt.Printf("Error occurred when writing: %w\n", err)
				}
				dialConn.Close()
			}

			bindConn.Close()
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		if strings.TrimRight(text, "\n") == "block" {
			fmt.Printf("Setting block to %t\n", block)
			block = !block
		}
	}

}

func unlinkIfExists(path string) {
	if _, err := os.Stat(path); err == nil {
		if err := syscall.Unlink(path); err != nil {
			fmt.Printf("Error occurred while unlinking existing socketfile: %w", err)
		}
	}

}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
