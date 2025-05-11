package main

import (
	"flag"
	"fmt"
	"yabc/network"
)

func main() {
	fmt.Println("Node started.")

	nodeAddressPtr := flag.String("node-address", "", "The address to listen on for HTTP requests.")
	flag.Parse()

	n := network.NewNetwork(*nodeAddressPtr)

	n.Start()

}
