package blockchain

type Transaction struct {
	Address   string   `json:"address"`
	To        []string `json:"to"`
	Amount    float64  `json:"amount"`
	Signature string   `json:"signature"`
	Data      string   `json:"data"`
}
