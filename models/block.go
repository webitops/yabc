package models

type Block struct {
	data         string
	transactions string
	hash         string
	prevHash     string
	nonce        int
	timestamp    int64
	difficulty   int
}
