package main

import "fmt"

func main() {
	repository := NewRepository()
	commandHandler := NewCommandHandler(repository)

	createUserCommand := &CreateUserCommand{"1", "Aidan Fewster", 20}
	err := commandHandler.handle(createUserCommand)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	increaseUsersAgeCommand := &IncreaseUsersAgeCommand{"1"}
	commandHandler.handle(increaseUsersAgeCommand)
	commandHandler.handle(increaseUsersAgeCommand)
	commandHandler.handle(increaseUsersAgeCommand)
	commandHandler.handle(increaseUsersAgeCommand)
	commandHandler.handle(increaseUsersAgeCommand)
	commandHandler.handle(increaseUsersAgeCommand)
	commandHandler.handle(increaseUsersAgeCommand)

	err = commandHandler.handle(increaseUsersAgeCommand)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	user, err := repository.GetUser("1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Users age: %v\n", user.GetAge())
}
