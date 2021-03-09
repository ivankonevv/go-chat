package main

import (
	"gRPCChatServer/chatserver"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" //default Port set if Port is not set in env
	}

	//init Listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}

	log.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	cs := chatserver.ChatServer{}
	chatserver.RegisterServiceServer(grpcserver, &cs)

	//gRPC listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}
}
