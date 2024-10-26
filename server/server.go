package main

import (
	proto "chitchat/grpc"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type chitchatServer struct {
	proto.UnimplementedChitchatServer
	mu      sync.Mutex
	clients map[string]proto.Chitchat_SendReceiveServer
}

var message string
var name string
var lamport uint32
var newLine string

func (s *chitchatServer) SendReceive(stream proto.Chitchat_SendReceiveServer) error {
	// Register client
	firstMsg, err := stream.Recv()
	if err != nil {
		log.Printf("Error receiving initial message: %v", err)
		return err
	}
	name := firstMsg.GetName()

	s.mu.Lock()
	s.clients[name] = stream
	s.mu.Unlock()
	if firstMsg.GetLamport() > lamport {
		lamport = firstMsg.GetLamport()
	}
	lamport++
	log.Println("%s recieved that", strings.TrimRight(name, newLine), " has joined us ", lamport)
	lamport += 1
	for clientName, clientStream := range s.clients {
		log.Println("%s sent message to ", strings.TrimRight(clientName, newLine), "that", strings.TrimRight(name, newLine), "has joined", lamport)
		if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: " Has joined the server.", Lamport: lamport}); err != nil {
			log.Printf("Error sending to %s: %v", strings.TrimRight(clientName, newLine), err)
		}

	}

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
			if msg.GetLamport() > lamport {
				lamport = msg.GetLamport()
			}
			lamport += 1

			if msg.GetMessage() == "/leave" {
				s.mu.Lock()
				log.Println("%s recieved message about leaving from", strings.TrimRight(name, newLine), lamport)
				delete(s.clients, name)
				lamport += 1
				for clientName, clientStream := range s.clients {
					log.Println("%s we send message to ", strings.TrimRight(clientName, newLine), " that ", strings.TrimRight(name, newLine), "has left", lamport)
					if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: " Has left the server.", Lamport: lamport}); err != nil {
						log.Printf("Error sending to %s: %v", strings.TrimRight(clientName, newLine), err)
					}

				}

				s.mu.Unlock()
			} else {
				s.mu.Lock()
				log.Println("%s recieved message from", strings.TrimRight(name, newLine), lamport)
				lamport += 1
				for clientName, clientStream := range s.clients {
					log.Println("%s we send message to ", strings.TrimRight(clientName, newLine), lamport)
					if err := clientStream.Send(&proto.ServerMessage{Name: name, Message: msg.GetMessage(), Lamport: lamport}); err != nil {
						log.Printf("Error sending to %s: %v", strings.TrimRight(clientName, newLine), err)
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
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	} else {
		newLine = "\n"
	}

	file, err := openLogFile("../mylog.log")
	if err != nil {
		log.Fatalf("Not working")
	}
	log.SetOutput(file)

	server := &chitchatServer{clients: make(map[string]proto.Chitchat_SendReceiveServer)}
	for name := range server.clients {
		log.Println(" -", name)
	}
	lamport = 0
	message = ""
	name = ""
	server.start_server()
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Failed")
	}
	return logFile, nil
}

func (s *chitchatServer) start_server() {

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":6060")
	if err != nil {
		log.Fatalf("Failed to listen on port 5050: %v", err)

	}
	fmt.Println("Server is active")

	proto.RegisterChitchatServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
