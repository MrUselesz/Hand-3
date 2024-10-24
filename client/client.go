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
var senderName = ""
var message = ""

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
		time.Sleep(5 * time.Millisecond)
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("error: %v", err)

		}
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		senderName = req.GetName()
		message = req.GetMessage()

		fmt.Println(senderName, " ", message)

	}

}

func Sender(client proto.ITUDatabaseClient) {
	for {
		time.Sleep(10 * time.Millisecond)
		stream, err := client.ClientSend(context.Background())
		if err != nil {
			log.Fatalf("Client failed to connect on sending a message: %v", err)
		}
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Did not work")
		}
		if len([]rune(line)) <= 128 {
			stream.Send(&proto.ClientMessage{
				Message: line,
			})

			stream.CloseSend()
		} else {
			fmt.Println("Message to long.")
		}

	}
}

func Join(client proto.ITUDatabaseClient) {

	stream, err := client.Join(context.Background())
	if err != nil {
		log.Fatalf("Client failed to connect on Join: %v", err)
	}
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
