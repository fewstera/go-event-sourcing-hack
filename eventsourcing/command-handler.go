package eventsourcing

import (
	"fmt"
	"reflect"
)

type CommandHandler struct {
	eventStore *EventStore
	projection *Projection
}

func NewCommandHandler(eventStore *EventStore, projection *Projection) *CommandHandler {
	ch := new(CommandHandler)
	ch.eventStore = eventStore
	ch.projection = projection
	return ch
}

func (commandHandler *CommandHandler) Handle(command Command) (Event, error) {
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

func (commandHandler *CommandHandler) handleCreateUserCommand(c *CreateUserCommand) (Event, error) {
	event, err := NewUser(c.Id, c.Name, c.Age)
	if err != nil {
		return nil, err
	}
	err = commandHandler.eventStore.SaveEvent(event)
	if err != nil {
		return nil, err
	}

	fmt.Println("Saved user event")
	return event, nil
}

func (commandHandler *CommandHandler) handleIncreaseUsersAgeCommand(c *IncreaseUsersAgeCommand) (Event, error) {
	user, err := commandHandler.projection.GetUser(c.Id)
	if err != nil {
		return nil, err
	}

	event := user.IncreaseAge()

	err = commandHandler.eventStore.SaveEvent(event)
	if err != nil {
		return nil, err
	}

	fmt.Println("Saved user event")

	return event, nil
}

func (commandHandler *CommandHandler) handleChangeUsersName(c *ChangeUsersNameCommand) (Event, error) {
	user, err := commandHandler.projection.GetUser(c.Id)
	if err != nil {
		return nil, err
	}

	event := user.ChangeName(c.NewName)
	err = commandHandler.eventStore.SaveEvent(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
