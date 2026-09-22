package main

import (
	"fmt"
	"net"
)

func main() {

	// Start a TCP server and listen on port 5001.
	//
	// 0.0.0.0 means the server accepts connections
	// coming through any IPv4 network interface.
	listener, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}

	// Close the TCP listener when the program exits.
	defer listener.Close()

	fmt.Println("Node A's duo TCP server listening on port 5001")

	// Keep accepting TCP connections from peers.
	for {

		// Wait until another node connects to this TCP server.
		//
		// conn represents the TCP connection with that peer.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Print the address of the peer that connected.
		fmt.Println("TCP connection received from:", conn.RemoteAddr())

		// Create a buffer for receiving data from the peer.
		buffer := make([]byte, 1024)

		// Read data sent by the peer.
		//
		// n = number of bytes received.
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			conn.Close()
			continue
		}

		// Convert the received bytes into a string and print it.
		fmt.Println("Received message:", string(buffer[:n]))

		// Prepare a response to send back to the peer.
		response := []byte("Hello from Node B")

		// Send the response through the same TCP connection.
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error sending response:", err)
			conn.Close()
			continue
		}

		fmt.Println("Response sent to Node A")

		// Close the TCP connection.
		//
		// Currently we close the connection after one
		// request-response exchange.
		conn.Close()
	}
}