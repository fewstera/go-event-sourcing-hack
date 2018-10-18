package user

import (
	"fmt"
	"sync"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

type User struct {
	EventNumber int    `json:"version"`
	StreamID    string `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	mutex       sync.RWMutex
}

// Constuctor
func NewUser(streamID string, name string, age int) (*UserCreatedEvent, error) {
	if age < 0 {
		return nil, &InvalidAgeError{"Age is negative"}
	}

	return NewUserCreatedEvent(streamID, 1, name, age), nil
}

// Apply methods - These should only mutate state, they are not allowed to error.
func (u *User) Apply(event eventstore.Event) {
	u.EventNumber = event.GetEventNumber()

	switch e := event.(type) {
	case *UserCreatedEvent:
		u.applyUserCreated(e)
	default:
		fmt.Println("Unkown event applied on user")
	}
}

func (u *User) applyUserCreated(e *UserCreatedEvent) {
	u.mutex.Lock()
	u.StreamID = e.StreamID
	u.Age = e.Age
	u.Name = e.Name
	u.mutex.Unlock()
}
