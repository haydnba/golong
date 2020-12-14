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

// Item - DDB item type
type Item struct {
	ConnectionID string
}

// Response - for readability
type Response = events.APIGatewayProxyResponse

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayWebsocketProxyRequest) (Response, error) {
	// Log execution
	log.Printf("Processing Lambda request %s\n", req.RequestContext.RequestID)
	log.Printf("Processing Lambda request %s\n", req.RequestContext.ConnectionID)

	// Initialize session using ~/.aws/credentials & ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Register the DDB client
	svc := dynamodb.New(sess)

	// Construct the `ConnectionID` storage item
	item := Item{ConnectionID: req.RequestContext.ConnectionID}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Error processing connection id: %v", err)
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	// Construct the `PutItemInput`
	table := os.Getenv("TABLE_NAME")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Error storing connection id: %v", err)
		return Response{StatusCode: http.StatusInternalServerError}, err
	}

	// Return OK
	return Response{StatusCode: http.StatusOK}, nil
}
