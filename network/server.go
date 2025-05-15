package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"yabc/protocol"
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

	request, err := protocol.Receive(conn)

	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading from connection: %v", err)
		}
		return
	}

	switch request.Type {
	case protocol.MsgPeerDiscovery:
		s.network.AddNewDiscoveredPeer(request.Sender, Peer{Connection: conn, Status: true})
		s.network.peerIsOnline(request.Sender)
		response := protocol.NewMessage(protocol.MsgPeerDiscovery, s.network.GetAllKnownPeersAddresses(), s.network.GetNodeAddress())
		err := protocol.Send(response, conn)

		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	default:
		log.Printf("Unknown message type: %s", request.Type)
	}
}
