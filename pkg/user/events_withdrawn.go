package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

// WithdrawnEvent is a struct representing a user created event
type WithdrawnEvent struct {
	eventstore.BaseEvent
	Amount float32 `json:"amount"`
}

// NewWithdrawnEvent creates an instance of WithdrawnEvent
func NewWithdrawnEvent(streamID string, eventNumber int, amount float32) *WithdrawnEvent {
	userCreatedEvent := new(WithdrawnEvent)
	userCreatedEvent.StreamID = streamID
	userCreatedEvent.EventNumber = eventNumber
	userCreatedEvent.Amount = amount
	return userCreatedEvent
}

type withdrawnEventData struct {
	Amount float32 `json:"amount"`
}

// Init initialises an event from an instance of Data.
func (e *WithdrawnEvent) Init(data *eventstore.EventData) {
	e.StreamID = data.StreamID
	e.EventNumber = data.EventNumber
	e.Timestamp = data.Timestamp

	var eventData withdrawnEventData
	err := json.Unmarshal([]byte(data.Data), &eventData)
	if err != nil {
		fmt.Printf("Error parsing eventstore data: %v", err)
	}
	e.Amount = eventData.Amount
}

// Data returns an instance of Data so the event can be saved.
func (e *WithdrawnEvent) Data() (*eventstore.EventData, error) {
	eventData := withdrawnEventData{e.Amount}
	bytes, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}
	return &eventstore.EventData{e.StreamID, e.EventNumber, EventTypeWithdrawn, string(bytes), time.Time{}}, nil
}
