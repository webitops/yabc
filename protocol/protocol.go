package protocol

import (
	"bufio"
	"encoding/json"
	"net"
	"time"
)

type Message struct {
	Type      MessageType `json:"type"`
	Payload   []byte      `json:"payload"`
	Sender    string      `json:"sender"`
	Timestamp int64       `json:"timestamp"`
}

type MessageType string

const (
	// Node-to-Node messages
	MsgNewBlock      MessageType = "NEW_BLOCK"
	MsgRequestBlocks MessageType = "REQUEST_BLOCKS"
	MsgBlocks        MessageType = "BLOCKS"
	MsgTransaction   MessageType = "TRANSACTION"
	MsgPeerDiscovery MessageType = "PEER_DISCOVERY"
	MsgSelfIntroduce MessageType = "SELF_INTRODUCE"

	// Wallet-to-Node messages
	MsgSubmitTx     MessageType = "SUBMIT_TX"
	MsgQueryBalance MessageType = "QUERY_BALANCE"
	MsgQueryTx      MessageType = "QUERY_TX"
	MsgQueryBlocks  MessageType = "QUERY_BLOCKS"
)

func NewMessage(msgType MessageType, payload interface{}, sender string) *Message {
	jsonData, _ := json.Marshal(payload)

	return &Message{
		Type:      msgType,
		Payload:   jsonData,
		Sender:    sender,
		Timestamp: time.Now().Unix(),
	}
}

func (m *Message) Encode() ([]byte, error) {
	result, err := json.Marshal(m)
	return append(result, '\n'), err
}

func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

func Send(msg *Message, conn net.Conn) error {
	message, err := msg.Encode()
	if err != nil {
		return err
	}

	_, err = conn.Write(message)

	return err
}

func Receive(conn net.Conn) (*Message, error) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	msg, err := DecodeMessage([]byte(message))

	return msg, err
}
