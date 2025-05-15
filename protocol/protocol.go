package protocol

import (
	"encoding/json"
	"time"
)

// Message represents the standard communication format
type Message struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	Sender    string      `json:"sender"`
	Timestamp int64       `json:"timestamp"`
	Signature []byte      `json:"signature,omitempty"`
}

type MessageType string

// Define message types
const (
	// Node-to-Node messages
	MsgNewBlock      MessageType = "NEW_BLOCK"
	MsgRequestBlocks MessageType = "REQUEST_BLOCKS"
	MsgBlocks        MessageType = "BLOCKS"
	MsgTransaction   MessageType = "TRANSACTION"
	MsgPeerDiscovery MessageType = "PEER_DISCOVERY"

	// Wallet-to-Node messages
	MsgSubmitTx     MessageType = "SUBMIT_TX"
	MsgQueryBalance MessageType = "QUERY_BALANCE"
	MsgQueryTx      MessageType = "QUERY_TX"
	MsgQueryBlocks  MessageType = "QUERY_BLOCKS"
)

// NewMessage creates a new protocol message
func NewMessage(msgType MessageType, payload interface{}, sender string) *Message {
	return &Message{
		Type:      msgType,
		Payload:   payload,
		Sender:    sender,
		Timestamp: time.Now().Unix(),
	}
}

// Encode serializes a message to JSON
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// DecodeMessage deserializes a message from JSON
func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
