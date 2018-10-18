package user

import (
	"fmt"
	"reflect"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

type CommandHandler struct {
	eventstore eventstore.EventStore
	projection *Projection
}

func NewCommandHandler(eventstore eventstore.EventStore, projection *Projection) *CommandHandler {
	ch := new(CommandHandler)
	ch.eventstore = eventstore
	ch.projection = projection
	return ch
}

func (ch *CommandHandler) Handle(c Command) (eventstore.Event, error) {
	switch c := c.(type) {
	case CreateUserCommand:
		return ch.handleCreateUserCommand(c)
	default:
		return nil, &UnkownCommandError{fmt.Sprintf("Unkown command (%s) sent to command handler", reflect.TypeOf(c))}
	}
}

func (ch *CommandHandler) handleCreateUserCommand(c CreateUserCommand) (eventstore.Event, error) {
	event, err := NewUser(c.StreamID, c.Name, c.Age)
	if err != nil {
		return nil, err
	}
	err = ch.eventstore.SaveEvent(event)
	if err != nil {
		return nil, err
	}

	fmt.Println("Saved user event")
	return event, nil
}
