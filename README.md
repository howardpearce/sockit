## Sockit
`Sockit` is a small golang client used to forward messages along a unix socketfile. It creates a new socket, and forwards any messages that it receives to another socket. It also exposes a control socket used to issue commands to sockit iteself.

I named it `sockit` so that you can force a process to 'put a sock in it'.

## Why?
`Sockit` allows you to observe messages that are being sent by the service that is connected to it, without interrupting communication. It also allows you to test specific conditions which are hard to re-crate (ex. communication timeout on the domain socket) 

## How to use it
Sockit can be used as a command-line executable which takes three arugments:
1. A path to a socket to receive messages on
2. A path to a socket to forward messages to
3. A path to a control socket where commands can be received

example:
socket <path_to_receive_socket> <path_to_send_socket> <path_to_control_socket>

Once `sockit` is running, you can send the word 'block' over the control socket and it will stop forwarding messages to simulate a communication timeout.
