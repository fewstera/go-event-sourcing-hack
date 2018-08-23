package main

import "fmt"

type User struct {
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
	u.apply(&UserCreatedEvent{id, name, age})
	return u, nil
}

// Actions - These are called by command handlers, they can error but should not mutate state (except through calling apply)
func (u *User) IncreaseAge() {
	u.apply(&UserGotOlderEvent{u.id})
}

func (u *User) ChangeName(name string) {
	u.apply(&UsersNameChangedEvent{u.id, name})
}

// Apply methods - These should only mutate state, they are not allowed to error.
func (u *User) apply(event Event) {
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
	u.id = e.id
	u.age = e.age
	u.name = e.name

	u.uncommitedEvents = append(u.uncommitedEvents, e)
}

func (u *User) applyUserGotOlder(e *UserGotOlderEvent) {
	u.age = u.age + 1

	u.uncommitedEvents = append(u.uncommitedEvents, e)
}

func (u *User) applyUsersNameChanged(e *UsersNameChangedEvent) {
	u.name = e.newName

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
