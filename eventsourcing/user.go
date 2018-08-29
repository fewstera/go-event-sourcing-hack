package eventsourcing

import (
	"encoding/json"
	"fmt"
)

type User struct {
	eventNumber      int
	id               string
	name             string
	age              int
	uncommitedEvents []Event
}

// Constuctor
func NewUser(id string, name string, age int) (*User, error) {
	u := new(User)
	if age < 0 {
		return nil, &InvalidAgeError{"Age is negative"}
	}
	u.apply(NewUserCreatedEvent(1, id, name, age))
	return u, nil
}

// Actions - These are called by command handlers, they can error but should not mutate state (except through calling apply)
func (u *User) IncreaseAge() {
	nextEventNumber := u.eventNumber + 1
	u.apply(NewUserGotOlderEvent(nextEventNumber, u.id))
}

func (u *User) ChangeName(name string) {
	nextEventNumber := u.eventNumber + 1
	u.apply(NewUsersNameChangedEvent(nextEventNumber, u.id, name))
}

// Apply methods - These should only mutate state, they are not allowed to error.
func (u *User) apply(event Event) {
	u.eventNumber = event.GetEventNumber()

	switch e := event.(type) {
	case *UserCreatedEvent:
		u.applyUserCreated(e)
	case *UserGotOlderEvent:
		u.applyUserGotOlder(e)
	case *UsersNameChangedEvent:
		u.applyUsersNameChanged(e)
	default:
		fmt.Println("Unkown event applied on user")
	}
}

func (u *User) applyUserCreated(e *UserCreatedEvent) {
	u.id = e.Id
	u.age = e.Age
	u.name = e.Name

	u.uncommitedEvents = append(u.uncommitedEvents, e)
}

func (u *User) applyUserGotOlder(e *UserGotOlderEvent) {
	u.age = u.age + 1

	u.uncommitedEvents = append(u.uncommitedEvents, e)
}

func (u *User) applyUsersNameChanged(e *UsersNameChangedEvent) {
	u.name = e.NewName

	u.uncommitedEvents = append(u.uncommitedEvents, e)
}

// Commited events methods
func (u *User) MarkChangesAsCommitted() {
	u.uncommitedEvents = nil
}

func (u *User) GetUncommitedEvents() []Event {
	return u.uncommitedEvents
}

// Getters
func (u *User) GetId() string {
	return u.id
}

func (u *User) GetAge() int {
	return u.age
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Id:   u.id,
		Age:  u.age,
		Name: u.name,
	})
}
