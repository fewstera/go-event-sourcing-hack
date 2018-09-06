package eventsourcing

import (
	"encoding/json"
	"fmt"
	"sync"
)

type User struct {
	eventNumber int
	id          string
	name        string
	age         int
	sync.RWMutex
}

// Constuctor
func NewUser(id string, name string, age int) (*UserCreatedEvent, error) {
	if age < 0 {
		return nil, &InvalidAgeError{"Age is negative"}
	}

	return NewUserCreatedEvent(1, id, name, age), nil
}

// Actions - These are called by command handlers, they can error but should not mutate state (except through calling apply)
func (u *User) IncreaseAge() *UserGotOlderEvent {
	u.RLock()
	nextEventNumber := u.eventNumber + 1
	event := NewUserGotOlderEvent(nextEventNumber, u.id)
	u.RUnlock()
	return event
}

func (u *User) ChangeName(name string) *UserNameChangedEvent {
	u.RLock()
	nextEventNumber := u.eventNumber + 1
	event := NewUserNameChangedEvent(nextEventNumber, u.id, name)
	u.RUnlock()
	return event
}

// Apply methods - These should only mutate state, they are not allowed to error.
func (u *User) Apply(event Event) {
	u.eventNumber = event.GetEventNumber()

	switch e := event.(type) {
	case *UserCreatedEvent:
		u.applyUserCreated(e)
	case *UserGotOlderEvent:
		u.applyUserGotOlder(e)
	case *UserNameChangedEvent:
		u.applyUsersNameChanged(e)
	default:
		fmt.Println("Unkown event applied on user")
	}
}

func (u *User) applyUserCreated(e *UserCreatedEvent) {
	u.Lock()
	u.id = e.Id
	u.age = e.Age
	u.name = e.Name
	u.Unlock()
}

func (u *User) applyUserGotOlder(e *UserGotOlderEvent) {
	u.Lock()
	u.age = u.age + 1
	u.Unlock()
}

func (u *User) applyUsersNameChanged(e *UserNameChangedEvent) {
	u.Lock()
	u.name = e.NewName
	u.Unlock()
}

// Getters
func (u *User) GetId() string {
	u.RLock()
	id := u.id
	u.RUnlock()
	return id
}

func (u *User) GetAge() int {
	u.RLock()
	age := u.age
	u.RUnlock()
	return age
}

func (u *User) GetName() string {
	u.RLock()
	name := u.name
	u.RUnlock()
	return name
}

func (u *User) MarshalJSON() ([]byte, error) {
	u.RLock()
	userJson, error := json.Marshal(struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Id:   u.id,
		Age:  u.age,
		Name: u.name,
	})
	u.RUnlock()
	return userJson, error
}
