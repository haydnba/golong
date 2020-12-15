package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Key - DDB key type
type Key struct {
	ConnectionID string
}

// Response - for readability
type Response = events.APIGatewayProxyResponse

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayWebsocketProxyRequest) (Response, error) {
	// Log execution
	log.Print("Processing wsapi client disconnect")
	log.Printf("Request ID: %s\n", req.RequestContext.RequestID)
	log.Printf("Connection ID: %s\n", req.RequestContext.ConnectionID)

	// Initialize session using ~/.aws/credentials & ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Register the DDB client
	svc := dynamodb.New(sess)

	// Construct the `ConnectionID` storage key
	key := Key{ConnectionID: req.RequestContext.ConnectionID}
	av, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		log.Fatalf("Error processing connection id: %v", err)
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	// Construct the `PutItemInput`
	table := os.Getenv("TABLE_NAME")
	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(table),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Error storing connection id: %v", err)
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	// Return OK
	return Response{StatusCode: http.StatusOK}, nil
}
