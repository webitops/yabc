package blockchain

import "yabc/network"

type Blockchain struct {
	network *network.Network
	config  *Config
}

func NewBlockchain(options map[string]string) *Blockchain {
	n := &Blockchain{
		config: NewBlockchainConfig(options),
	}

	n.network = network.NewNetwork(n.config.getNodeAddress(), options)

	return n
}

func (b *Blockchain) Start() {
	b.network.Start()
}
