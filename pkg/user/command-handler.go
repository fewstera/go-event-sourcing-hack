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
		return ch.handleCreateUser(c)
	case DepositCommand:
		return ch.handleDeposit(c)
	default:
		return nil, &UnkownCommandError{fmt.Sprintf("Unkown command (%s) sent to command handler", reflect.TypeOf(c))}
	}
}

func (ch *CommandHandler) handleCreateUser(c CreateUserCommand) (eventstore.Event, error) {
	user := &User{}
	event, err := user.Create(c.StreamID, c.Name, c.Age)
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

func (ch *CommandHandler) handleDeposit(c DepositCommand) (eventstore.Event, error) {
	usr, err := ch.projection.GetUser(c.StreamID)
	if err != nil {
		return nil, err
	}

	event, err := usr.Deposit(c.ClientEventNumber, c.Amount)
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
