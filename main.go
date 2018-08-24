package main

import (
	"fmt"

	"github.com/fewstera/go-event-sourcing/eventsourcing"
	"github.com/fewstera/go-event-sourcing/server"
)

func main() {
	repository := eventsourcing.NewRepository()
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
