package main

import (
	proto "ITUServer/grpc"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type ITU_databaseServer struct {
	proto.UnimplementedITUDatabaseServer
	mu      sync.Mutex
	clients map[string]proto.ITUDatabase_SendReceiveServer
}

var message string
var name string

func (s *ITU_databaseServer) SendReceive(stream proto.ITUDatabase_SendReceiveServer) error {
	// Register client
	firstMsg, err := stream.Recv()
	if err != nil {
		log.Printf("Error receiving initial message: %v", err)
		return err
	}
	name := firstMsg.GetName()
	log.Printf("%s joined", name)

	s.mu.Lock()
	s.clients[name] = stream
	s.mu.Unlock()

	for clientName, clientStream := range s.clients {
		if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: " Has joined the server."}); err != nil {
			log.Printf("Error sending to %s: %v", clientName, err)
		}
	}

	defer func() {
		// Handle client disconnection
		log.Printf("%s left", name)
		s.mu.Lock()
		delete(s.clients, name)
		for clientName, clientStream := range s.clients {
			if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: " Has left the server."}); err != nil {
				log.Printf("Error sending to %s: %v", clientName, err)
			}
		}
		s.mu.Unlock()
	}()

	// Goroutine for receiving messages from the client and broadcasting
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Error receiving from %s: %v", name, err)
				return
			}
			log.Printf("Received message from %s: %s", name, msg.Message)

			// Broadcast received message to all clients

			if msg.GetMessage() == "/leave\n" {
				log.Printf("%s left", name)
				s.mu.Lock()
				delete(s.clients, name)
				for clientName, clientStream := range s.clients {
					if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: " Has left the server."}); err != nil {
						log.Printf("Error sending to %s: %v", clientName, err)
					}
				}

				s.mu.Unlock()
			} else {
				s.mu.Lock()
				for clientName, clientStream := range s.clients {
					if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: msg.GetMessage()}); err != nil {
						log.Printf("Error sending to %s: %v", clientName, err)
					}
				}
				s.mu.Unlock()
			}

		}
	}()

	<-stream.Context().Done()
	return nil
}

func main() {

	server := &ITU_databaseServer{clients: make(map[string]proto.ITUDatabase_SendReceiveServer)}
	for name := range server.clients {
		log.Println(" -", name)
	}
	message = ""
	name = ""
	server.start_server()
}

func (s *ITU_databaseServer) start_server() {

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":6060")
	if err != nil {
		log.Fatalf("Failed to listen on port 5050: %v", err)

	}
	fmt.Println("Server is active")

	proto.RegisterITUDatabaseServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
