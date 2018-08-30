package eventsourcing

import "encoding/json"

type Event interface {
	GetStreamId() string
	GetEventNumber() int
	GetEventType() string
	GetData() ([]byte, error)
	InitFromDbEvent(streamId string, eventNumber int, data []byte) error
}

type EventImplementation struct {
	EventNumber int
	Id          string
}

// Event types
const (
	EventTypeUserCreated     string = "USER_CREATED"
	EventTypeUserNameChanged string = "USER_NAME_CHANGED"
	EventTypeUserGotOlder    string = "USER_GOT_OLDER"
)

// EVENT METHODS
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
	return EventTypeUserCreated
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

func (e *UserCreatedEvent) InitFromDbEvent(streamId string, eventNumber int, data []byte) error {
	var eventData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err := json.Unmarshal(data, &eventData)
	if err != nil {
		return err
	}

	e.Id = streamId
	e.EventNumber = eventNumber
	e.Name = eventData.Name
	e.Age = eventData.Age
	return nil
}

// USER NAME CHANGE
type UserNameChangedEvent struct {
	EventImplementation
	NewName string
}

func NewUserNameChangedEvent(eventNumber int, id string, newName string) *UserNameChangedEvent {
	userNameChangedEvent := new(UserNameChangedEvent)
	userNameChangedEvent.EventNumber = eventNumber
	userNameChangedEvent.Id = id
	userNameChangedEvent.NewName = newName

	return userNameChangedEvent
}

func (e *UserNameChangedEvent) GetEventType() string {
	return EventTypeUserNameChanged
}

func (e *UserNameChangedEvent) GetData() ([]byte, error) {
	return json.Marshal(struct {
		NewName string `json:"newName"`
	}{
		NewName: e.NewName,
	})
}

func (e *UserNameChangedEvent) InitFromDbEvent(streamId string, eventNumber int, data []byte) error {
	var eventData struct {
		NewName string `json:"newName"`
	}

	err := json.Unmarshal(data, &eventData)
	if err != nil {
		return err
	}

	e.Id = streamId
	e.EventNumber = eventNumber
	e.NewName = eventData.NewName
	return nil
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
	return EventTypeUserGotOlder
}

func (e *UserGotOlderEvent) GetData() ([]byte, error) {
	return []byte("{}"), nil
}

func (e *UserGotOlderEvent) InitFromDbEvent(streamId string, eventNumber int, data []byte) error {
	e.Id = streamId
	e.EventNumber = eventNumber
	return nil
}
