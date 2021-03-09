package main

import (
	"bufio"
	"context"
	"fmt"
	"gRPCChatServer/chatserver"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Enter Server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}

	serverID = strings.Trim(serverID, "\r\n")

	log.Printf("Connecting: " + serverID)

	//connect to gRPC server
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server :: %v", err)
	}
	defer conn.Close()

	//call ChatService to create a stream
	client := chatserver.NewServiceClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	//implement communication with gRPC server
	ch := clientHandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	//blocker
	bl := make(chan bool)
	<- bl
}

type clientHandle struct {
	stream chatserver.Service_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console :: %v", err)
	}

	ch.clientName = strings.Trim(name, "\r\n")
}

func (ch *clientHandle) sendMessage() {
	//create a loop
	for {
		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console :: %v", err)
		}

		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)
		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}
	}
}

func (ch *clientHandle) receiveMessage() {
	//create a loop
	for{
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}
		//print message to console
		fmt.Printf("%s: %s", mssg.Name, mssg.Body)
	}
}