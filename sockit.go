// Listens to a socket
package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

var block bool = false

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Please provide three paths as input.\n")
		os.Exit(1)
	}

	receiveSocketPath := os.Args[1]
	sendSocketPath := os.Args[2]
	controlSocketPath := os.Args[3]

	unlinkIfExists(receiveSocketPath)
	unlinkIfExists(controlSocketPath)

	fmt.Println("Receiving messages on", receiveSocketPath, "and forwarding them to", sendSocketPath)

	receiveSocket := must(net.Listen("unix", receiveSocketPath))
	controlSocket := must(net.Listen("unix", controlSocketPath))

	go func() {
		for {

			// Receive a message
			receiveConn := must(receiveSocket.Accept())
			buf := make([]byte, 4096)
			n := must(receiveConn.Read(buf))
			fmt.Printf("Message Received size=%d: %s\n", n, string(buf[:]))

			if !block {
				sendSocket := must(net.Dial("unix", sendSocketPath))

				// Forward the message
				fmt.Printf("Writing message to other socket.\n")
				_, err := sendSocket.Write(buf)
				if err != nil {
					fmt.Printf("Error occurred when writing: %w\n", err.Error())
				}

				sendBuf := make([]byte, 4096)
				n, err := sendSocket.Read(sendBuf)
				fmt.Println("Received response:", string(sendBuf[0:n]))
				fmt.Println("Returning response to receiver")
				receiveConn.Write(sendBuf[0:n])

				sendSocket.Close()

			} else {
				fmt.Println("Block enabled. Not forwarding message.")
			}

			receiveConn.Close()

		}
	}()

	for {
		controlConn := must(controlSocket.Accept())
		buf := make([]byte, 4096)
		n := must(controlConn.Read(buf))
		fmt.Printf("Message Received on control socket size=%d: %s\n", n, string(buf[:]))

		if string(buf[0:n]) == "block" {
			block = !block
		}

		controlConn.Close()
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
