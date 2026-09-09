package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Type     string `json:"type"`
	SenderID string `json:"sender_id"`
	Payload  string `json:"payload"`
}

func SendMessage(peer *Peer, message Message) error {
	if peer == nil || peer.Conn == nil {
		return fmt.Errorf("invalid peer connection")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	peer.WriteMux.Lock()
	defer peer.WriteMux.Unlock()

	_, err = peer.Conn.Write(data)
	return err
}

func ReceiveMessages(conn net.Conn, handler func(Message), onDisconnect func()) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		var message Message

		err := json.Unmarshal(scanner.Bytes(), &message)
		if err != nil {
			fmt.Println("Invalid message:", err)
			continue
		}

		handler(message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection read error:", err)
	}

	fmt.Println("Peer disconnected:", conn.RemoteAddr())

	onDisconnect()
}
