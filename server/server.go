package main

import (
	proto "ITUServer/grpc"
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type ITU_databaseServer struct {
	proto.UnimplementedITUDatabaseServer
	smessages []string
}

func (s *ITU_databaseServer) SendRecieve(ctx context.Context, in *proto.Empty) (*proto.ServerMessage, error) {
	return &proto.ServerMessage{Smessages: s.smessages}, nil
}

func main() {

	server := &ITU_databaseServer{smessages: []string{}}
	server.start_server()

	time.Sleep(30*time.Second)

	cmessage, err := server.SendRecieve(context.Background(), &proto.ClientMessage{})
	if err != nil {
		log.Fatalf("Not working")
	}

	//server.smessages = append(server.smessages, cmessage)
	log.Fatalln(cmessage)

}

func (s *ITU_databaseServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterITUDatabaseServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
