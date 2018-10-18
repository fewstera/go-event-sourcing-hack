package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

// DepositedEvent is a struct representing a user created event
type DepositedEvent struct {
	eventstore.BaseEvent
	Amount float32 `json:"amount"`
}

// NewDepositedEvent creates an instance of DepositedEvent
func NewDepositedEvent(streamID string, eventNumber int, amount float32) *DepositedEvent {
	userCreatedEvent := new(DepositedEvent)
	userCreatedEvent.StreamID = streamID
	userCreatedEvent.EventNumber = eventNumber
	userCreatedEvent.Amount = amount
	return userCreatedEvent
}

type depositedEventData struct {
	Amount float32 `json:"amount"`
}

// Init initialises an event from an instance of Data.
func (e *DepositedEvent) Init(data *eventstore.EventData) {
	e.StreamID = data.StreamID
	e.EventNumber = data.EventNumber
	e.Timestamp = data.Timestamp

	var eventData depositedEventData
	err := json.Unmarshal([]byte(data.Data), &eventData)
	if err != nil {
		fmt.Printf("Error parsing eventstore data: %v", err)
	}
	e.Amount = eventData.Amount
}

// Data returns an instance of Data so the event can be saved.
func (e *DepositedEvent) Data() (*eventstore.EventData, error) {
	eventData := depositedEventData{e.Amount}
	bytes, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}
	return &eventstore.EventData{e.StreamID, e.EventNumber, EventTypeDeposited, string(bytes), time.Time{}}, nil
}
