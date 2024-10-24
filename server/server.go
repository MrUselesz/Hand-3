package main

import (
	proto "ITUServer/grpc"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ITU_databaseServer struct {
	proto.UnimplementedITUDatabaseServer
}

var message string
var name string

func (joining *ITU_databaseServer) Join(stream proto.ITUDatabase_JoinServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		name = req.GetName()
		fmt.Println("Recieved ", name)

	}

	return nil
}

func (recieve *ITU_databaseServer) ClientSend(stream proto.ITUDatabase_ClientSendServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		name = req.GetName()
		message = req.GetMessage()
		fmt.Println("Recieved ", name, " ", message)

	}

	return nil
}

func (sender *ITU_databaseServer) ServerSend(req *proto.Empty, stream proto.ITUDatabase_ServerSendServer) error {
	for {

		if message != "" {
			if err := stream.Send(&proto.ServerMessage{Message: message}); err != nil {
				log.Fatalf("Failed to send: %v", err)
			}
			fmt.Println("Sending ", message)
			message = ""
			name = ""
			break
		} else if name != "" && message == "" {
			if err := stream.Send(&proto.ServerMessage{Name: name + "Has joined the server"}); err != nil {
				log.Fatalf("Failed to send: %v", err)
			}
			fmt.Println("Sending ", name)
			name = ""
			message = ""
			break
		}

	}

	return nil
}

func main() {

	server := &ITU_databaseServer{}
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
