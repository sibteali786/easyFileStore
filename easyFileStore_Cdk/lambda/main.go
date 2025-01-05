package main

import (
	"fmt"
	"lambda-func/app"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Take in a payload and do something with it maybe
func HandleRequest(event MyEvent) (string, error) {

	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Successfully called by - %v\n", event.Username), nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterApiHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginHandler(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: 404,
			}, nil
		}
	})
}
