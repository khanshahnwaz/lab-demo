package main

import (
	"fmt"
	"log"
	"time"

	"lab-demo/internal/network"
)

// ============================================================
// CHANGE ONLY THIS VALUE ON EACH COMPUTER
// ============================================================
//
// PC 01 -> 1
// PC 02 -> 2
// PC 03 -> 3
// PC 04 -> 4
// PC 05 -> 5
// PC 06 -> 6
// PC 07 -> 7
// PC 08 -> 8
// PC 09 -> 9
// PC 10 -> 10
// PC 11 -> 11
// PC 12 -> 12
// PC 13 -> 13
// PC 14 -> 14
// PC 15 -> 15
//
const THIS_NODE = 1

// ============================================================
// ALL NODE CONFIGURATIONS
// DO NOT CHANGE THESE AFTER VERIFYING THE IP LIST
// ============================================================

type NodeConfig struct {
	ID   string
	IP   string
	Port string
}

var nodes = []NodeConfig{
	{"node-01", "10.107.6.183", "5000"},
	{"node-02", "10.107.2.111", "5001"},
	{"node-03", "10.107.2.178", "5002"},
	{"node-04", "10.107.6.191", "5003"},
	{"node-05", "10.107.2.179", "5004"},
	{"node-06", "10.107.2.160", "5005"},
	{"node-07", "10.107.6.198", "5006"},
	{"node-08", "10.107.2.158", "5007"},
	{"node-09", "10.107.2.148", "5008"},
	{"node-10", "10.107.2.201", "5009"},
	{"node-11", "10.107.2.150", "5010"},
	{"node-12", "10.107.2.232", "5011"},
	{"node-13", "10.107.3.191", "5012"},
	{"node-14", "10.107.3.138", "5013"},
	{"node-15", "10.107.5.191", "5014"},
}

// ============================================================
// MAIN
// ============================================================

func main() {

	// --------------------------------------------------------
	// Validate node number
	// --------------------------------------------------------

	if THIS_NODE < 1 || THIS_NODE > len(nodes) {
		log.Fatalf(
			"Invalid THIS_NODE=%d. Must be between 1 and %d",
			THIS_NODE,
			len(nodes),
		)
	}

	// --------------------------------------------------------
	// Get configuration for THIS computer
	// --------------------------------------------------------

	myConfig := nodes[THIS_NODE-1]

	myAddress := myConfig.IP + ":" + myConfig.Port

	fmt.Println("========================================")
	fmt.Println("       P2P BLOCKCHAIN NETWORK NODE")
	fmt.Println("========================================")
	fmt.Println("Node ID :", myConfig.ID)
	fmt.Println("Address :", myAddress)
	fmt.Println("========================================")

	// --------------------------------------------------------
	// Create node
	// --------------------------------------------------------

	node := &network.Node{
		ID:      myConfig.ID,
		Address: myAddress,
	}

	// --------------------------------------------------------
	// Start TCP server
	// --------------------------------------------------------

	if err := node.Start(); err != nil {
		log.Fatal("Failed to start node:", err)
	}

	fmt.Println("Node started successfully.")
	fmt.Println()

	// --------------------------------------------------------
	// Connect to other nodes
	// --------------------------------------------------------
	//
	// Current topology:
	//
	// Each node connects to the previous two nodes.
	//
	// Node-01 -> none
	// Node-02 -> Node-01
	// Node-03 -> Node-01, Node-02
	// Node-04 -> Node-02, Node-03
	// Node-05 -> Node-03, Node-04
	// ...
	//
	// This avoids creating 15 x 14 connections.
	//
	// --------------------------------------------------------

	connectToPreviousNodes(node, THIS_NODE)

	fmt.Println()
	fmt.Println("P2P node is running.")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// --------------------------------------------------------
	// Keep program alive
	// --------------------------------------------------------

	for {
		time.Sleep(1 * time.Second)
	}
}

// ============================================================
// CONNECT TO PREVIOUS NODES
// ============================================================

func connectToPreviousNodes(node *network.Node, currentNode int) {

	// Node 01 has no previous node.
	if currentNode == 1 {
		fmt.Println("No outgoing peer connections required for Node-01.")
		return
	}

	// Connect to previous node.
	previousNode := nodes[currentNode-2]

	connectToNode(node, previousNode)

	// Connect to second previous node if available.
	if currentNode >= 3 {

		secondPreviousNode := nodes[currentNode-3]

		connectToNode(node, secondPreviousNode)
	}
}

// ============================================================
// CONNECT HELPER
// ============================================================

func connectToNode(node *network.Node, peer NodeConfig) {

	address := peer.IP + ":" + peer.Port

	fmt.Println("Connecting to", peer.ID, "at", address)

	err := node.Connect(address)

	if err != nil {
		fmt.Println(
			"Could not connect to",
			peer.ID,
			":",
			err,
		)

		return
	}

	fmt.Println(
		"Connected to",
		peer.ID,
		"at",
		address,
	)
}