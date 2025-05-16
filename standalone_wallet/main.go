package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"yabc/models"
	"yabc/protocol"
)

func main() {
	socketPort := flag.String("port", "7071", "The port to listen on for tcp socket requests.")
	address := flag.String("addr", "", "The address to send the transaction to.")
	amount := flag.Float64("amount", 0, "The amount to send.")
	data := flag.String("data", "", "Other data to send.")

	flag.Parse()

	transaction := &models.Transaction{
		Address: "current_wallet_address",
		To:      []string{*address},
		Amount:  *amount,
		Data:    *data,
	}

	msg := protocol.NewMessage(protocol.MsgSubmitTx, transaction, "current_wallet_address")

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

	protocol.Send(msg, conn)

	if err != nil {
		panic(err)
	}
}
