package main

import (
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	socketPath := filepath.Join(os.TempDir(), "yabc_wallet_127.0.0.1:7071.sock")

	conn, err := net.Dial("unix", socketPath)
	defer conn.Close()

	if err != nil {
		panic(err)
	}

	conn.Write([]byte("transaction sent: " + time.Now().String()))
}
