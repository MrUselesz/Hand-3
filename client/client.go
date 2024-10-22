package main

import (
	proto "ITUServer/grpc"
	"bufio"
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ITU_databaseServer struct {
	proto.UnimplementedITUDatabaseServer
	messages []string
}

func main() {
	
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Did not work")
	}	
	
	
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewITUDatabaseClient(conn)
	
	stream, err :=client.SendRecieve(context.Background(), line)
	if err != nil {
		log.Fatalf("Not working")
	}
	
	
	smessage, err := client.SendRecieve(context.Background(), &proto.ClientMessage{})
	if err != nil {
		log.Fatalf("Not working")
	}

	for _, smessage := range smessage.Smessages {
		log.Println(" - " + smessage)
	}

}
