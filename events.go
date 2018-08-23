package main

type Event interface {
	ImplementsEvent()
}

type UserCreatedEvent struct {
	id   string
	name string
	age  int
}

func (e UserCreatedEvent) ImplementsEvent() {}

type UsersNameChangedEvent struct {
	id      string
	newName string
}

func (e UsersNameChangedEvent) ImplementsEvent() {}

type UserGotOlderEvent struct {
	id string
}

func (e UserGotOlderEvent) ImplementsEvent() {}
