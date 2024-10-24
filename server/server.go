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
			return err
		}

		fmt.Println("Recieved ", req.GetName())
		name = req.GetName()
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
			return err
		}

		fmt.Println("Recieved ", req.GetMessage())
		message = req.GetMessage()
	}

	return nil
}

func (sender *ITU_databaseServer) ServerSend(req *proto.Empty, stream proto.ITUDatabase_ServerSendServer) error {
	for {
		if name != "Vagina" {
			if err := stream.Send(&proto.ServerMessage{Message: name}); err != nil {
				return err
			}
			fmt.Println("Sending ", name)
			name = "Vagina"
		}
		if message != "Penis" {
			if err := stream.Send(&proto.ServerMessage{Message: message}); err != nil {
				return err
			}
			fmt.Println("Sending ", message)
			message = "Penis"

		}

	}

	return nil
}

func main() {

	server := &ITU_databaseServer{}
	message = "Penis"
	name = "Vagina"
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
