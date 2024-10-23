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

func (sender *ITU_databaseServer) Sender(stream proto.ITUDatabase_SenderServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("Recieved ", req.GetMessage())

		if err := stream.Send(&proto.ServerMessage{Message: req.GetMessage()}); err != nil {
			return nil
		}

	}

	return nil
}

/*
recieves a message from the Client.

	func (recieve *ITU_databaseServer) ClientSend(stream proto.ITUDatabase_ClientSendServer) error {
		value, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Empty{})
		}
		if err != nil {
			log.Fatalf("Not working")

		}
		fmt.Println(value.GetMessage())
		return err
	}

	func (send *ITU_databaseServer) ServerSend(req *proto.Empty, stream proto.ITUDatabase_ServerSendServer) error {
		for {
			select {
			case <-stream.Context().Done():
				return status.Error(codes.Canceled, "Stream has ended")
			default:
				time.Sleep(1 * time.Minute)

				value := "Complete"

				err := stream.SendMsg(&proto.ServerMessage{
					Message: value,
				})
				if err != nil {
					log.Fatalf("Not working")

				}
			}

		}
	}
*/
func main() {

	server := &ITU_databaseServer{}
	server.start_server()
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
