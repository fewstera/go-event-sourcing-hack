package main

import (
	"fmt"
	"reflect"
)

type CommandHandler struct {
	repository *Repository
}

func NewCommandHandler(repository *Repository) *CommandHandler {
	ch := new(CommandHandler)
	ch.repository = repository
	return ch
}

func (commandHandler *CommandHandler) handle(command Command) error {
	switch c := command.(type) {
	case *CreateUserCommand:
		return commandHandler.handleCreateUserCommand(c)
	case *IncreaseUsersAgeCommand:
		return commandHandler.handleIncreaseUsersAgeCommand(c)
	case *ChangeUsersNameCommand:
		return commandHandler.handleChangeUsersName(c)
	default:
		return &UnkownCommandError{fmt.Sprintf("Unkown command (%s) sent to command handler", reflect.TypeOf(c))}
	}
}

func (commandHandler *CommandHandler) handleCreateUserCommand(c *CreateUserCommand) error {
	user, err := NewUser(c.id, c.name, c.age)
	if err == nil {
		commandHandler.repository.SaveUser(user)
	}
	return err
}

func (commandHandler *CommandHandler) handleIncreaseUsersAgeCommand(c *IncreaseUsersAgeCommand) error {
	user, err := commandHandler.repository.GetUser(c.id)
	if err != nil {
		return err
	}

	user.IncreaseAge()
	commandHandler.repository.SaveUser(user)
	return nil
}

func (commandHandler *CommandHandler) handleChangeUsersName(c *ChangeUsersNameCommand) error {
	user, err := commandHandler.repository.GetUser(c.id)
	if err != nil {
		return err
	}

	user.ChangeName(c.newName)
	commandHandler.repository.SaveUser(user)
	return nil
}
