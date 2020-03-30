package wrapper

import (
	"../../proto"
	responseCode "../../util"
	"../KeyValueStore"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"strings"
	"time"
)

type Server struct{}

var CoordinatorClient proto.TwoPhaseServiceClient
var ConnectionToCoordinator *grpc.ClientConn

func CreateClientForCoordinator() {
	connection, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	ConnectionToCoordinator = connection

	// err occurred during establishing connection with the Server
	if err != nil {
		log.Fatalf("Error occurred while establishing the connection with the Server: %v ", err)
	}

	// Creating client to call the server's service
	CoordinatorClient = proto.NewTwoPhaseServiceClient(connection)
}

func checkForConnection() {

	for ConnectionToCoordinator.GetState() == connectivity.TransientFailure {
		log.Println("Server seems to be done. Let's Wait for 5 second and try connecting again...")
		time.Sleep(time.Second * 5)
		log.Println("Reconnecting....")
	}
}

// Client Functions

func (s *Server) PUT(ctx context.Context, request *proto.PutRequest) (*proto.Response, error) {

	log.Printf("PUT call is received from the client")

	checkForConnection()

	twoPhaseProtocolResponse, err := CoordinatorClient.InitiateTwoPhaseProtocol(context.Background(), &proto.TwoPhaseRequest{Key: request.Key, Value: request.Value, Operation: "PUT"})

	if err != nil {
		log.Fatalf("Error happened when two phase protocol is initiated... %v", err)
	}

	return &proto.Response{ResponseCode: twoPhaseProtocolResponse.OperationResponseCode, Message: twoPhaseProtocolResponse.OperationResponseMessage}, nil

}

func (s *Server) GET(ctx context.Context, request *proto.GetAndDeleteRequest) (*proto.Response, error) {

	key := request.Key

	if key == "" {
		return &proto.Response{ResponseCode: responseCode.INVALID_INPUT, Message: "Key cannot be null or empty. GET is not initiated"}, nil
	}

	value := KeyValueStore.GET(key)

	if value == "" {
		return &proto.Response{ResponseCode: responseCode.INVALID_INPUT, Message: "Value is not found for the given key"}, nil
	}

	return &proto.Response{ResponseCode: responseCode.SUCCESS, Message: fmt.Sprintf("GET has been successful. Value for the key %s is %s", key, value)}, nil

}

func (s *Server) DELETE(ctx context.Context, request *proto.GetAndDeleteRequest) (*proto.Response, error) {

	log.Printf("DELETE call is received from the client")

	checkForConnection()

	twoPhaseProtocolResponse, err := CoordinatorClient.InitiateTwoPhaseProtocol(context.Background(), &proto.TwoPhaseRequest{Key: request.Key, Value: "Empty Such Empty", Operation: "DELETE"})

	if err != nil {
		log.Fatalf("Error happened when two phase protocol is initiated... %v", err)
	}

	return &proto.Response{ResponseCode: twoPhaseProtocolResponse.OperationResponseCode, Message: twoPhaseProtocolResponse.OperationResponseMessage}, nil

}

// Coordinator Functions

func (s *Server) CanCommit(ctx context.Context, request *proto.PutRequest) (*proto.CanCommitResponse, error) {
	key := request.Key
	value := request.Value

	actualValue := KeyValueStore.GET(key)

	checkForConnection()

	if key == "" {
		return &proto.CanCommitResponse{CanCommit: false, Message: "Cannot be committed as the key is empty."}, nil
	}

	if value == actualValue {
		return &proto.CanCommitResponse{CanCommit: false, Message: "Cannot be committed as the key-value pair already exists."}, nil
	}

	return &proto.CanCommitResponse{CanCommit: true, Message: "I am ready for committment"}, nil
}

func (s *Server) Commit(ctx context.Context, request *proto.CommitRequest) (*proto.CommitACK, error) {
	key := request.Key
	value := request.Value

	if strings.ToUpper(request.Operation) == "PUT" {
		KeyValueStore.PUT(key, value)

		return &proto.CommitACK{
			OperationResponseCode:    200,
			OperationResponseMessage: fmt.Sprintf("PUT has been successful.  Key-Value pair %s-%s is updated into the DB successfully", request.Key, request.Value),
			IsCommitted:              true,
			Message:                  fmt.Sprintf("PUT has been successful.  Key-Value pair %s-%s is updated into the DB successfully", key, value)}, nil
	}

	originalValue := KeyValueStore.GET(key)
	KeyValueStore.DELETE(key)

	return &proto.CommitACK{
		IsCommitted:              true,
		Message:                  "Commit is successful",
		OperationResponseCode:    200,
		OperationResponseMessage: fmt.Sprintf("DELETE has been successful. Key-Value paid %s - %s no longer exist", key, originalValue)}, nil
}

func (s *Server) Abort(ctx context.Context, request *proto.Empty) (*proto.AbortACK, error) {
	return &proto.AbortACK{IsAborted: true, Message: "Operation aborted successfully"}, nil
}
