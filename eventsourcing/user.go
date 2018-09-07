package eventsourcing

import (
	"encoding/json"
	"fmt"
	"sync"
)

type User struct {
	EventNumber int
	Id          string
	Name        string
	Age         int
	mutex       sync.RWMutex
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
	u.mutex.RLock()
	nextEventNumber := u.EventNumber + 1
	event := NewUserGotOlderEvent(nextEventNumber, u.Id)
	u.mutex.RUnlock()
	return event
}

func (u *User) ChangeName(name string) *UserNameChangedEvent {
	u.mutex.RLock()
	nextEventNumber := u.EventNumber + 1
	event := NewUserNameChangedEvent(nextEventNumber, u.Id, name)
	u.mutex.RUnlock()
	return event
}

// Apply methods - These should only mutate state, they are not allowed to error.
func (u *User) Apply(event Event) {
	u.EventNumber = event.GetEventNumber()

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
	u.mutex.Lock()
	u.Id = e.Id
	u.Age = e.Age
	u.Name = e.Name
	u.mutex.Unlock()
}

func (u *User) applyUserGotOlder(e *UserGotOlderEvent) {
	u.mutex.Lock()
	u.Age = u.Age + 1
	u.mutex.Unlock()
}

func (u *User) applyUsersNameChanged(e *UserNameChangedEvent) {
	u.mutex.Lock()
	u.Name = e.NewName
	u.mutex.Unlock()
}

// Getters
func (u *User) GetId() string {
	u.mutex.RLock()
	id := u.Id
	u.mutex.RUnlock()
	return id
}

func (u *User) GetAge() int {
	u.mutex.RLock()
	age := u.Age
	u.mutex.RUnlock()
	return age
}

func (u *User) GetName() string {
	u.mutex.RLock()
	name := u.Name
	u.mutex.RUnlock()
	return name
}

func (u *User) MarshalJSON() ([]byte, error) {
	u.mutex.RLock()
	userJson, error := json.Marshal(struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Id:   u.Id,
		Age:  u.Age,
		Name: u.Name,
	})
	u.mutex.RUnlock()
	return userJson, error
}
