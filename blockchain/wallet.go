package blockchain

type Wallet struct {
	address    string // hashed public key
	privateKey string
	publicKey  string
}
