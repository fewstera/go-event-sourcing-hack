package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
	"github.com/fewstera/go-event-sourcing-hack/pkg/user"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type CreateUserPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (s *Server) GetUserRouteHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	usr, err := s.projection.GetUser(userId)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usr); err != nil {
		writeErrorResponse(err, w)
	}
}

func (s *Server) PostUserRouteHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		writeErrorResponse(err, w)
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		writeErrorResponse(err, w)
		return
	}

	streamID := uuid.NewV4().String()
	createUserCommand := user.NewCreateUserCommand(streamID, payload.Name, payload.Age)
	event, err := s.cmdHandler.Handle(createUserCommand)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeEventJsonResponse(event, w)
}

func writeEventJsonResponse(event eventstore.Event, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		writeErrorResponse(err, w)
	}
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	switch err.(type) {
	case *user.UserNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "UserNotFoundError: %v\n", err)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
	}
}
