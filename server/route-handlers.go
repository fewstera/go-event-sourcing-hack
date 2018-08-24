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
		switch err.(type) {
		case *eventsourcing.UserNotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "UserNotFoundError: %v\n", err)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

func GetUsersRouteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Multiple users")
}

func PostUserRouteHandler(w http.ResponseWriter, r *http.Request) {
	var createUserPayload CreateUserPayload
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &createUserPayload); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	id := uuid.Must(uuid.NewV4()).String()
	createUserCommand := &eventsourcing.CreateUserCommand{id, createUserPayload.Name, createUserPayload.Age}
	user, err := commandHandler.Handle(createUserCommand)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}
