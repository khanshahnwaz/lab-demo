package main

import (
	"fmt"
	"net"
)

func main() {

	// Define the UDP address on which this node will listen.
	// 0.0.0.0 means listen on all available IPv4 network interfaces.
	// Port 4000 is our UDP peer-discovery port.
	addr := net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 4000,
	}

	// Create a UDP socket and bind it to the above address.
	// This allows the node to receive UDP discovery messages.
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting discovery:", err)
		return
	}

	// Close the UDP connection when the program exits.
	defer conn.Close()

	fmt.Println("Listening for peer discovery on UDP port 4000...")

	// Create a buffer to store incoming UDP messages.
	// Maximum message size we will process here is 1024 bytes.
	buffer := make([]byte, 1024)

	// Keep listening continuously for discovery messages.
	for {

		// Read an incoming UDP packet.
		//
		// n          = number of bytes received
		// remoteAddr = address of the node that sent the packet
		// err        = possible error
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving:", err)
			return
		}

		// Convert only the received bytes into a string.
		message := string(buffer[:n])

		// Print the received message and the sender's address.
		fmt.Println(
			"Received:",
			message,
			"from:",
			remoteAddr,
		)

		// Check whether the received message is a peer-discovery request.
		if message == "DISCOVER" {

			// Send this node's identity and TCP port to the
			// node that requested discovery.
			//
			// NODE_A's_Duo = this node's ID
			// 5001          = TCP port used for P2P communication
			response := []byte("NODE_A's_Duo:5001")

			// Send the response directly back to the node
			// that sent the DISCOVER message.
			_, err = conn.WriteToUDP(response, remoteAddr)
			if err != nil {
				fmt.Println("Error sending response:", err)
				continue
			}

			// Confirm that the discovery response was sent.
			fmt.Println("Response sent to:", remoteAddr)
		}
	}
}