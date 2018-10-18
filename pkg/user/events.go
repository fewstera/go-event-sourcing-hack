package user

import (
	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

// Event types
const (
	EventTypeUserCreated string = "USER_CREATED"
	EventTypeDeposited   string = "DEPOSITED"
	EventTypeWithdrawn   string = "WITHDRAWN"
)

// EmptyEventCreators returns a mapping of EventTypes to funcs that return
// empty instance of those events.
//
// The map is sent to the event factory to allow event instances to be created
// from event data.
func EmptyEventCreators() eventstore.EmptyEventCreatorsMap {
	return eventstore.EmptyEventCreatorsMap{
		EventTypeUserCreated: func() eventstore.Event { return &UserCreatedEvent{} },
		EventTypeDeposited:   func() eventstore.Event { return &DepositedEvent{} },
		EventTypeWithdrawn:   func() eventstore.Event { return &WithdrawnEvent{} },
	}
}
