package server

import (
	"net/http"

	"github.com/fewstera/go-event-sourcing/eventsourcing"
)

var commandHandler *eventsourcing.CommandHandler
var repository *eventsourcing.Repository

func StartServer(ch *eventsourcing.CommandHandler, repo *eventsourcing.Repository) {
	commandHandler = ch
	repository = repo

	router := NewRouter()
	http.ListenAndServe(":8080", router)
}
