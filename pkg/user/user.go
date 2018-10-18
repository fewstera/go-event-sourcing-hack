package user

import (
	"fmt"
	"sync"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

// User is for representing a User.
type User struct {
	EventNumber int     `json:"version"`
	StreamID    string  `json:"id"`
	Name        string  `json:"name"`
	Age         int     `json:"age"`
	Balance     float32 `json:"balance"`
	mutex       sync.RWMutex
}

// Apply applies a given event to the user.
func (u *User) Apply(event eventstore.Event) {
	u.EventNumber = event.GetEventNumber()

	switch e := event.(type) {
	case *UserCreatedEvent:
		u.applyUserCreated(e)
	case *DepositedEvent:
		u.applyDeposited(e)
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

func (u *User) applyDeposited(e *DepositedEvent) {
	u.mutex.Lock()
	u.Balance = u.Balance + e.Amount
	u.mutex.Unlock()
}

// Create returns a UserCreatedEvent when validation passes and an error otherwise.
func (u *User) Create(streamID string, name string, age int) (*UserCreatedEvent, error) {
	if age < 0 {
		return nil, &InvalidAgeError{"Age is negative"}
	}

	return NewUserCreatedEvent(streamID, 1, name, age), nil
}

// Deposit returns a new DepositedEvent if validation passes.
func (u *User) Deposit(eventNumber int, amount float32) (*DepositedEvent, error) {
	if eventNumber != u.EventNumber {
		return nil, &EventNumberSyncError{u.EventNumber, eventNumber}
	}

	nextEventNumber := eventNumber + 1

	return NewDepositedEvent(u.StreamID, nextEventNumber, amount), nil
}
