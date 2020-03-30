package main

import (
	proto "../../proto"
	sw "../wrapper"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	// create a listener on TCP port 7777
	portNumber, _ := strconv.Atoi(os.Args[1])
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Connecting with Coordinator
	sw.CreateClientForCoordinator()

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	//creating server instance
	s := sw.Server{}

	proto.RegisterKeyValueStoreServiceServer(grpcServer, &s)
	reflection.Register(grpcServer)

	log.Printf("Going to connect with Coordinator...")
	// Sending the port number to Coordinator.
	_, err = sw.CoordinatorClient.JoinConnection(context.Background(), &proto.JoinConnectionRequest{PortNumber: int64(portNumber)})

	if err != nil {
		log.Fatalf("Could not connect to Coordinator: %v", err)
	}

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
