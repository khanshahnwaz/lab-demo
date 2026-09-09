package main

import (
	"fmt"
	"lab-demo/internal/network"
	"time"
)

func main() {
	node := &network.Node{
		ID:      "node-C",
		Address: ":5002",
	}

	fmt.Println("Starting Node C...")

	if err := node.Start(); err != nil {
		fmt.Println("Node failed:", err)
		return
	}

	time.Sleep(2 * time.Second)

	if err := node.Connect("10.107.3.191:5000"); err != nil {
		fmt.Println("Connection to Node A failed:", err)
	}

	if err := node.Connect("10.107.2.151:5001"); err != nil {
		fmt.Println("Connection to Node B failed:", err)
	}

	select {}
}
