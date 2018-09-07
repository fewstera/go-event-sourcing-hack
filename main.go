package main

import (
	"database/sql"
	"fmt"

	"github.com/fewstera/go-event-sourcing-hack/eventsourcing"
	"github.com/fewstera/go-event-sourcing-hack/server"
)

import _ "github.com/go-sql-driver/mysql"

func main() {
	var projection *eventsourcing.Projection
	db := initDb()

	snapshots := eventsourcing.NewSnapshots(db)
	snapshotState := snapshots.GetStateFromLatestSnapshot()
	if snapshotState.Projection != nil {
		projection = snapshotState.Projection
		fmt.Printf("Start up from a snapshot. Current position %v\n", snapshotState.Position)
	} else {
		projection = eventsourcing.NewProjection()
	}

	eventStore := eventsourcing.NewEventStore(db, projection, snapshots, snapshotState.Position)
	commandHandler := eventsourcing.NewCommandHandler(eventStore, projection)

	fmt.Println("Starting HTTP server")
	server.StartServer(commandHandler, projection)
}

func handleCommandOrPanic(commandHandler *eventsourcing.CommandHandler, command eventsourcing.Command) {
	_, err := commandHandler.Handle(command)
	if err != nil {
		panic(err)
	}
}

func initDb() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/events")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}
