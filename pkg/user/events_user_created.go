package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

// UserCreatedEvent is a struct representing a user created event
type UserCreatedEvent struct {
	eventstore.BaseEvent
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// NewUserCreatedEvent creates an instance of UserCreatedEvent
func NewUserCreatedEvent(streamID string, eventNumber int, name string, age int) *UserCreatedEvent {
	userCreatedEvent := new(UserCreatedEvent)
	userCreatedEvent.StreamID = streamID
	userCreatedEvent.EventNumber = eventNumber
	userCreatedEvent.Name = name
	userCreatedEvent.Age = age
	return userCreatedEvent
}

type userCreatedEventData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Init initialises an event from an instance of Data.
func (e *UserCreatedEvent) Init(data *eventstore.EventData) {
	e.EventNumber = data.EventNumber
	e.StreamID = data.StreamID
	e.Timestamp = data.Timestamp

	var eventData userCreatedEventData
	err := json.Unmarshal([]byte(data.Data), &eventData)
	if err != nil {
		fmt.Printf("Error parsing eventstore data: %v", err)
	}
	e.Name = eventData.Name
	e.Age = eventData.Age
}

// Data returns an instance of Data so the event can be saved.
func (e *UserCreatedEvent) Data() (*eventstore.EventData, error) {
	eventData := userCreatedEventData{e.Name, e.Age}
	bytes, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}
	return &eventstore.EventData{e.StreamID, e.EventNumber, EventTypeUserCreated, string(bytes), time.Time{}}, nil
}
