package server

import (
	"net/http"

	"github.com/fewstera/go-event-sourcing/eventsourcing"
)

var commandHandler *eventsourcing.CommandHandler
var projection *eventsourcing.Projection

func StartServer(ch *eventsourcing.CommandHandler, prjction *eventsourcing.Projection) {
	commandHandler = ch
	projection = prjction

	router := NewRouter()
	http.ListenAndServe(":8080", router)
}
