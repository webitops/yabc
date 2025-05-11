package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	network *Network
}

func NewServer(network *Network) *Server {
	return &Server{
		network: network,
	}
}

func (s *Server) Serve() {
	ln, err := net.Listen("tcp", s.network.nodeAddress)
	fmt.Println("Listening on " + s.network.nodeAddress)

	if err != nil {
		panic(err)
	}

	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			log.Fatal("Error closing listener: ", err)
		}
	}(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection: ", err)
		}

		go s.handlePeerRequest(conn)
	}
}

func (s *Server) handlePeerRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Error closing connection: ", err)
		}
	}(conn)

	if peer, exists := s.network.discoveredPeers[conn.RemoteAddr().String()]; exists {
		peer.Connection = conn
		peer.Status = true
		s.network.AddNewDiscoveredPeer(conn.RemoteAddr().String(), peer)
	} /* else {
		fmt.Println("Requesting node address from peer:" + conn.RemoteAddr().String())
		_, err := conn.Write([]byte(RequestNodeAddress + CommandDelimiter + Eol))
		if err != nil {
			log.Fatal("Error writing to peer: ", err)
		}

		nodeAddress, err := bufio.NewReader(conn).ReadString(Eol[0])

		s.network.AddNewDiscoveredPeer(nodeAddress, Peer{Connection: nil, Status: true})
		return
	}*/

	command, err := bufio.NewReader(conn).ReadString(Eol[0])

	if err != nil {
		fmt.Println("Peer disconnected." + conn.RemoteAddr().String())
		s.network.peerDisconnected(conn.RemoteAddr().String())
		return
	}

	s.HandleRequestCommand(conn, command)
}

func (s *Server) HandleRequestCommand(conn net.Conn, rawCommand string) {
	if strings.Contains(rawCommand, Identify+CommandDelimiter) {
		newPeer := strings.Split(rawCommand, CommandDelimiter)[1]
		s.network.AddNewDiscoveredPeer(strings.Trim(newPeer, Eol), Peer{Connection: conn, Status: true})
		_, err := conn.Write([]byte(Done + CommandDelimiter + Identify + CommandDelimiter + newPeer + Eol))
		if err != nil {
			log.Fatal("Error writing to peer: ", err)
		}
	} else if strings.Contains(rawCommand, RequestPeersList+CommandDelimiter) {
		peersJson, err := json.Marshal(s.network.GetAllKnownPeersAddresses())
		fmt.Println("RP sending this RAW: ", string(peersJson))
		if err != nil {
			log.Fatal("Error marshalling peers list: ", err)
		}
		_, err = conn.Write([]byte(Done + CommandDelimiter + RequestPeersList + CommandDelimiter + string(peersJson) + Eol))
		if err != nil {
			log.Fatal("Error writing to peer: ", err)
		}
	} else {
		_, err := conn.Write([]byte(Done + Eol))
		if err != nil {
			log.Fatal("Error writing to peer: ", err)
		}
	}
}
