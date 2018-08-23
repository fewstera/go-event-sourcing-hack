package main

import "fmt"

func main() {
	repository := NewRepository()
	commandHandler := NewCommandHandler(repository)

	createUserCommand := &CreateUserCommand{"1", "Aidan Fewster", 25}
	handleCommandOrPanic(commandHandler, createUserCommand)

	increaseUsersAgeCommand := &IncreaseUsersAgeCommand{"1"}
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)
	handleCommandOrPanic(commandHandler, increaseUsersAgeCommand)

	userNameChangeCommand := &ChangeUsersNameCommand{"1", "Bob Smith"}
	handleCommandOrPanic(commandHandler, userNameChangeCommand)

	user, _ := repository.GetUser("1")
	fmt.Printf("Users name: %v\n", user.GetName())
	fmt.Printf("Users age: %v\n", user.GetAge())
}

func handleCommandOrPanic(commandHandler *CommandHandler, command Command) {
	err := commandHandler.handle(command)
	if err != nil {
		panic(err)
	}
}
