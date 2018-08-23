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

// Apply methods - These should only mutate state, they are not allowed to error.
func (u *User) apply(event Event) {
	switch e := event.(type) {
	case *UserCreatedEvent:
		u.applyUserCreated(e)
	default:
		fmt.Println(e)
		fmt.Println("Unkown event applied on user")
	}
}

func (u *User) applyUserCreated(e *UserCreatedEvent) {
	u.id = e.id
	u.age = e.age
	u.name = e.name

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
