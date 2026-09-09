# P2P Network — Run & Deployment Guide

## 1. Purpose

This document explains how to deploy, run, and test the Go-based P2P networking system across multiple physical computers connected to the same LAN.

The current system is designed so that:

* Every computer runs one P2P node.
* Every node can act as both a TCP server and TCP client.
* Nodes communicate using TCP.
* Messages use JSON format.
* Messages use newline-delimited framing.
* Peers are maintained by `PeerManager`.
* Nodes can broadcast messages.
* PING/PONG heartbeat is used to monitor peers.
* Disconnected peers are detected and removed.

The networking implementation is common across all computers.

Only:

```text
cmd/node/main.go
```

is different for each node.

---

# 2. Network Setup

Each physical computer must be connected to the same LAN.

Example:

```text
PC 01 → Node 01
PC 02 → Node 02
PC 03 → Node 03
...
PC 16 → Node 16
```

Each node must have:

* A unique LAN IPv4 address
* A unique TCP port
* A unique Node ID

Example:

| PC    | Node ID | Port |
| ----- | ------- | ---: |
| PC 01 | node-01 | 5000 |
| PC 02 | node-02 | 5001 |
| PC 03 | node-03 | 5002 |
| PC 04 | node-04 | 5003 |
| PC 05 | node-05 | 5004 |
| PC 06 | node-06 | 5005 |
| PC 07 | node-07 | 5006 |
| PC 08 | node-08 | 5007 |
| PC 09 | node-09 | 5008 |
| PC 10 | node-10 | 5009 |
| PC 11 | node-11 | 5010 |
| PC 12 | node-12 | 5011 |
| PC 13 | node-13 | 5012 |
| PC 14 | node-14 | 5013 |
| PC 15 | node-15 | 5014 |
| PC 16 | node-16 | 5015 |

---

# 3. Requirements

Every computer should have:

* Windows
* Go
* VS Code (optional)
* The project folder

Check Go:

```powershell
go version
```

Expected:

```text
go version go1.x.x windows/amd64
```

The exact Go version may differ.

---

# 4. Project Structure

The project should look like:

```text
D:\lab-demo
│
├── go.mod
│
├── cmd
│   └── node
│       └── main.go
│
└── internal
    └── network
        ├── node.go
        ├── peer.go
        ├── peer_manager.go
        ├── message.go
        ├── server.go
        └── heartbeat.go
```

The following files are COMMON and should be identical on all computers:

```text
internal/network/node.go
internal/network/peer.go
internal/network/peer_manager.go
internal/network/message.go
internal/network/server.go
internal/network/heartbeat.go
go.mod
```

Only this file is node-specific:

```text
cmd/node/main.go
```

---

# 5. Find the IPv4 Address

On each computer open PowerShell:

```powershell
ipconfig
```

Find the active Ethernet adapter.

Look for:

```text
IPv4 Address
```

Example:

```text
Ethernet adapter Ethernet:

   IPv4 Address. . . . . . : 10.107.3.191
   Subnet Mask . . . . . . : 255.255.224.0
```

Use:

```text
10.107.3.191
```

as the LAN address.

Do NOT use:

```text
127.0.0.1
192.168.56.x
169.254.x.x
IPv6 address
```

The exact IP addresses may change depending on the LAN/DHCP configuration, so check `ipconfig` before deployment.

---

# 6. Maintain the Node Information Table

Before starting the experiment, collect the IPv4 address of every computer.

Create a table:

```text
PC 01 → IP → Port 5000 → node-01
PC 02 → IP → Port 5001 → node-02
PC 03 → IP → Port 5002 → node-03
...
PC 16 → IP → Port 5015 → node-16
```

Example:

```text
node-01 → 10.107.3.191:5000
node-02 → 10.107.2.151:5001
node-03 → 10.107.5.138:5002
```

Use the actual current IP addresses collected from `ipconfig`.

---

# 7. Copy the Project to a Computer

Copy the complete project folder:

```text
lab-demo
```

to:

```text
D:\lab-demo
```

on the target computer.

For example:

```text
D:\lab-demo
```

---

# 8. Verify the Project

Open PowerShell:

```powershell
cd D:\lab-demo
```

Check the Go module:

```powershell
go mod tidy
```

Format the code:

