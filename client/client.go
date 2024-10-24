package main

import (
	proto "ITUServer/grpc"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = ""

func main() {

	conn, err := grpc.NewClient("localhost:6060", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewITUDatabaseClient(conn)
	Join(client)
	fmt.Println("We have moved on")
	go reciever(client)
	go Sender(client)

	time.Sleep(1 * time.Hour)
}

func reciever(client proto.ITUDatabaseClient) {
	for {
		stream, err := client.ServerSend(context.Background(), &proto.Empty{})
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		fmt.Println(res.GetMessage())
	}
}

func Sender(client proto.ITUDatabaseClient) {
	for {
		stream, err := client.ClientSend(context.Background())
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Did not work")
		}

		stream.Send(&proto.ClientMessage{
			Message: line,
		})

		stream.CloseSend()
	}
}

func Join(client proto.ITUDatabaseClient) {

	stream, err := client.Join(context.Background())
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Did not work")
	}

	stream.Send(&proto.Joining{
		Name: line,
	})
	name = line
	stream.CloseSend()

}
