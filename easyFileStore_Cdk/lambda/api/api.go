package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterApiHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}

	// does user exists
	userExists, err := api.dbStore.DoesUserExists(event.Username)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if userExists {
		return fmt.Errorf("user already exists")
	}

	// we know that does not exists
	err = api.dbStore.InsertUser(event)
	if err != nil {
		fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}