```powershell
go fmt ./...
```

Build:

```powershell
go build ./...
```

Run tests:

```powershell
go test ./...
```

All commands should complete without errors.

---

# 9. Configure Node ID and Port

Open:

```text
D:\lab-demo\cmd\node\main.go
```

Each computer must have its own Node ID and listening port.

Example Node 01:

```go
node := &network.Node{
    ID:      "node-01",
    Address: ":5000",
}
```

Node 02:

```go
node := &network.Node{
    ID:      "node-02",
    Address: ":5001",
}
```

Node 03:

```go
node := &network.Node{
    ID:      "node-03",
    Address: ":5002",
}
```

Continue the same pattern up to Node 16.

---

# 10. Configure Peer Connections

The current implementation uses explicit peer addresses in `main.go`.

Example:

```go
err = node.Connect("10.107.2.151:5001")

if err != nil {
    fmt.Println("Connection to Node B failed:", err)
}
```

Use the actual IP address and port of the target peer.

Example:

```text
node-01 → 10.107.2.151:5001
```

means Node 01 connects to Node 02.

---

# 11. Example Node 01

Example configuration:

```go
package main

import (
    "fmt"
    "time"

    "lab-demo/internal/network"
)

func main() {

    node := &network.Node{
        ID:      "node-01",
        Address: ":5000",
    }

    err := node.Start()

    if err != nil {
        fmt.Println("Failed to start node:", err)
        return
    }

    time.Sleep(1 * time.Second)

    err = node.Connect("NODE_02_IP:5001")

    if err != nil {
        fmt.Println("Connection to Node 02 failed:", err)
    }

    err = node.Connect("NODE_03_IP:5002")

    if err != nil {
        fmt.Println("Connection to Node 03 failed:", err)
    }

    select {}
}
```

Replace:

```text
NODE_02_IP
NODE_03_IP
```

with the actual IP addresses.

---

# 12. Example Node 02

```go
package main

import (
    "fmt"

    "lab-demo/internal/network"
)

func main() {

    node := &network.Node{
        ID:      "node-02",
        Address: ":5001",
    }

    err := node.Start()

    if err != nil {
        fmt.Println("Failed to start node:", err)
        return
    }

    err = node.Connect("NODE_01_IP:5000")

    if err != nil {
        fmt.Println("Connection to Node 01 failed:", err)
    }

    err = node.Connect("NODE_03_IP:5002")

    if err != nil {
        fmt.Println("Connection to Node 03 failed:", err)
    }

    select {}
}
```

---

# 13. Start Order

For the current implementation, start nodes in a controlled order.

Recommended:

```text
1. Node 01
2. Node 02
3. Node 03
4. Node 04
...
16. Node 16
```

If a node attempts to connect before its target node is running, the connection may fail.

Start the target node first when necessary.

---

# 14. Run a Node

On the target computer:

```powershell
cd D:\lab-demo
```

Run:

```powershell
go run ./cmd/node
```

Expected output:

```text
TCP server listening on :5000
Node node-01 started
```

If connections are configured, you should also see:

```text
Node node-01 connected to ...
```

---

# 15. Verify Incoming Connections

When another node connects to this node, output should show something similar to:

```text
Incoming connection from: ...
```

This confirms that the TCP server accepted the peer connection.

---

# 16. Verify PING/PONG

After nodes are connected, heartbeat messages should appear.

Example:

```text
PING → ...
```

Remote node:

```text
PING received from: node-01
```

Then:

```text
PONG received from: node-01
```

This confirms that the heartbeat mechanism is operating.

---

# 17. Verify Broadcast

A broadcast message should look like:

```text
Message received:
  Type: BROADCAST
  Sender: node-01
  Payload: Hello from Node 01
```

The receiving node should display the sender and payload.

---

# 18. Verify Peer Disconnect

Stop a node using:

```text
Ctrl + C
```

A connected peer should eventually detect the disconnection.

Expected output may include:

```text
Peer disconnected: ...
```

and:

```text
Removed disconnected peer: ...
```

If heartbeat detects the failure:

```text
Heartbeat failed: ...
```

or:

```text
Heartbeat timeout: ...
```

---

# 19. TCP Connectivity Test

If a node is running on a particular port, test it from another computer:

