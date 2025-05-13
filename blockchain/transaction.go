package blockchain

type Transaction struct {
	address   string
	to        []string
	amount    float64
	signature string
	data      string
}
