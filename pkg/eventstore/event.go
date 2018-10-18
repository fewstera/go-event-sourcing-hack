package eventstore

import "time"

// Event is the interface implemented by event sourcing events.
type Event interface {
	// Init initialises an event from an instance of data.
	Init(data *EventData)
	// Data returns an Data struct so the event can be persisted
	Data() (*EventData, error)

	// EventNumber returns the event number of the event.
	GetEventNumber() int
	// StreamID returns the StreamID (aggregate ID) of the event
	GetStreamID() string
	GetTimestamp() time.Time
}

// EventData is a struct for representing an event in a persistable manner.
// EventData is to be used by implementers of the Event interface so that the event
// can be saved to an event store or loaded from the event store.
//
// StreamID, EventNumber and EventType are the values for that of the event.
//
// The Data field is used for representing data relating to an event. For example,
// for an UserUpdatedEvent the Data field might contain JSON with the updated users
// details, like
//		{"name": "Bob", "age": 12}
type EventData struct {
	StreamID    string    `json:"streamID"`
	EventNumber int       `json:"eventNumber"`
	EventType   string    `json:"eventType"`
	Data        string    `json:"data"`
	Timestamp   time.Time `json:"-"`
}

// BaseEvent can be embedded by implementers of the Event interface to gain
// implementations for EventNumber() and StreamID().
type BaseEvent struct {
	StreamID    string    `json:"streamID"`
	EventNumber int       `json:"eventNumber"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetEventNumber returns the event number of the event.
func (e *BaseEvent) GetEventNumber() int {
	return e.EventNumber
}

// GetStreamID returns the StreamID (aggregate ID) of the event
func (e *BaseEvent) GetStreamID() string {
	return e.StreamID
}

// GetTimestamp returns the StreamID (aggregate ID) of the event
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}
