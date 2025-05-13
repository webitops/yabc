package blockchain

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"yabc/network"
)

type Blockchain struct {
	blocks  []*Block
	network *network.Network
	config  *Config
}

func NewBlockchain(options map[string]string) *Blockchain {
	n := &Blockchain{
		config: NewBlockchainConfig(options),
	}

	n.network = network.NewNetwork(n.config.getNodeAddress(), options)

	return n
}

func (b *Blockchain) Start() {
	go b.listenForTransactions()
	b.network.Start()
}

func (b *Blockchain) listenForTransactions() {
	fmt.Println("Listening for transactions...")

	socketPath := filepath.Join(os.TempDir(), "yabc_wallet.sock")

	err := os.Remove(socketPath)
	if err != nil && !os.IsNotExist(err) {
		log.Print("Error removing socket file: ", err)
	}

	ln, err := net.Listen("unix", socketPath)

	if err != nil {
		log.Print("Error listening for transactions: ", err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Print("Error accepting connection: ", err)
		}
		go b.handleWalletConnection(conn)
	}
}

func (b *Blockchain) handleWalletConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Print("Error closing connection: ", err)
		}
	}(conn)

	transactions, _ := io.ReadAll(conn)

	fmt.Println("Received transaction: ")
	fmt.Println(string(transactions))
}

func (b *Blockchain) BroadcastTransaction(transaction string) {
	b.network.BroadcastMessage(transaction)
}

func (b *Blockchain) GetNodeAddress() string {
	return b.network.GetNodeAddress()
}
