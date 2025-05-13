package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"yabc/blockchain"
)

func main() {
	fmt.Println("Node started.")

	nodeAddressPtr := flag.String("node-address", "", "The address to listen on for HTTP requests.")
	debugMode := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	options := make(map[string]string)
	options["debug"] = fmt.Sprintf("%t", *debugMode)
	options["node-address"] = *nodeAddressPtr
	bc := blockchain.NewBlockchain(options)

	go bc.Start()

	// handle wallet requests
	go func() {
		fmt.Println("Listening for wallet requests...")

		socketPath := filepath.Join(os.TempDir(), "yabc_wallet_"+bc.GetNodeAddress()+".sock")

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

			go func(conn net.Conn) {
				defer func(conn net.Conn) {
					err := conn.Close()
					if err != nil {

					}
				}(conn)

				message, err := io.ReadAll(conn)

				if err != nil {
					log.Println("Error reading from connection: ", err)
				}

				bc.BroadcastTransaction(string(message))
			}(conn)
		}

	}()

	// Wait.
	wait := make(chan bool)
	<-wait
}
