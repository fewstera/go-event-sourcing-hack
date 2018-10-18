package user

import (
	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

// Event types
const (
	EventTypeUserCreated string = "USER_CREATED"
)

// EmptyEventCreators returns a mapping of EventTypes to funcs that return
// empty instance of those events.
//
// The map is sent to the event factory to allow event instances to be created
// from event data.
func EmptyEventCreators() eventstore.EmptyEventCreatorsMap {
	return eventstore.EmptyEventCreatorsMap{
		EventTypeUserCreated: func() eventstore.Event { return &UserCreatedEvent{} },
	}
}
