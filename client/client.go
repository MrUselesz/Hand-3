package main

import (
	"bufio"
	proto "chitchat/grpc"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = ""
var lamport uint32
var left bool

func main() {
	for {

		left = false
		lamport = 0
		file, err := openLogFile("./mylog.log")
		if err != nil {
			log.Fatalf("Not working")
		}
		log.SetOutput(file)

		conn, err := grpc.NewClient("localhost:6060", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Not working")
		}
		fmt.Println("Please write your name.")
		client := proto.NewChitchatClient(conn)
		stream, err := client.SendReceive(context.Background())
		if err != nil {
			log.Fatalf("Client failed to connect on Join: %v", err)
		}
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Did not work")
		}
		lamport += 1
		log.Println(strings.TrimRight(line, "\n"), ": Sending message to join server", strconv.FormatUint(uint64(lamport), 10))
		stream.Send(&proto.ClientMessage{
			Name:    line,
			Lamport: lamport,
		})
		name = line
		go receiver(stream)
		go Sender(client, stream)

		for {
			if left {
				break
			}
		}

	}
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Failed")
	}
	return logFile, nil
}

func receiver(stream proto.Chitchat_SendReceiveClient) {
	for {
		res, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		if res.GetLamport() > lamport {
			lamport = res.GetLamport()
		}
		lamport += 1
		log.Println(strings.TrimRight(name, "\n"), "Recieved", strings.TrimRight(res.GetName(), "\n"), ": ", strings.TrimRight(res.GetMessage(), "\n"), " ", strconv.FormatUint(uint64(lamport), 10))

	}
}

func Sender(client proto.ChitchatClient, stream proto.Chitchat_SendReceiveClient) {
	for {

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Did not work")
		}
		if len([]rune(line)) <= 128 {
			lamport += 1
			stream.Send(&proto.ClientMessage{
				Name:    name,
				Message: line,
				Lamport: lamport,
			})
			log.Println(strings.TrimRight(name, "\n"), ": Sends message", strings.TrimRight(line, "\n"), " ", strconv.FormatUint(uint64(lamport), 10))
			if line == "/leave\n" {
				break
			}
		} else {
			fmt.Println(strings.TrimRight(name, "\n"), " Message to long.")
		}

	}
	left = true
}