```powershell
Test-NetConnection <IP> -Port <PORT>
```

Example:

```powershell
Test-NetConnection 10.107.2.151 -Port 5001
```

Successful result:

```text
TcpTestSucceeded : True
```

This confirms that the TCP port is reachable.

---

# 20. Troubleshooting

## Problem: Node cannot start

Check whether the port is already being used:

```powershell
netstat -ano | findstr :5000
```

Change the port if necessary.

---

## Problem: Connection refused

Check:

1. Target PC is running.
2. Target node is running.
3. IP address is correct.
4. Port is correct.

Test:

```powershell
Test-NetConnection TARGET_IP -Port TARGET_PORT
```

---

## Problem: Wrong IP address

Run:

```powershell
ipconfig
```

Again and verify the active Ethernet adapter.

---

## Problem: Connection works one way but not another

Check that both nodes are running as both:

```text
TCP server
+
TCP client
```

Also verify the target IP and port.

---

# 21. Common Mistakes

Do not use:

```text
127.0.0.1
```

for communication between physical computers.

Do not use:

```text
192.168.56.x
```

if that address belongs to a virtual adapter.

Do not assign the same port to multiple nodes on the same machine.

Do not accidentally copy the wrong `main.go` to a node.

Do not modify the common `internal/network` files independently on different PCs.

---

# 22. Deployment Rule

The common networking code must remain identical.

```text
COMMON
├── node.go
├── peer.go
├── peer_manager.go
├── message.go
├── server.go
└── heartbeat.go
```

Node-specific:

```text
cmd/node/main.go
```

Therefore:

```text
Main PC
   │
   │ Final common code
   ▼
USB / Pendrive
   │
   ├── PC 01
   ├── PC 02
   ├── PC 03
   ├── ...
   └── PC 16
```

---

# 23. Final 16-Node Checklist

Before starting the experiment:

```text
[ ] All PCs connected to LAN
[ ] IPv4 collected for every PC
[ ] Node IDs assigned
[ ] Ports assigned
[ ] Common code copied
[ ] Correct main.go copied to each PC
[ ] go build ./... passes
[ ] go test ./... passes
```

During execution:

```text
[ ] Node 01 running
[ ] Node 02 running
[ ] Node 03 running
[ ] Node 04 running
[ ] Node 05 running
[ ] Node 06 running
[ ] Node 07 running
[ ] Node 08 running
[ ] Node 09 running
[ ] Node 10 running
[ ] Node 11 running
[ ] Node 12 running
[ ] Node 13 running
[ ] Node 14 running
[ ] Node 15 running
[ ] Node 16 running
```

Network verification:

```text
[ ] TCP connections established
[ ] Incoming connections visible
[ ] PeerManager contains peers
[ ] PING works
[ ] PONG works
[ ] Broadcast works
[ ] Disconnect detection works
[ ] Disconnected peer is removed
```

---

# 24. Useful Commands

### Enter project

```powershell
cd D:\lab-demo
```

### Format

```powershell
go fmt ./...
```

### Dependencies

```powershell
go mod tidy
```

### Build

```powershell
go build ./...
```

### Tests

```powershell
go test ./...
```

### Run node

```powershell
go run ./cmd/node
```

### Check IP

```powershell
ipconfig
```

### Check listening port

```powershell
netstat -ano | findstr :5000
```

### Test remote TCP port

```powershell
Test-NetConnection <IP> -Port <PORT>
```

### Stop node

```text
Ctrl + C
```

---

# 25. Current P2P Architecture

```text
                    LAN
                     │
        ┌────────────┼────────────┐
        │            │            │
      Node 01      Node 02      Node 03
        │            │            │
        └────────────┼────────────┘
                     │
                  TCP P2P
                     │
        ┌────────────┼────────────┐
        │            │            │
      Node 04      Node 05     ... Node 16
```

Each node contains:

```text
Node
 │
 ├── TCP Server
 │
 ├── TCP Client
 │
 ├── PeerManager
 │
 ├── Peer
 │
 ├── Message Protocol
 │
 ├── Broadcast
 │
 └── Heartbeat
```

This document describes the current operational procedure for deploying and testing the P2P networking system across the physical LAN nodes.


go fmt ./...
go mod tidy
go build ./...
go test ./...