package main

import (
	"../proto"
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	ipAddress := os.Args[1]
	portNumber := os.Args[2]

	target := ipAddress + ":" + portNumber

	// Connection to the Server
	connection, err := grpc.Dial(target, grpc.WithInsecure())

	// err occurred during establishing connection with the Server
	if err != nil {
		log.Fatalf("Error occurred while establishing the connection with the Server: %v ", err)
	}

	// close the connection when the function is exiting..
	defer connection.Close()

	// Creating client to call the server's service
	client := proto.NewKeyValueStoreServiceClient(connection)

	// Call the user input handler
	UserInputHandler(client, connection)

}

//UserInputHandler for handling the user input values.
func UserInputHandler(client proto.KeyValueStoreServiceClient, connection *grpc.ClientConn) {
	reader := bufio.NewReader(os.Stdin)

	// Infinite loop
	for {
		fmt.Println("Please Enter one of the following operations: \nPUT/GET/DELETE :")
		userInput, _ := reader.ReadString('\n')

		switch formatUpper(userInput) {
		case "PUT":
			{
				fmt.Println("Please enter a Key:")
				key, _:= reader.ReadString('\n')

				fmt.Println("Please enter a Value")
				value, _:= reader.ReadString('\n')

				checkForConnection(connection)
				response, err := client.PUT(context.Background(), &proto.PutRequest{Key: format(key), Value: format(value)})

				if err != nil {
					log.Fatalf("Error occurred when calling PUT %v ", err)
				}

				log.Printf("\nResponse is received from server "+
					"\nResponse Code : %d , \nResponse Message : %s \n", response.ResponseCode, response.Message)
			}
		case "GET":
			{
				fmt.Println("Please enter the Key:")
				key, _:= reader.ReadString('\n')

				checkForConnection(connection)
				response, err := client.GET(context.Background(), &proto.GetAndDeleteRequest{Key: format(key)})

				if err != nil {
					log.Fatalf("Error occurred when calling GET %v ", err)
				}

				log.Printf("\nResponse is received from server "+
					"\nResponse Code : %d , \nResponse Message : %s \n", response.ResponseCode, response.Message)
			}
		case "DELETE":
			{
				fmt.Println("Please enter the Key:")
				key, _:= reader.ReadString('\n')

				response, err := client.DELETE(context.Background(), &proto.GetAndDeleteRequest{Key: format(key)})

				if err != nil {
					log.Fatalf("Error occurred when calling GET %v ", err)
				}

				log.Printf("\nResponse is received from server "+
					"\nResponse Code : %d  \nResponse Message : %s \n", response.ResponseCode, response.Message)
			}
		default:
			{
				log.Printf("\nOperation entered in invalid.\n")
			}
		}

	}

}

func checkForConnection(connection *grpc.ClientConn) {

	for connection.GetState() == connectivity.TransientFailure {
		log.Println("Server seems to be done. Let's Wait for 5 second and try connecting again...")
		time.Sleep(time.Second * 5)
		log.Println("Reconnecting....")
	}
}

func format(text string) string {
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func formatUpper(text string) string {
	text = strings.Replace(text, "\n", "", -1)
	return strings.ToUpper(text)
}
