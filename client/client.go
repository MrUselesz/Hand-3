package main

import (
	proto "ITUServer/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = ""
var senderName = ""
var message = ""

func main() {

	file, err := openLogFile("./mylog.log")
	if err != nil {
		log.Fatalf("Not working")
	}
	log.SetOutput(file)
	log.Println("log created")

	conn, err := grpc.NewClient("localhost:6060", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}
	fmt.Println("Please write your name.")
	client := proto.NewITUDatabaseClient(conn)
	stream, err := client.SendReceive(context.Background())
	if err != nil {
		log.Fatalf("Client failed to connect on Join: %v", err)
	}
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Did not work")
	}
	stream.Send(&proto.ClientMessage{
		Name: line,
	})

	go receiver(stream)
	go Sender(client, stream)

	time.Sleep(1 * time.Hour)
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Failed")
	}
	return logFile, nil
}

func receiver(stream proto.ITUDatabase_SendReceiveClient) {
	for {
		res, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		log.Println(strings.TrimRight(res.GetName(), "\n"), ": ", res.GetMessage())
	}
}

func Sender(client proto.ITUDatabaseClient, stream proto.ITUDatabase_SendReceiveClient) {
	for {

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Did not work")
		}
		if len([]rune(line)) <= 128 {
			stream.Send(&proto.ClientMessage{
				Name:    name,
				Message: line,
			})

		} else {
			log.Println(name, " Message to long.")
		}

	}
}
