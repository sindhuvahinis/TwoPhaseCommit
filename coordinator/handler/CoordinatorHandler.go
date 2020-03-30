package handler

import (
	"../../proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"strconv"
	"strings"
)

type Coordinator struct {
}

type Connection struct {
	PortNumber       int64
	ClientForServer  proto.KeyValueStoreServiceClient
	ServerConnection *grpc.ClientConn
}

var Connections []Connection

func (c *Coordinator) JoinConnection(ctx context.Context, request *proto.JoinConnectionRequest) (*proto.EmptyTwoPhase, error) {

	log.Printf("A server is requested to connect..")
	portNumber := request.PortNumber

	// create a client for connection that is made.
	clientForServer, serverConnection := CreateClient(portNumber)
	Connections = append(Connections, Connection{PortNumber: portNumber, ClientForServer: clientForServer, ServerConnection: serverConnection})

	return &proto.EmptyTwoPhase{}, nil
}

func (c *Coordinator) InitiateTwoPhaseProtocol(ctx context.Context, request *proto.TwoPhaseRequest) (*proto.TwoPhaseResponse, error) {

	log.Printf("Initiated Two Phase protocol... \n")

	var accumulate = true

	log.Printf("The Connections %v are ", Connections)
	// Iterate through all the connections and collect canCommit
	for _, connect := range Connections {

		if !isConnectionAlive(connect.ServerConnection) {
			accumulate = false
			continue
		}

		canCommitResponse, err := connect.ClientForServer.CanCommit(context.Background(), &proto.PutRequest{Key: request.Key, Value: request.Value})

		if err != nil {
			log.Fatalf("Error happened in the can commit phase %v ", err)
		}

		log.Printf("CanCommit response from connection port number %d is %v and %s", connect.PortNumber, canCommitResponse.CanCommit, canCommitResponse.Message)
		accumulate = accumulate && canCommitResponse.CanCommit
	}

	if !accumulate {
		for _, connect := range Connections {

			if !isConnectionAlive(connect.ServerConnection) {
				continue
			}

			abortResponse, err := connect.ClientForServer.Abort(context.Background(), &proto.Empty{})

			if err != nil {
				log.Fatalf("Error happened in the abort phase %v ", err)
			}

			log.Printf("Abort response from connection port number %d is %v and %s", connect.PortNumber, abortResponse.IsAborted, abortResponse.Message)
		}

		return &proto.TwoPhaseResponse{OperationResponseCode: 400, OperationResponseMessage: "PUT/DELETE operation cannot be done, as in two phase protocol some voting are against this"}, nil

	} else {

		for _, connect := range Connections {

			if !isConnectionAlive(connect.ServerConnection) {
				continue
			}

			commitResponse, err := connect.ClientForServer.Commit(context.Background(), &proto.CommitRequest{Key: request.Key, Value: request.Value, Operation: request.Operation})

			if err != nil {
				log.Fatalf("Error happened in the commit phase %v ", err)
			}

			log.Printf("Commit response from connection port number %d is %v and %s", connect.PortNumber, commitResponse.IsCommitted, commitResponse.Message)
		}
	}

	var responseMessage string

	if strings.ToUpper(request.Operation) == "PUT" {
		responseMessage = fmt.Sprintf("PUT has been successful.  Key-Value pair %s-%s is updated into the DB successfully", request.Key, request.Value)
	} else {
		responseMessage = fmt.Sprintf("DELETE has been successful. value for the key %s is %s.", request.Key, request.Value)
	}

	return &proto.TwoPhaseResponse{OperationResponseCode: 200, OperationResponseMessage: responseMessage}, nil
}

func CreateClient(portNumber int64) (proto.KeyValueStoreServiceClient, *grpc.ClientConn) {

	target := "localhost:" + strconv.FormatInt(portNumber, 10)

	log.Printf("target name is %s ", target)
	connection, err := grpc.Dial(target, grpc.WithInsecure())

	// err occurred during establishing connection with the Server
	if err != nil {
		log.Fatalf("Error occurred while establishing the connection with the Server: %v ", err)
	}

	// Creating client to call the server's service
	return proto.NewKeyValueStoreServiceClient(connection), connection

}

func isConnectionAlive(connection *grpc.ClientConn) bool {

	if connection.GetState() == connectivity.TransientFailure {
		log.Println("Server seems to be down")
		return false
	}
	return true
}
