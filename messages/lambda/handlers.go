package main

import (
	"encoding/json"
	"errors"

	// "fmt"

	"github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/kewlamogh/kidgo-backend/messages"
)

type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	if request.Resource == "/inbox" {
		username := request.QueryStringParameters["username"]
		classroom := request.QueryStringParameters["classroom"]
		return getInboxHandler(request, username, classroom)
	} else if request.Resource == "/outbox" {
		username := request.QueryStringParameters["username"]
		classroom := request.QueryStringParameters["classroom"]
	return getOutboxHandler(request, username, classroom)
	} else if request.Resource == "/send" {
		from := request.QueryStringParameters["from"]
		to := request.QueryStringParameters["to"]
		s3url := request.QueryStringParameters["s3url"]
		classroom := request.QueryStringParameters["classroom"]
		return sendMessage(request, from, to, s3url, classroom)
	} else {
		return Response{}, errors.New("Invalid thingy")
	}
}

func getInboxHandler(request events.APIGatewayProxyRequest, username string, classroom string) (Response, error) {
	inboxItems, err := messages.GetUserInbox(username, classroom)

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(inboxItems)
	jsonStr := string(jsonBytes)

	// You would typically serialize your response to JSON here.
	response := jsonStr

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            response,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}, nil
}

func getOutboxHandler(request events.APIGatewayProxyRequest, username string, classroom string) (Response, error) {
	inboxItems, err := messages.GetUserOutbox(username, classroom)

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(inboxItems)
	jsonStr := string(jsonBytes)

	// You would typically serialize your response to JSON here.
	response := jsonStr

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            response,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}, nil
}

func sendMessage(request events.APIGatewayProxyRequest, from, to, s3url, classID string) (Response, error) {
	err := messages.SendMessage(from, to, s3url, classID)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "the deed is done",
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}, nil
}
