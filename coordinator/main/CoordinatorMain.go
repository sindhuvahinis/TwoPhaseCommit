package main

import (
	"../../proto"
	"../handler"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	//creating server instance
	coordinatorServer := handler.Coordinator{}

	proto.RegisterTwoPhaseServiceServer(grpcServer, &coordinatorServer)
	reflection.Register(grpcServer)

	log.Printf("\nCoordinator is up and running in the port 9000.")
	log.Printf("\nWaiting for servers to connect..")

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
