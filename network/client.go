package network

import (
	"encoding/json"
	"log"
	"net"
	"time"
	"yabc/protocol"
)

type Client struct {
	network *Network
}

func NewClient(network *Network) *Client {
	return &Client{
		network: network,
	}
}

func (c *Client) RequestPeersList() {
	msg := protocol.NewMessage(protocol.MsgPeerDiscovery, c.network.nodeAddress, c.network.GetNodeAddress())

	for {
		c.SendToPeers(msg)
		time.Sleep(10 * time.Second)
		if c.IsDebugEnabled() {
			c.network.PrintPeersList()
		}
	}
}

func (c *Client) IsDebugEnabled() bool {
	return c.network.config.IsDebugEnabled()
}

func (c *Client) SendToPeers(msg *protocol.Message) {

	for peerAddress := range c.network.GetAllKnownPeersAddresses() {
		go func(peerAddress string) {
			if peerAddress == c.network.nodeAddress {
				return
			}

			conn, sendErr := net.Dial("tcp", peerAddress)
			if sendErr != nil {
				log.Println(sendErr)
				return
			}
			defer conn.Close()

			sendErr = protocol.Send(msg, conn)

			if sendErr != nil {
				log.Println(sendErr)
			}

			response, receiveErr := protocol.Receive(conn)

			c.handleResponse(response)

			if receiveErr != nil {
				log.Println(receiveErr)
			}
		}(peerAddress)
	}

}

func (c *Client) handleResponse(response *protocol.Message) {
	// First, marshal the payload back to JSON
	switch response.Type {
	case protocol.MsgPeerDiscovery:
		receivedPeers := make(map[string]Peer)
		if err := json.Unmarshal(response.Payload, &receivedPeers); err != nil {
			log.Println("Error unmarshaling to map[string]Peer:", err)
			return
		}

		// Now you can use receivedPeers as a map[string]Peer
		for addr, peer := range receivedPeers {
			c.network.AddNewDiscoveredPeer(addr, peer)
		}
		break
	case protocol.MsgSubmitTx:
		log.Println("Received transaction: ACCEPTED")
		break

	default:
		log.Printf("Unknown message type: %s", response.Type)
	}
}
