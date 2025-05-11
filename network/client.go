package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	network *Network
}

func NewClient(network *Network) *Client {
	return &Client{
		network: network,
	}
}

func (c *Client) IntroduceSelf() {
	c.SendToPeers(Identify, c.network.nodeAddress)
}

func (c *Client) RequestPeersList() {
	for {
		c.SendToPeers(RequestPeersList, "")
		time.Sleep(5 * time.Second)
		c.network.PrintPeersList()
	}
}

func (c *Client) SendToPeers(command string, params string) {

	for peerAddress := range c.network.GetAllKnownPeersAddresses() {
		if peerAddress == c.network.nodeAddress {
			continue
		}
		conn, err := net.Dial("tcp", peerAddress)

		if err != nil {
			c.network.peerDisconnected(peerAddress)
			continue
		}

		c.network.peerIsOnline(peerAddress)

		_, err = conn.Write([]byte(command + CommandDelimiter + params + Eol))
		if err != nil {
			log.Print("Error writing to peer: ", err)
		}

		response, _ := bufio.NewReader(conn).ReadString(Eol[0])

		c.handlePeerResponseForRequest(command, response, conn)

		err = conn.Close()
		if err != nil {
			log.Print("Error closing connection: ", err)
		}

	}

}

func (c *Client) handlePeerResponseForRequest(command string, response string, conn net.Conn) {
	switch command {
	case RequestNodeAddress:
		fmt.Println("sending my node address to peer: " + conn.RemoteAddr().String())
		_, err := conn.Write([]byte(Identify + CommandDelimiter + c.network.nodeAddress + Eol))
		if err != nil {
			log.Print("Error writing to peer: ", err)
		}
		break
	case Identify:
		break
	case RequestPeersList:
		responseElements := strings.Split(response, CommandDelimiter)
		newPeersJson := responseElements[len(responseElements)-1]
		newPeers := make(map[string]struct{})
		err := json.Unmarshal([]byte(newPeersJson), &newPeers)
		if err != nil {
			log.Print("Error parsing peers list: ", err)
			break
		}
		for newPeer := range newPeers {
			c.network.AddNewDiscoveredPeer(newPeer, Peer{Connection: nil, Status: false})
		}
		break
	default:
		fmt.Println("Received other [SRC CMD:" + command + " ]: " + response)
		break
	}
}
