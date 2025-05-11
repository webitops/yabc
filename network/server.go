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
			log.Print("Error closing listener: ", err)
		}
	}(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Error accepting connection: ", err)
		}

		go s.handlePeerRequest(conn)
	}
}

func (s *Server) IsDebugEnabled() bool {
	return s.network.config.IsDebugEnabled()
}

func (s *Server) handlePeerRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Print("Error closing connection: ", err)
		}
	}(conn)

	command, err := bufio.NewReader(conn).ReadString(Eol[0])

	if err != nil {
		if s.IsDebugEnabled() {
			fmt.Println("Peer disconnected." + conn.RemoteAddr().String())
		}
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
			log.Print("Error writing to peer: ", err)
		}
	} else if strings.Contains(rawCommand, RequestPeersList+CommandDelimiter) {
		peersJson, err := json.Marshal(s.network.GetAllKnownPeersAddresses())
		if err != nil {
			log.Print("Error marshalling peers list: ", err)
		}
		_, err = conn.Write([]byte(Done + CommandDelimiter + RequestPeersList + CommandDelimiter + string(peersJson) + Eol))
		if err != nil {
			log.Print("Error writing to peer: ", err)
		}
	} else {
		_, err := conn.Write([]byte(Done + Eol))
		if err != nil {
			log.Print("Error writing to peer: ", err)
		}
	}
}
