package network

import (
	"fmt"
	"net"
	"sync"
)

const DefaultNodeAddress = "127.0.0.1:7071"
const (
	Done               = "DONE"
	Identify           = "IDENTIFY"
	RequestPeersList   = "REQUEST_PEERS_LIST"
	RequestNodeAddress = "REQUEST_NODE_ADDRESS"
	Eol                = "\n"
	CommandDelimiter   = ">>>"
)

type Network struct {
	server          *Server
	client          *Client
	mutex           sync.Mutex
	discoveredPeers map[string]Peer
	nodeAddress     string
	config          *Config
}

type Peer struct {
	Connection net.Conn
	Status     bool
}

func NewNetwork(nodeAddress string, options map[string]string) *Network {
	finalNodeAddress := DefaultNodeAddress

	if nodeAddress != "" {
		finalNodeAddress = nodeAddress
	}

	network := &Network{
		nodeAddress:     finalNodeAddress,
		discoveredPeers: make(map[string]Peer),
	}

	network.config = NewNetworkConfig(options)

	network.server = NewServer(network)
	network.client = NewClient(network)

	network.initDefaultPeers()

	return network
}

func (n *Network) initDefaultPeers() {
	for _, peer := range n.getDefaultPeers() {
		n.AddNewDiscoveredPeer(peer, Peer{Connection: nil, Status: false})
	}
}

func (n *Network) getDefaultPeers() []string {
	return []string{"127.0.0.1:7071", "127.0.0.1:1234"}
}

func (n *Network) Start() {
	if n.config.IsDebugEnabled() {
		fmt.Printf("\n\n*** DEBUG ENABLED ***\n\n")
	}

	go n.client.IntroduceSelf()
	go n.client.RequestPeersList()

	n.server.Serve()
}

func (n *Network) AddNewDiscoveredPeer(newPeerAddress string, newPeer Peer) {
	func() {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		if _, exists := n.discoveredPeers[newPeerAddress]; !exists {
			n.discoveredPeers[newPeerAddress] = newPeer
		}
	}()
}

func (n *Network) GetAllKnownPeersAddresses() map[string]Peer {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.discoveredPeers
}

func (n *Network) peerDisconnected(peerAddress string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if peer, exists := n.discoveredPeers[peerAddress]; exists {
		peer.Status = false
		n.discoveredPeers[peerAddress] = peer
	}
}

func (n *Network) PrintPeersList() {
	fmt.Printf("\n### START Peers List ###\n")
	for peerAddress, peer := range n.GetAllKnownPeersAddresses() {
		if peerAddress == n.nodeAddress {
			continue
		}
		fmt.Printf("- [%s]:\t %t\n", peerAddress, peer.Status)
	}
	fmt.Printf("###  END  Peers  List ###\n\n")
}

func (n *Network) peerIsOnline(peerAddress string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if peer, exists := n.discoveredPeers[peerAddress]; exists {
		peer.Status = true
		n.discoveredPeers[peerAddress] = peer
	}
}
