package eventstore

import "fmt"

// EventNumberConflictError is an error for when an event is attempted to
// be written to the store with an already used event number.
type EventNumberConflictError struct {
	StreamID    string
	EventNumber int
}

func (e *EventNumberConflictError) Error() string {
	return fmt.Sprintf("eventstore event number conflict: event number %d already exists on streamID %v", e.EventNumber, e.StreamID)
}
