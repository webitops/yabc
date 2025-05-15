package wallet_socket

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"yabc/blockchain"
	"yabc/protocol"
)

type WalletSocket struct {
	blockchain *blockchain.Blockchain
}

func NewWalletSocket(blockchain blockchain.Blockchain) *WalletSocket {
	return &WalletSocket{
		blockchain: &blockchain,
	}
}

func (s *WalletSocket) Listen() {
	fmt.Println("Listening for wallet requests...")

	socketPath := filepath.Join(os.TempDir(), "yabc_wallet_"+s.blockchain.GetNodeAddress()+".sock")

	fmt.Println("Socket path: ", socketPath)
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Error removing existing socket file: %v", err)
	}

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Printf("Error listening on socket: %v", err)
		return
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err)
		}

		go s.handleWalletRequest(conn)
	}
}

func (s *WalletSocket) handleWalletRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	message, err := io.ReadAll(conn)

	msg := &protocol.Message{}

	msg, _ = protocol.DecodeMessage(message)

	if err != nil {
		log.Println("Error reading from connection: ", err)
	}

	s.blockchain.BroadcastTransaction(msg)
}
