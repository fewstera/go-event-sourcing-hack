package eventsourcing

import (
	"fmt"

	"database/sql"
)

type EventStore struct {
	db *sql.DB
}

func NewEventStore(db *sql.DB) *EventStore {
	eventStore := new(EventStore)
	eventStore.db = db
	return eventStore
}

func (eventStore *EventStore) SaveEvent(event Event) {
	dataBytes, err := event.GetData()
	if err != nil {
		fmt.Printf("WE ERRORED SERIALIZING: %v\n", err)
		return
	}

	_, err = eventStore.db.Exec(
		"INSERT INTO `event` (stream_category, stream_id, event_number, event_type, data) VALUES (?, ?, ?, ?, ?)",
		event.GetStreamCategory(), event.GetStreamId(), event.GetEventNumber(), event.GetEventType(), string(dataBytes),
	)

	if err != nil {
		fmt.Printf("WE ERRORED: %v\n", err)
		return
	}

	fmt.Println("Wrote event to store")
}
