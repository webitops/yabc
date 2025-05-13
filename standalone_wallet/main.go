package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"
)

func main() {
	socketPort := flag.String("port", "7071", "The port to listen on for tcp socket requests.")
	transaction := flag.String("tx", "<empty-transaction>", "The transaction to send.")
	flag.Parse()

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

	_, err = conn.Write([]byte(*transaction))
	if err != nil {
		panic(err)
	}
}
