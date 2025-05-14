package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"yabc/blockchain"
)

func main() {
	socketPort := flag.String("port", "7071", "The port to listen on for tcp socket requests.")
	address := flag.String("addr", "", "The address to send the transaction to.")
	amount := flag.Float64("amount", 0, "The amount to send.")
	data := flag.String("data", "", "Other data to send.")

	flag.Parse()

	transaction := &blockchain.Transaction{
		Address: "current_wallet_address",
		To:      []string{*address},
		Amount:  *amount,
		Data:    *data,
	}

	socketPath := filepath.Join(os.TempDir(), "yabc_wallet_127.0.0.1:"+*socketPort+".sock")

	conn, err := net.Dial("unix", socketPath)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	if err != nil {
		panic(err)
	}

	serializedTransaction, err := json.Marshal(transaction)

	if err != nil {
		panic(err)
	}
	fmt.Println(string(serializedTransaction))

	_, err = conn.Write([]byte(serializedTransaction))
	if err != nil {
		panic(err)
	}
}
