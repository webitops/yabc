package main

import (
	"flag"
	"fmt"
	"yabc/network"
)

func main() {
	fmt.Println("Node started.")

	nodeAddressPtr := flag.String("node-address", "", "The address to listen on for HTTP requests.")
	debugMode := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	options := make(map[string]string)
	options["debug"] = fmt.Sprintf("%t", *debugMode)

	n := network.NewNetwork(*nodeAddressPtr, options)

	n.Start()

}
