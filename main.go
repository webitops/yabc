package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"yabc/blockchain"
)

func main() {
	fmt.Println("Node started.")

	nodeAddressPtr := flag.String("node-address", "", "The address to listen on for HTTP requests.")
	debugMode := flag.Bool("debug", false, "Enable debug mode.")
	walletMode := flag.Bool("wallet", false, "Enable wallet mode.")
	flag.Parse()

	options := make(map[string]string)
	options["debug"] = fmt.Sprintf("%t", *debugMode)
	options["node-address"] = *nodeAddressPtr
	bc := blockchain.NewBlockchain(options)

	go bc.Start()

	if *walletMode {
		fmt.Println("Wallet started.")
		for {
			fmt.Println("Transaction to send: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := scanner.Text()
			bc.BroadcastTransaction(answer)
		}

	}

	// Wait.
	wait := make(chan bool)
	<-wait
}
