package eventsourcing

import "encoding/json"

type Event interface {
	GetStreamCategory() string
	GetStreamId() string
	GetEventNumber() int
	GetEventType() string
	GetData() ([]byte, error)
}

type EventImplementation struct {
	EventNumber int
	Id          string
}

// EVENT METHODS
func (e *EventImplementation) GetStreamCategory() string {
	return "USER"
}

func (e *EventImplementation) GetStreamId() string {
	return e.Id
}

func (e *EventImplementation) GetEventNumber() int {
	return e.EventNumber
}

// CREATE USER EVENT
type UserCreatedEvent struct {
	EventImplementation
	Name string
	Age  int
}

func NewUserCreatedEvent(eventNumber int, id string, name string, age int) *UserCreatedEvent {
	userCreatedEvent := new(UserCreatedEvent)
	userCreatedEvent.EventNumber = eventNumber
	userCreatedEvent.Id = id
	userCreatedEvent.Name = name
	userCreatedEvent.Age = age
	return userCreatedEvent
}

func (e *UserCreatedEvent) GetEventType() string {
	return "USER_CREATED"
}

func (e *UserCreatedEvent) GetData() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Age:  e.Age,
		Name: e.Name,
	})
}

// USER NAME CHANGE
type UsersNameChangedEvent struct {
	EventImplementation
	NewName string
}

func NewUsersNameChangedEvent(eventNumber int, id string, newName string) *UsersNameChangedEvent {
	usersNameChangedEvent := new(UsersNameChangedEvent)
	usersNameChangedEvent.EventNumber = eventNumber
	usersNameChangedEvent.Id = id
	usersNameChangedEvent.NewName = newName
	return usersNameChangedEvent
}

func (e *UsersNameChangedEvent) GetEventType() string {
	return "USER_NAME_CHANGE"
}

func (e *UsersNameChangedEvent) GetData() ([]byte, error) {
	return json.Marshal(struct {
		NewName string `json:"name"`
	}{
		NewName: e.NewName,
	})
}

// USER GOT OLDER
type UserGotOlderEvent struct {
	EventImplementation
}

func NewUserGotOlderEvent(eventNumber int, id string) *UserGotOlderEvent {
	userGotOlderEvent := new(UserGotOlderEvent)
	userGotOlderEvent.EventNumber = eventNumber
	userGotOlderEvent.Id = id
	return userGotOlderEvent
}

func (e *UserGotOlderEvent) GetEventType() string {
	return "USER_GOT_OLDER"
}

func (e *UserGotOlderEvent) GetData() ([]byte, error) {
	return []byte("{}"), nil
}
