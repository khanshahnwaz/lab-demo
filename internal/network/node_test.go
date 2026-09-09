package network

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestHandleConnectionExchangesHello(t *testing.T) {
	node := &Node{ID: "node-A", PeerManager: NewPeerManager()}
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	done := make(chan struct{})
	go func() {
		node.handleConnection(serverConn)
		close(done)
	}()

	clientReader := bufio.NewReader(clientConn)
	if deadlineErr := clientConn.SetReadDeadline(time.Now().Add(time.Second)); deadlineErr != nil {
		t.Fatal(deadlineErr)
	}

	line, err := clientReader.ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}

	var hello Message
	if err := json.Unmarshal(line, &hello); err != nil {
		t.Fatal(err)
	}
	if hello.Type != "HELLO" || hello.SenderID != "node-A" {
		t.Fatalf("unexpected HELLO: %+v", hello)
	}

	remoteHello, err := json.Marshal(Message{Type: "HELLO", SenderID: "node-B"})
	if err != nil {
		t.Fatal(err)
	}
	remoteHello = append(remoteHello, '\n')
	if _, err := clientConn.Write(remoteHello); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		peers := node.PeerManager.GetPeers()
		if len(peers) == 1 && peers[0].ID == "node-B" {
			_ = clientConn.Close()
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatalf("peer identity was not recorded: %+v", node.PeerManager.GetPeers())
}
