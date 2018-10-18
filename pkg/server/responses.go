package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
	"github.com/fewstera/go-event-sourcing-hack/pkg/user"
)

func writeEventJsonResponse(event eventstore.Event, w http.ResponseWriter) {
	eventData, err := event.Data()
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(eventData); err != nil {
		writeErrorResponse(err, w)
	}
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	switch err.(type) {
	case *user.UserNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "UserNotFoundError: %v\n", err)
	case *user.InsufficientFundsError:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Insufficient funds: %v\n", err)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v\n", err)
	}
}
