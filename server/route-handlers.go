package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fewstera/go-event-sourcing/eventsourcing"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type CreateUserPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetUserRouteHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	user, err := repository.GetUser(userId)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeUserJsonResponse(user, w)
}

func PostUserRouteHandler(w http.ResponseWriter, r *http.Request) {
	var createUserPayload CreateUserPayload
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		writeErrorResponse(err, w)
		return
	}
	if err := json.Unmarshal(body, &createUserPayload); err != nil {
		writeErrorResponse(err, w)
		return
	}

	id := uuid.Must(uuid.NewV4()).String()
	createUserCommand := &eventsourcing.CreateUserCommand{id, createUserPayload.Name, createUserPayload.Age}
	user, err := commandHandler.Handle(createUserCommand)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeUserJsonResponse(user, w)
}

func PatchUserRouteHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	var patchMap map[string]*json.RawMessage
	if err := json.Unmarshal(body, &patchMap); err != nil {
		writeErrorResponse(err, w)
		return
	}

	nameJson, patchHasName := patchMap["name"]
	if !patchHasName || len(patchMap) != 1 {
		w.WriteHeader(422)
		fmt.Fprintf(w, "Only name value is patchable")
		return
	}

	var newName string
	if err := json.Unmarshal(*nameJson, &newName); err != nil {
		writeErrorResponse(err, w)
		return
	}

	userId := mux.Vars(r)["id"]
	changeUsersNameCommand := &eventsourcing.ChangeUsersNameCommand{userId, newName}
	user, err := commandHandler.Handle(changeUsersNameCommand)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeUserJsonResponse(user, w)
}

func PostIncreaseUserAgeRouteHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	increaseUsersAgeCommand := &eventsourcing.IncreaseUsersAgeCommand{userId}
	user, err := commandHandler.Handle(increaseUsersAgeCommand)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeUserJsonResponse(user, w)
}

func writeUserJsonResponse(user *eventsourcing.User, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeErrorResponse(err, w)
	}
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	switch err.(type) {
	case *eventsourcing.UserNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "UserNotFoundError: %v\n", err)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
	}
}
