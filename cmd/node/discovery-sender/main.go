package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Peer stores information about a discovered P2P node.
type Peer struct {
	ID   string // Peer/node identifier
	IP   string // Peer's IP address
	Port string // Peer's TCP listening port
}

func main() {

	// ------------------------------------------------
	// STEP 1: CREATE UDP SOCKET FOR NODE A
	// ------------------------------------------------

	// Node A will use UDP port 4001.
	//
	// 0.0.0.0 means the socket can use any available
	// IPv4 network interface on this machine.
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 4001,
	}

	// ------------------------------------------------
	// STEP 2: DEFINE LAN BROADCAST ADDRESS
	// ------------------------------------------------

	// Send the discovery packet to the broadcast address
	// of our current LAN.
	//
	// Nodes listening on UDP port 4000 can receive it.
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.2.15.255"),
		Port: 4000,
	}

	// ------------------------------------------------
	// STEP 3: START UDP SOCKET
	// ------------------------------------------------

	// Bind Node A to UDP port 4001.
	//
	// We use the same socket to:
	// 1. Send DISCOVER
	// 2. Receive discovery responses
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Println("Error starting UDP:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Node A listening on UDP port 4001")

	// Store all peers discovered during this discovery round.
	peers := []Peer{}

	// ------------------------------------------------
	// STEP 4: BROADCAST DISCOVERY MESSAGE
	// ------------------------------------------------

	message := []byte("DISCOVER")

	// Broadcast DISCOVER to UDP port 4000.
	_, err = conn.WriteToUDP(message, broadcastAddr)
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}

	fmt.Println("Discovery message sent!")

	// ------------------------------------------------
	// STEP 5: RECEIVE DISCOVERY RESPONSES
	// ------------------------------------------------

	// Buffer for incoming UDP responses.
	buffer := make([]byte, 1024)

	// Listen for responses for 5 seconds.
	//
	// After 5 seconds, ReadFromUDP will return an error
	// and we finish this discovery round.
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {

		// Wait for a peer to respond.
		//
		// n          = number of received bytes
		// remoteAddr = UDP address of responding peer
		// err        = possible network error
		n, remoteAddr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Finished discovery:", err)
			break
		}

		// Convert received bytes into a string.
		message := string(buffer[:n])

		fmt.Println(
			"Received:",
			message,
			"from:",
			remoteAddr,
		)

		// ------------------------------------------------
		// STEP 6: PARSE PEER INFORMATION
		// ------------------------------------------------

		// Expected response:
		//
		// NODE_B:5001
		//
		// NODE_B = peer ID
		// 5001   = peer's TCP listening port
		parts := strings.Split(message, ":")

		if len(parts) != 2 {
			fmt.Println("Invalid peer response")
			continue
		}

		// Create a Peer object.
		//
		// We get the IP from remoteAddr.
		// We get the ID and TCP port from the response.
		peer := Peer{
			ID:   parts[0],
			IP:   remoteAddr.IP.String(),
			Port: parts[1],
		}

		// Store the discovered peer.
		peers = append(peers, peer)

		fmt.Println("Discovered peer:", peer)
	}

	// ------------------------------------------------
	// STEP 7: CONNECT TO DISCOVERED PEERS USING TCP
	// ------------------------------------------------

	fmt.Println("\nDiscovered peers:")

	for _, peer := range peers {

		fmt.Println(peer)

		// Combine peer IP and TCP port.
		//
		// Example:
		// 10.2.12.231 + 5001
		//
		// becomes:
		// 10.2.12.231:5001
		tcpAddress := net.JoinHostPort(peer.IP, peer.Port)

		fmt.Println("Connecting to:", tcpAddress)

		// Establish TCP connection with the discovered peer.
		tcpConn, err := net.Dial("tcp", tcpAddress)
		if err != nil {
			fmt.Println("TCP connection failed:", err)
			continue
		}

		fmt.Println("TCP connection established with:", peer.ID)

		// ------------------------------------------------
		// STEP 8: SEND TCP MESSAGE
		// ------------------------------------------------

		message := []byte("Hello from Node A")

		_, err = tcpConn.Write(message)
		if err != nil {
			fmt.Println("Error sending TCP message:", err)
			tcpConn.Close()
			continue
		}

		fmt.Println("Message sent to:", peer.ID)

		// ------------------------------------------------
		// STEP 9: RECEIVE TCP RESPONSE
		// ------------------------------------------------

		// Create buffer for the TCP response.
		buffer := make([]byte, 1024)

		// Wait for the peer to respond.
		n, err := tcpConn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			tcpConn.Close()
			continue
		}

		fmt.Println(
			"Received from",
			peer.ID,
			":",
			string(buffer[:n]),
		)

		// ------------------------------------------------
		// STEP 10: CLOSE TCP CONNECTION
		// ------------------------------------------------

		// For now, we close the connection after
		// one message-response exchange.
		//
		// Later we can keep this connection alive
		// for persistent P2P communication.
		tcpConn.Close()
	}
}