package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterApiHandler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(event.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request - empty fields",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	// does user exists
	userExists, err := api.dbStore.DoesUserExists(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("could not create a new user %v", err)
	}

	// we know that user does not exists
	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user - %v", err)
	}
	return events.APIGatewayProxyResponse{
		Body:       "User registered successfully",
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginRequest LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// validate password with one stored in database
	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid User Credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Successfully logged in",
		StatusCode: http.StatusOK,
	}, nil
}
