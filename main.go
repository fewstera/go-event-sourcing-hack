package main

import (
	"database/sql"
	"fmt"

	"github.com/fewstera/go-event-sourcing/eventsourcing"
	"github.com/fewstera/go-event-sourcing/server"
)

import _ "github.com/go-sql-driver/mysql"

func main() {
	db := initDb()

	eventStore := eventsourcing.NewEventStore(db)
	repository := eventsourcing.NewRepository(eventStore)
	commandHandler := eventsourcing.NewCommandHandler(repository)

	createUserCommand := &eventsourcing.CreateUserCommand{"1", "Aidan Fewster", 25}
	handleCommandOrPanic(commandHandler, createUserCommand)

	increaseUsersAgeCommand := &eventsourcing.IncreaseUsersAgeCommand{"1"}
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)

	userNameChangeCommand := &eventsourcing.ChangeUsersNameCommand{"1", "Bob Smith"}
	handleCommandOrPanic(commandHandler, userNameChangeCommand)

	user, _ := repository.GetUser("1")
	fmt.Printf("Users name: %v\n", user.GetName())
	fmt.Printf("Users age: %v\n", user.GetAge())

	server.StartServer(commandHandler, repository)
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
