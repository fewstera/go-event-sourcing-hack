package eventsourcing

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

func (commandHandler *CommandHandler) Handle(command Command) (*User, error) {
	switch c := command.(type) {
	case *CreateUserCommand:
		return commandHandler.handleCreateUserCommand(c)
	case *IncreaseUsersAgeCommand:
		return commandHandler.handleIncreaseUsersAgeCommand(c)
	case *ChangeUsersNameCommand:
		return commandHandler.handleChangeUsersName(c)
	default:
		return nil, &UnkownCommandError{fmt.Sprintf("Unkown command (%s) sent to command handler", reflect.TypeOf(c))}
	}
}

func (commandHandler *CommandHandler) handleCreateUserCommand(c *CreateUserCommand) (*User, error) {
	user, err := NewUser(c.Id, c.Name, c.Age)
	if err != nil {
		return nil, err
	}
	commandHandler.repository.SaveUser(user)
	return user, nil
}

func (commandHandler *CommandHandler) handleIncreaseUsersAgeCommand(c *IncreaseUsersAgeCommand) (*User, error) {
	user, err := commandHandler.repository.GetUser(c.Id)
	if err != nil {
		return nil, err
	}

	user.IncreaseAge()
	commandHandler.repository.SaveUser(user)
	return user, nil
}

func (commandHandler *CommandHandler) handleChangeUsersName(c *ChangeUsersNameCommand) (*User, error) {
	user, err := commandHandler.repository.GetUser(c.Id)
	if err != nil {
		return nil, err
	}

	user.ChangeName(c.NewName)
	commandHandler.repository.SaveUser(user)
	return user, nil
}
