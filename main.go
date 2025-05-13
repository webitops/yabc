package main

import (
	"flag"
	"fmt"
	"yabc/blockchain"
	"yabc/wallet_socket"
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

	ws := wallet_socket.NewWalletSocket(*bc)
	go ws.Listen()

	// Wait.
	wait := make(chan bool)
	<-wait
}
