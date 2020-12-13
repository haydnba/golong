package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response - will export to lib soon
type Response = events.APIGatewayProxyResponse

func main() {
	lambda.Start(handler)
}

func handler(req *events.APIGatewayWebsocketProxyRequest) (Response, error) {
	// Just do this for now...
	if false {
		return Response{StatusCode: http.StatusInternalServerError}, nil
	}

	// Log some stuff
	log.Printf("Processing Lambda request %s\n", req.RequestContext.RequestID)
	log.Printf("Processing Lambda request %s\n", req.RequestContext.ConnectionID)

	// Return OK
	return Response{StatusCode: http.StatusOK}, nil
}
