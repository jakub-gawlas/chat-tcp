package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	address     string
	listener    net.Listener
	messages    chan Message
	clients     map[string]*Client
	clientsLock sync.RWMutex
}

type Message struct {
	senderId string
	data     []byte
}

func NewServer(host, port string) Server {
	return Server{
		address:     fmt.Sprintf("%s:%s", host, port),
		messages:    make(chan Message),
		clients:     map[string]*Client{},
		clientsLock: sync.RWMutex{},
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	s.listener = listener

	fmt.Printf("Server started on address %s.\n", s.address)

	go s.broadcastMessages()

	s.handleConnections()

	return nil
}

func (s *Server) broadcastMessages() {
	for {
		select {
		case msg := <-s.messages:
			s.clientsLock.RLock()
			for _, client := range s.clients {
				if client.id == msg.senderId {
					continue
				}
				if err := client.Send(msg.data); err != nil {
					client.Log("error while send message: " + err.Error())
					continue
				}
			}
			s.clientsLock.RUnlock()
		default:
		}
	}
}

func (s *Server) handleConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Printf("error while accept connection: %s", err.Error())
			continue
		}

		client := NewClient(conn)
		s.addClient(client)
		go func() {
			client.HandleConnection(s.messages)
			s.removeClient(client)
		}()

		client.Log("client connected")
	}
}

func (s *Server) addClient(c Client) {
	s.clientsLock.Lock()
	s.clients[c.id] = &c
	s.clientsLock.Unlock()
}

func (s *Server) removeClient(c Client) {
	s.clientsLock.Lock()
	delete(s.clients, c.id)
	s.clientsLock.Unlock()
}

// Close gracefully terminate server
func (s *Server) Close() {
	message := []byte("Server is terminating.")
	s.clientsLock.RLock()
	for _, client := range s.clients {
		if err := client.Send(message); err != nil {
			client.Log("error while send message: " + err.Error())
		}
		client.Close()
	}
	s.clientsLock.RUnlock()

	defer s.listener.Close()
}
