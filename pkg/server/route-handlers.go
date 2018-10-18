package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fewstera/go-event-sourcing-hack/pkg/user"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

func (s *Server) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := s.projection.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		writeErrorResponse(err, w)
	}
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	streamID := mux.Vars(r)["id"]
	usr, err := s.projection.GetUser(streamID)
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

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

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

func (s *Server) depositHandler(w http.ResponseWriter, r *http.Request) {
	streamID := mux.Vars(r)["id"]
	payload := struct {
		Version int     `json:"version"`
		Amount  float32 `json:"amount"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&payload)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	cmd := user.NewDepositCommand(streamID, payload.Version, payload.Amount)
	event, err := s.cmdHandler.Handle(cmd)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeEventJsonResponse(event, w)
}

func (s *Server) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	streamID := mux.Vars(r)["id"]
	payload := struct {
		Version int     `json:"version"`
		Amount  float32 `json:"amount"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&payload)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	cmd := user.NewWithdrawCommand(streamID, payload.Version, payload.Amount)
	event, err := s.cmdHandler.Handle(cmd)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	writeEventJsonResponse(event, w)
}
